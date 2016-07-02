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
}

var WorldACLPermAll = zk.WorldACL(zk.PermAll)

func NewStore() *Storage {
	conn, connEvent, err := zk.Connect([]string{"127.0.0.1"}, defaultLeaderTime)

	if err != nil {
		panic(err)
		log.Error("Zk连接创建失败...", err)
	}

	log.Info("Zk连接创建成功...")

	execIdLock := zk.NewLock(conn, "/swiss/ids/execid", WorldACLPermAll)
	jobIdLock := zk.NewLock(conn, "/swiss/ids/jobid", WorldACLPermAll)
	return &Storage{conn:conn, connEvent:connEvent, execIdLock:execIdLock, jobIdLock:jobIdLock}
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