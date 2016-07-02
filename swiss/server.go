package swiss

import (
	"time"
	"os"
	"os/signal"
	"syscall"
)

const (
	gracefulTimeout = 3 * time.Second
	defaultRecoverTime = 10 * time.Second
	defaultLeaderTime = 10 * time.Second
)

type Server struct {
	Store      *Storage
	ShutdownCh <-chan struct{}
	candidate  *Candidate //用于选举leader
}

func (s *Server) Run(args []string) int {
	log.Info("连接zk Store...")
	s.Store = NewStore()
	log.Info("连接zk Store...成功")

	log.Info("创建候选人...")
	s.candidate = NewCandidate(s.Store, s.Store.LeaderKey())
	log.Info("创建候选人...成功")

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

func (s *Server) runForElection() {
	log.Info("参与竞选...")
	electedCh, errCh := s.candidate.RunForElection()

	log.Info("等待竞选结果...")
	for {
		select {
		case isElected := <-electedCh:
			if isElected {
				log.Info("server: Cluster leadership acquired")
				// If this server is elected as the leader, start the scheduler
				//s.schedule()
			} else {
				log.Info("server: Cluster leadership lost")
				// Stop the schedule of this server
				//s.stopSchedule()
			}
		case err := <-errCh:
			log.WithError(err).Debug("Leader election failed, channel is probably closed")
			return
		}
	}
}

// handleSignals blocks until we get an exit-causing signal
func (s *Server) handleSignals() int {
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
	//log.Info("agent: Gracefully shutting down agent...")
	go func() {
		// If we're exiting a server
		s.candidate.Stop()
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