package main

import (
	"github.com/mitchellh/cli"
	"github.com/hashicorp/serf/serf"
	"fmt"
	"github.com/hashicorp/memberlist"
	"strings"
	"io/ioutil"
	"github.com/Sirupsen/logrus"
	"time"
)

func main(){
	agent := AgentCommand{}
	agent.Run(nil)

	wait := make(chan int)
	<- wait
}

var log = logrus.NewEntry(logrus.New())

func InitLogger(logLevel string, node string) {
	formattedLogger := logrus.New()
	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}

	formattedLogger.Level = level
	log = logrus.NewEntry(formattedLogger).WithField("node", node)
}


type AgentCommand struct {
	Ui               cli.Ui
	Version          string
	ShutdownCh       <-chan struct{}

	serf      *serf.Serf
	eventCh   chan serf.Event
}

func (a *AgentCommand) Help() string {
	helpText := `
Usage: dkron agent [options]
	Run dkron agent

Options:

  -bind=0.0.0.0:8946              Address to bind network listeners to.
  -advertise=bind_addr            Address used to advertise to other nodes in the cluster. By default, the bind address is advertised.
  -http-addr=0.0.0.0:8080         Address to bind the UI web server to. Only used when server.
  -discover=cluster               A cluster name used to discovery peers. On
                                  networks that support multicast, this can be used to have
                                  peers join each other without an explicit join.
  -join=addr                      An initial agent to join with. This flag can be
                                  specified multiple times.
  -node=hostname                  Name of this node. Must be unique in the cluster
  -profile=[lan|wan|local]        Profile is used to control the timing profiles used.
                                  The default if not provided is lan.
  -server=false                   This node is running in server mode.
  -tag key=value                  Tag can be specified multiple times to attach multiple
                                  key/value tag pairs to the given node.
  -keyspace=dkron                 The keyspace to use. A prefix under all data is stored
                                  for this instance.
  -backend=[etcd|consul|zk]       Backend storage to use, etcd, consul or zookeeper. The default
                                  is etcd.
  -backend-machine=127.0.0.1:2379 Backend storage servers addresses to connect to. This flag can be
                                  specified multiple times.
  -encrypt                        Key for encrypting network traffic.
                                  Must be a base64-encoded 16-byte key.
  -ui-dir                         Directory from where to serve Web UI
  -rpc-port=6868                  RPC Port used to communicate with clients. Only used when server.
                                  The RPC IP Address will be the same as the bind address.

  -mail-host                      Mail server host address to use for notifications.
  -mail-port                      Mail server port.
  -mail-username                  Mail server username used for authentication.
  -mail-password                  Mail server password to use.
  -mail-from                      From email address to use.

  -webhook-url                    Webhook url to call for notifications.
  -webhook-payload                Body of the POST request to send on webhook call.
  -webhook-header                 Headers to use when calling the webhook URL. Can be specified multiple times.

  -log-level=info                 Log level (debug, info, warn, error, fatal, panic). Default to info.
`
	return strings.TrimSpace(helpText)
}

// setupSerf is used to create the agent we use
func (a *AgentCommand) setupSerf() *serf.Serf {
	InitLogger(logrus.DebugLevel.String(),"YANG1")

	serfConfig := serf.DefaultConfig()
	serfConfig.NodeName = "yws12"
	serfConfig.MemberlistConfig = memberlist.DefaultLANConfig()

	serfConfig.MemberlistConfig.BindAddr = "127.0.0.1"
	serfConfig.MemberlistConfig.BindPort = 7373
	serfConfig.MemberlistConfig.AdvertiseAddr = "127.0.0.1"
	serfConfig.MemberlistConfig.AdvertisePort = 5000

	eventCh := make(chan serf.Event, 64)
	serfConfig.EventCh = eventCh
	serfConfig.LogOutput = ioutil.Discard
	serfConfig.MemberlistConfig.LogOutput = ioutil.Discard

	s, err := serf.Create(serfConfig)
	if err != nil {
		fmt.Println(err)
	}

	return s
}

func (a *AgentCommand) Run(args []string) int {
	if a.serf = a.setupSerf(); a.serf == nil {
		log.Fatal("agent: Can not setup serf")
	}
	a.join([]string{"127.0.0.1:7373"}, true)

	go a.eventLoop()

	time.Sleep(100000)
	return 1
}

func (a *AgentCommand) eventLoop() {
	serfShutdownCh := a.serf.ShutdownCh()
	log.Info("agent: Listen for events")
	for {
		select {
		case e := <-a.eventCh:
			log.WithFields(logrus.Fields{
				"event": e.String(),
			}).Debug("agent: Received event")

		// Log all member events
			if failed, ok := e.(serf.MemberEvent); ok {
				for _, member := range failed.Members {
					log.WithFields(logrus.Fields{
						"node":   "yang",
						"member": member.Name,
						"event":  e.EventType(),
					}).Debug("agent: Member event")
				}
			}

			if e.EventType() == serf.EventQuery {
				query := e.(*serf.Query)
				fmt.Println(query)
			}

		case <-serfShutdownCh:
			log.Warn("agent: Serf shutdown detected, quitting")
			return
		}
	}
}


// Join asks the Serf instance to join. See the Serf.Join function.
func (a *AgentCommand) join(addrs []string, replay bool) (n int, err error) {
	log.Infof("agent: joining: %v replay: %v", addrs, replay)
	ignoreOld := !replay
	n, err = a.serf.Join(addrs, ignoreOld)
	if n > 0 {
		log.Infof("agent: joined: %d nodes", n)
	}
	if err != nil {
		log.Warnf("agent: error joining: %v", err)
	}
	return
}
