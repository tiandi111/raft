package raft

import (
	"context"
	"github.com/tiandi111/raft/pkg/rpc/raft"
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
	Clients                map[int32]raft.RaftClient
	Mux                    sync.Mutex
}

func (n *Node) AppendEntries(ctx context.Context, req *raft.AppendEntriesRequest) (*raft.AppendEntriesResponse, error) {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if req.GetTerm() < n.CurrentTerm {
		return &raft.AppendEntriesResponse{Term: n.CurrentTerm, Success: false}, nil
	}

	//todo: log part

	n.LatestHeartbeatAt = time.Now()
	n.LatestHeartbeatFrom = req.GetLeaderId()
	n.State = FOLLOWER

	return &raft.AppendEntriesResponse{Term: n.CurrentTerm, Success: true}, nil
}

func (n *Node) RequestVote(ctx context.Context, req *raft.RequestVoteRequest) (*raft.RequestVoteResponse, error) {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if req.GetTerm() < n.CurrentTerm {
		return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
	}

	if req.GetTerm() == n.CurrentTerm {
		if n.State == CANDIDATE || n.State == LEADER {
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
		}
	}

	if n.State == FOLLOWER {
		// 没vote过或者vote给了同样的candidate，则继续vote
		if n.VotedFor < 0 || n.VotedFor == req.GetCandidateId() {
			n.VotedFor = req.GetCandidateId()
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: true}, nil
		} else {
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
		}
	} else {
		n.State = FOLLOWER
		n.VotedFor = req.CandidateId
	}

	return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
}

func (n *Node) HeartbeatMonitor() {
	t := time.NewTicker(n.HeartbeatCheckInterval)
	for range t.C {
		n.Mux.Lock()
		if n.State != LEADER && n.LatestHeartbeatAt.Add(n.ElectionTimeOut).After(time.Now()) {
			n.State = CANDIDATE
			n.CurrentTerm++
			n.BeginElection()
		}
		n.Mux.Unlock()
	}
}
