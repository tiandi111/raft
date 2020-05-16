package raft

import (
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
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

type Node struct {
	ID                     int32
	CurrentTerm            int32
	State                  int
	VotedFor               int32
	ReceivedVotes          int
	HeartbeatInterval      time.Duration
	HeartbeatCheckInterval time.Duration
	ElectionTimeout        time.Duration
	LatestHeartbeatAt      time.Time
	LatestHeartbeatFrom    int32
	Config                 *Config
	Clients                map[int32]*Client
	Mux                    *sync.Mutex
}

type Client struct {
	C    raft.RaftClient
	conn *grpc.ClientConn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewNode(config *Config) *Node {
	return &Node{
		ID:                     config.ID,
		CurrentTerm:            0,
		State:                  FOLLOWER,
		VotedFor:               -1,
		ReceivedVotes:          0,
		HeartbeatInterval:      config.HeartbeatInterval * time.Millisecond,
		HeartbeatCheckInterval: config.HeartbeatCheckInterval * time.Millisecond,
		ElectionTimeout:        RandomElectionTimeout(config) * time.Millisecond,
		LatestHeartbeatAt:      time.Now(),
		LatestHeartbeatFrom:    config.ID,
		Config:                 config,
		Mux:                    &sync.Mutex{},
	}
}

func RandomElectionTimeout(config *Config) time.Duration {
	et := rand.Int63n(int64(config.MaxElectionTimeout))
	if et < int64(3*config.HeartbeatInterval) {
		et = 3 * int64(config.HeartbeatInterval)
	}
	return time.Duration(et)
}
