package swiss

import (
	"github.com/docker/libkv/store"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store/etcd"
	"github.com/coreos/etcd/client"
	"time"
	"strconv"
)

const MaxExecutions = 100
//
type Storage struct {
	Client   store.Store
	Kapi     client.KeysAPI
	command  *ServerCommand
	keyspace string
}

//
func NewStore(machines []string, cmd *ServerCommand, keyspace string) *Storage {
	etcd.Register()

	clt, err := libkv.NewStore(store.ETCD, cmd.Config.BackendMachines, &store.Config{})
	if err != nil {
		panic(err)
	}

	_, err = clt.List(keyspace)
	if err != store.ErrKeyNotFound && err != nil {
		log.WithError(err).Fatal("store: Store backend not reachable")
	}

	cfg := client.Config{
		Endpoints:               machines,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	return &Storage{Client: clt, Kapi : kapi, command: cmd, keyspace: keyspace}
}
//
// Retrieve the leader from the store
func (s *Storage) GetLeader() []byte {
	res, err := s.Client.Get(s.LeaderKey())
	if err != nil {
		if err == store.ErrNotReachable {
			log.Fatal("store: Store not reachable, be sure you have an existing key-value store running is running and is reachable.")
		} else if err != store.ErrKeyNotFound {
			log.Error(err)
		}
		return nil
	}

	log.WithField("node", string(res.Value)).Debug("store: Retrieved leader from datastore")

	return res.Value
}

//// Retrieve the leader key used in the KV store to store the leader node
func (s *Storage) LeaderKey() string {
	return s.keyspace + "/leader"
}

func (s *Storage) ListJobs() []*Job {
	//s.Client.List(s.keyspace + "/jobs")
	jobs := make([]*Job, 1)
	jobs[0] = &Job{
		Id:1,
		Name:"test",
		ScheduleType: 0,
		Cron : "*/5 * * * * *",
		ScriptType: 0,
		Server: s.command,
	}
	return jobs
}

func (s *Storage) GetJob(jobId string) *Job {
	s.Client.List(s.keyspace + "/jobs")
	return nil
}

func (s *Storage) getJobScript(jobId string) string {
	s.Client.List(s.keyspace + "/jobs")
	return ""
}

func (s *Storage) UpdateJob() error {
	s.Client.List(s.keyspace + "/jobs")
	return nil
}

func (s *Storage) NextExecutionId(key string) (int, error) {
	lock, error := s.Client.NewLock(key, nil)
	defer lock.Unlock()
	if error != nil {
		return 0, error
	}

	prev, error := s.Client.Get(key)
	prevInt, _ := strconv.Atoi(string(prev.Value[:len(prev.Value)]))
	newInt := prevInt + 1
	newStr := strconv.Itoa(newInt)
	s.Client.Put(key, []byte(newStr), nil)

	return newInt, nil
}

func (s *Storage) Set(key string, value []byte) error {
	log.Info(s.keyspace + key)
	return s.Client.Put(s.keyspace + key, value, &store.WriteOptions{IsDir:false})
}