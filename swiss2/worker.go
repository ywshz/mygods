package swiss

import (
"github.com/docker/libkv"
"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

type Worker struct {
	Config     *Config
}

func (s *Worker) Run(args []string) int {
	s.Config = NewConfig()
	etcd.Register()
	client, _ := libkv.NewStore(store.ETCD, s.Config.BackendMachines, &store.Config{})

	changeChan,_ := client.WatchTree("swiss/readytorun/", make(chan struct{}))
	for {
		node := <-changeChan
		log.Info(string(node[len(node)-1].Value))
	}
	return -1;
}

func (s *Worker) Help() string {
	return "Worker"
}

func (s *Worker) Synopsis() string {
	return "Run Worker"
}