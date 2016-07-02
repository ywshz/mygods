package swiss

import (
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

type Candidate struct {
	zkLock     *zk.Lock
	mutex      *sync.RWMutex
	leader     bool
	key        string
	store      *Storage
	electedChn chan bool
	stopChn    chan struct{}
	stopRenew  chan struct{}
	resignChn  chan bool
	errChn     chan error
}

func NewCandidate(store *Storage, leaderKey string) *Candidate {
	return &Candidate{
		store: store,
		key:    leaderKey,

		leader:   false,
		resignChn: make(chan bool),
		stopChn:   make(chan struct{}),
		errChn: make(chan error),
		electedChn: make(chan bool),
	}
}

func (c *Candidate) RunForElection() (<-chan bool, <-chan error) {
	c.mutex = &sync.RWMutex{}
	go c.campaign()
	return c.electedChn, c.errChn
}

func (c *Candidate) Stop() {

}

func (c *Candidate) IsLeader() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.leader
}

func (c *Candidate) setLeader(leader bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.leader = leader
	c.electedChn <- leader
}

func (c *Candidate) Resign() {
	close(c.stopChn)
	c.zkLock.Unlock()
	c.zkLock = nil
	c.setLeader(false)
}

func (c *Candidate) campaign() {
	defer close(c.electedChn)
	defer close(c.errChn)

	for {

		c.setLeader(false)

		if c.zkLock == nil {
			c.zkLock = zk.NewLock(c.store.conn, c.key, zk.WorldACL(zk.PermAll))
		}

		err := c.zkLock.Lock()
		if err != nil {
			log.Info("Get Lock error :", err)
			c.errChn <- err
			return
		}

		//good, we are now the leader
		log.Info("Get the leader lock success")
		c.setLeader(true)

		select {
		case <-c.resignChn:
			c.zkLock.Unlock()
		case <-c.stopChn:
			if c.IsLeader() {
				c.zkLock.Lock()
			}
			return
		}
	}

}
