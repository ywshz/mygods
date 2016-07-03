package swiss

import (
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
)

type Storage struct {
	conn       *zk.Conn
	connEvent  <- chan zk.Event
	execIdLock *zk.Lock
	jobIdLock  *zk.Lock
	server     *Server
}

var WorldACLPermAll = zk.WorldACL(zk.PermAll)

func NewStore(s *Server) *Storage {
	conn, connEvent, err := zk.Connect([]string{"127.0.0.1"}, defaultLeaderTime)

	if err != nil {
		panic(err)
		log.Error("Zk连接创建失败...", err)
	}

	log.Info("Zk连接创建成功...")

	execIdLock := zk.NewLock(conn, "/swiss/ids/execid", WorldACLPermAll)
	jobIdLock := zk.NewLock(conn, "/swiss/ids/jobid", WorldACLPermAll)
	store := &Storage{conn:conn, connEvent:connEvent, execIdLock:execIdLock, jobIdLock:jobIdLock, server:s}
	store.initBasePath()
	return store
}

func (s *Storage) initBasePath() {
	if exist, _, _ := s.conn.Exists("/swiss/readytorun"); !exist {
		log.Info("/swiss/readytorun not exist, 创建中")
		log.Info(s.conn.Create("/swiss/readytorun", []byte(""), 0, WorldACLPermAll))
	}
}

func (s *Storage) LeaderKey() string {
	return "/swiss/leader"
}

func (s *Storage) NextExecutionId() (int, error) {

	err := s.execIdLock.Lock()
	defer s.execIdLock.Unlock()

	if err != nil {
		return 0, err
	}

	prev, _, _ := s.conn.Get("/swiss/ids/execid")
	prevInt, _ := strconv.Atoi(string(prev[:len(prev)]))
	newInt := prevInt + 1
	newStr := strconv.Itoa(newInt)
	s.conn.Set("/swiss/ids/execid", []byte(newStr), -1)
	return newInt, nil
}

func (s *Storage) NextJobId() (int, error) {

	err := s.execIdLock.Lock()
	defer s.execIdLock.Unlock()

	if err != nil {
		return 0, err
	}

	prev, _, _ := s.conn.Get("/swiss/ids/jobid")
	prevInt, _ := strconv.Atoi(string(prev[:len(prev)]))
	newInt := prevInt + 1
	newStr := strconv.Itoa(newInt)
	s.conn.Set("/swiss/ids/jobid", []byte(newStr), -1)
	return newInt, nil
}

func (s *Storage) ListJobs() []*Job {
	//s.Client.List(s.keyspace + "/jobs")
	jobs := make([]*Job, 1)
	jobs[0] = &Job{
		Id:1,
		Name:"test",
		ScheduleType: 0,
		Cron : "*/2 * * * * *",
		ScriptType: 0,
		Server: s.server,
	}
	return jobs
}

func (s *Storage) Create(path string, data []byte) (string, error) {
	log.Debug("Create key,value:", path, "data", data)
	return s.conn.Create(path, data, 0, WorldACLPermAll)
}

func (s *Storage) Set(path string, data []byte) error {
	log.Debug("Set key,value:", path, "data", data)
	_, err := s.conn.Set(path, []byte(data), -1)
	return err
}

func (s *Storage) Get(path string) (string, error) {
	data, _, err := s.conn.Get(path)
	return string(data), err
}