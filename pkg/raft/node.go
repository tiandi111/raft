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
	CurrentTerm           int32
	State                 int32
	VotedFor              int32
	ReceivedVotes         int32
	InElection            bool // InElection == !ElectionResultChannel.Closed
	ElectionResultChannel chan bool
	Node                  *Node
	mux                   *sync.Mutex
}

func (s *Status) LockStatus() {
	s.mux.Lock()
}

func (s *Status) UnlockStatus() {
	s.mux.Unlock()
}

// - PromoteToCandidate
// 		case1: follower -> candidate, start new election term
//		case2: candidate -> candidate, start new election term
// - PromoteToLeader
// 		case1: candidate -> leader, receive majority of votes
// - DemoteToFollower
//		case1: candidate -> follower, discover current leader or new term
// 		case2: leader -> follower, discover new term
func (s *Status) PromoteToCandidate(term int32) {
	s.State = CANDIDATE
	s.CurrentTerm = term
	s.ReceivedVotes = 0
	s.VotedFor = 0
	if s.InElection && s.ElectionResultChannel != nil {
		close(s.ElectionResultChannel)
	}
	s.InElection = true
	s.ElectionResultChannel = make(chan bool, 1)
}

func (s *Status) PromoteToLeader() {
	s.State = LEADER
	s.ReceivedVotes = 0
	s.VotedFor = 0
	if s.InElection && s.ElectionResultChannel != nil {
		close(s.ElectionResultChannel)
	}
	s.InElection = false
	s.ElectionResultChannel = nil
}

func (s *Status) DemoteToFollower(term int32) {
	s.State = FOLLOWER
	s.CurrentTerm = term
	s.ReceivedVotes = 0
	s.VotedFor = 0
	if s.InElection && s.ElectionResultChannel != nil {
		close(s.ElectionResultChannel)
	}
	s.InElection = false
	s.ElectionResultChannel = nil
}

func (s *Status) VoteFor(votedFor int32) {
	s.VotedFor = votedFor
}

func (s *Status) ReceiveVote() {
	s.ReceivedVotes += 1
	if int(s.ReceivedVotes) >= len(s.Node.Clients)/2+1 {
		s.ElectionResultChannel <- true
	}
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
	HeartbeatInterval      time.Duration
	HeartbeatCheckInterval time.Duration
	ElectionTimeout        time.Duration
	*Status
	*Heartbeat
	Config  *Config
	Clients map[int32]*Client
	DoneC   chan struct{}
}

// similar to compare and swap
// Notice: never pass in functions that call EnsureAndDo inside to avoid dead-lock
func (n *Node) EnsureAndDo(term, state int32, fs ...func(node *Node)) bool {
	n.LockStatus()
	defer n.UnlockStatus()

	if n.CurrentTerm == term && n.State == state {
		for _, f := range fs {
			f(n)
		}
		return true
	}

	return false
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
			CurrentTerm:           0,
			State:                 FOLLOWER,
			VotedFor:              -1,
			ReceivedVotes:         0,
			ElectionResultChannel: make(chan bool, 1),
			mux:                   &sync.Mutex{},
		},
		HeartbeatInterval:      config.HeartbeatInterval * time.Millisecond,
		HeartbeatCheckInterval: config.HeartbeatCheckInterval * time.Millisecond,
		ElectionTimeout:        RandomElectionTimeout(config) * time.Millisecond,
		Heartbeat: &Heartbeat{
			LatestHeartbeatAt: time.Now(),
			mux:               &sync.RWMutex{},
		},
		Config:  config,
		Clients: map[int32]*Client{},
		DoneC:   make(chan struct{}),
	}
	node.Status.Node = node
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
	max := int64(config.MaxElectionTimeout)
	et := rand.Int63n(max/2) + max/2
	//if et < int64(30*config.HeartbeatInterval) {
	//	et = 30 * int64(config.HeartbeatInterval)
	//}
	return time.Duration(et)
}
