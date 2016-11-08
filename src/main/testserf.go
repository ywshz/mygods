package main

import (
	"github.com/hashicorp/serf/serf"
	"fmt"
	"github.com/hashicorp/memberlist"
	"io/ioutil"
)

func main() {
	go func() {
		serfConfig := serf.DefaultConfig()
		serfConfig.NodeName = "yws5"
		serfConfig.MemberlistConfig = memberlist.DefaultLANConfig()

		serfConfig.MemberlistConfig.BindAddr = "127.0.0.1"
		serfConfig.MemberlistConfig.BindPort = 7377
		serfConfig.MemberlistConfig.AdvertiseAddr = "127.0.0.1"
		serfConfig.MemberlistConfig.AdvertisePort = 5005

		eventCh := make(chan serf.Event, 64)
		serfConfig.EventCh = eventCh
		serfConfig.LogOutput = ioutil.Discard
		serfConfig.MemberlistConfig.LogOutput = ioutil.Discard

		s, err := serf.Create(serfConfig)
		if err != nil {
			fmt.Println(err)
		}

		n, err := s.Join([]string{"127.0.0.1:5003"}, false)

		if n > 0 {
			fmt.Println("agent: joined: %d nodes", n)
		}
		if err != nil {
			fmt.Println("agent: error joining: %v", err)
		}

		fmt.Println(s.Members())
		go func() {

			serfShutdownCh := s.ShutdownCh()
			fmt.Println("agent: Listen for events")

			for {
				select {
				case e := <-eventCh:
					fmt.Println("agent: Received event", e)

					if failed, ok := e.(serf.MemberEvent); ok {
						for _, member := range failed.Members {
							fmt.Println("agent: Member event", member)
						}
					}

					if e.EventType() == serf.EventQuery {
						query := e.(*serf.Query)
						query.Respond([]byte("127.0.0.1"))
					}
				case <-serfShutdownCh:
					return
				}
			}
		}()
	}()

	wait := make(chan struct{})
	<-wait
}