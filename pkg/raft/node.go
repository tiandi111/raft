package raft

import (
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	FOLLOWER = iota
	CANDIDATE
	LEADER
)

type Config struct {
	ID                     int32
	Addr                   string
	HeartbeatInterval      time.Duration
	HeartbeatCheckInterval time.Duration
	MaxElectionTimeout     time.Duration
	Others                 map[int32]string
}

type Status struct {
	CurrentTerm int32
	State       int
	mux         *sync.RWMutex
}

func (s *Status) LockStatus() {
	s.mux.Lock()
}

func (s *Status) UnlockStatus() {
	s.mux.Unlock()
}

func (s *Status) RLockStatus() {
	s.mux.RLock()
}

func (s *Status) RUnlockStatus() {
	s.mux.RUnlock()
}

type Heartbeat struct {
	LatestHeartbeatAt   time.Time
	LatestHeartbeatFrom int32
	mux                 *sync.RWMutex
}

func (h *Heartbeat) LockHeartbeat() {
	h.mux.Lock()
}

func (h *Heartbeat) UnlockHeartbeat() {
	h.mux.Unlock()
}

func (h *Heartbeat) RLockHeartbeat() {
	h.mux.RLock()
}

func (h *Heartbeat) RUnlockHeartbeat() {
	h.mux.RUnlock()
}

type Node struct {
	ID                     int32
	VotedFor               int32
	ReceivedVotes          int32
	HeartbeatInterval      time.Duration
	HeartbeatCheckInterval time.Duration
	ElectionTimeout        time.Duration
	*Status
	*Heartbeat
	Config  *Config
	Clients map[int32]*Client
}

func (n *Node) LoadVotedFor() int32 {
	return atomic.LoadInt32(&n.VotedFor)
}

func (n *Node) StoreVotedFor(v int32) {
	atomic.StoreInt32(&n.VotedFor, v)
}

func (n *Node) LoadReceivedVotes() int32 {
	return atomic.LoadInt32(&n.ReceivedVotes)
}

func (n *Node) StoreReceivedVotes(v int32) {
	atomic.StoreInt32(&n.ReceivedVotes, v)
}

func (n *Node) AddReceivedVotes(delta int32) {
	atomic.AddInt32(&n.ReceivedVotes, delta)
}

type Client struct {
	C    raft.RaftClient
	conn *grpc.ClientConn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewNode(config *Config) *Node {
	node := &Node{
		ID: config.ID,
		Status: &Status{
			CurrentTerm: 0,
			State:       FOLLOWER,
			mux:         &sync.RWMutex{},
		},
		VotedFor:               -1,
		ReceivedVotes:          0,
		HeartbeatInterval:      config.HeartbeatInterval * time.Millisecond,
		HeartbeatCheckInterval: config.HeartbeatCheckInterval * time.Millisecond,
		ElectionTimeout:        RandomElectionTimeout(config) * time.Millisecond,
		Heartbeat: &Heartbeat{
			LatestHeartbeatAt: time.Now(),
			mux:               &sync.RWMutex{},
		},
		Config:  config,
		Clients: map[int32]*Client{},
	}
	log.Printf("node config:\n"+
		"id:[%d]\n"+
		"addr:[%s]\n"+
		"heartbeat_interval:[%d]\n"+
		"heartbeat_check_interval:[%d]\n"+
		"max_election_timeout:[%d]", node.ID, node.Config.Addr,
		node.HeartbeatInterval, node.HeartbeatCheckInterval, node.ElectionTimeout)
	return node
}

func RandomElectionTimeout(config *Config) time.Duration {
	et := rand.Int63n(int64(config.MaxElectionTimeout))
	if et < int64(3*config.HeartbeatInterval) {
		et = 3 * int64(config.HeartbeatInterval)
	}
	return time.Duration(et)
}
