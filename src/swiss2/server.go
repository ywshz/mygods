package swiss

import (
	"github.com/mitchellh/cli"
	"github.com/docker/leadership"

	"time"
	"expvar"
	"os/signal"
	"os"
	"syscall"
	"fmt"
	"strings"
)

const (
	gracefulTimeout = 3 * time.Second
	defaultRecoverTime = 10 * time.Second
	defaultLeaderTTL = 10 * time.Second
)

var (
	expNode = expvar.NewString("node")
)

type ServerCommand struct {
	Ui         cli.Ui
	Version    string
	Config     *Config
	Store      *Storage
	Scheduler  *Scheduler
	ShutdownCh <-chan struct{}
	Candidate  *leadership.Candidate //用于选举leader
}

func (s *ServerCommand) Run(args []string) int {
	s.Config = NewConfig()
	s.Store = NewStore(s.Config.BackendMachines, s, s.Config.Keyspace)

	s.Scheduler = NewScheduler()

	log.Info(s.Store.Client)
	s.Candidate = leadership.NewCandidate(s.Store.Client, s.Store.LeaderKey(), "underwood", defaultLeaderTTL)
	go func() {
		for {
			s.runForElection()
			//当服务器出现故障之类的问题,休息X时间后重试
			time.Sleep(defaultRecoverTime)
		}
	}()
	//监听exit
	return s.handleSignals()
}

func (s *ServerCommand) Help() string {
	helpText := `HELP`
	return strings.TrimSpace(helpText)
}

func (s *ServerCommand) Synopsis() string {
	return "Run dkron"
}

func (s *ServerCommand) schedule() {
	if s.Scheduler.Started {
		s.Scheduler.Restart(s.Store.ListJobs())
	} else {
		s.Scheduler.Start(s.Store.ListJobs())
	}
}

func (s *ServerCommand) stopSchedule() {
	if s.Scheduler.Started {
		s.Scheduler.Stop()
	}
}
// Leader election routine
func (s *ServerCommand) runForElection() {
	log.Info("server: Running for election")
	electedCh, errCh := s.Candidate.RunForElection()

	for {
		select {
		case isElected := <-electedCh:
			if isElected {
				log.Info("1server: Cluster leadership acquired")
				// If this server is elected as the leader, start the scheduler
				s.schedule()
			} else {
				log.Info("server: Cluster leadership lost")
				// Stop the schedule of this server
				s.stopSchedule()
			}

		case err := <-errCh:
			log.WithError(err).Debug("Leader election failed, channel is probably closed")
			return
		}
	}
}


// handleSignals blocks until we get an exit-causing signal
func (s *ServerCommand) handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait for a signal
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-s.ShutdownCh:
		sig = os.Interrupt
	}
	s.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

	// Check if we should do a graceful leave
	graceful := false
	if sig == syscall.SIGTERM || sig == os.Interrupt {
		graceful = true
	}

	// Bail fast if not doing a graceful leave
	if !graceful {
		return 1
	}

	// Attempt a graceful leave
	gracefulCh := make(chan struct{})
	s.Ui.Output("Gracefully shutting down agent...")
	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		// If we're exiting a server
		s.Candidate.Stop()
		close(gracefulCh)
	}()

	// Wait for leave or another signal
	select {
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}