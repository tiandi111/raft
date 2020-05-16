package raft

import (
	"github.com/grpc/grpc-go"
	"github.com/tiandi111/raft/pkg/rpc/raft"
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
	ID     int32
	Addr   string
	Others map[int32]string
}

type Node struct {
	ID                     int32
	CurrentTerm            int32
	State                  int
	VotedFor               int32
	ReceivedVotes          int
	HeartbeatCheckInterval time.Duration
	ElectionTimeOut        time.Duration
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
		HeartbeatCheckInterval: 2 * time.Second,
		ElectionTimeOut:        time.Duration(rand.Intn(10)+3) * time.Second,
		LatestHeartbeatAt:      time.Now(),
		LatestHeartbeatFrom:    config.ID,
		Mux:                    &sync.Mutex{},
	}
}
