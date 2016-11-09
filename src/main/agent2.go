package main

import (
	"github.com/hashicorp/serf/serf"
	"github.com/hashicorp/memberlist"
	"fmt"
	"io/ioutil"
	"github.com/Sirupsen/logrus"
)

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

func main(){
	InitLogger(logrus.DebugLevel.String(),"yws2")

	serfConfig := serf.DefaultConfig()
	serfConfig.NodeName = "yws2"
	serfConfig.MemberlistConfig = memberlist.DefaultLANConfig()

	serfConfig.MemberlistConfig.BindAddr = "127.0.0.1"
	serfConfig.MemberlistConfig.BindPort = 7374
	serfConfig.MemberlistConfig.AdvertiseAddr = "127.0.0.1"
	serfConfig.MemberlistConfig.AdvertisePort = 5001

	eventCh := make(chan serf.Event, 64)
	serfConfig.EventCh = eventCh
	serfConfig.LogOutput = ioutil.Discard
	serfConfig.MemberlistConfig.LogOutput = ioutil.Discard

	s, err := serf.Create(serfConfig)
	if err != nil {
		fmt.Println(err)
	}

	go func(){
		serfShutdownCh := s.ShutdownCh()
		log.Info("agent: Listen for events")
		for {
			select {
			case e := <-eventCh:
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
	}()

	fmt.Println(s)

	_,e :=s.Join([]string{"127.0.0.1:7373"},true)
	fmt.Println(e)

	wait := make(chan int)
	<- wait
}
