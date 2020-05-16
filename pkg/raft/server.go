package raft

import (
	"context"
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"log"
	"time"
)

func (n *Node) AppendEntries(ctx context.Context, req *raft.AppendEntriesRequest) (*raft.AppendEntriesResponse, error) {
	log.Printf("AppendEntries received from %d, term %d", req.GetLeaderId(), req.GetTerm())
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if req.GetTerm() < n.CurrentTerm {
		log.Printf("AppendEntries from %d, term %d is staler than me, term %d", req.GetLeaderId(),
			req.GetTerm(), n.CurrentTerm)
		return &raft.AppendEntriesResponse{Term: n.CurrentTerm, Success: false}, nil
	}

	//todo: log part

	log.Printf("AppendEntries record heartbeat from %d at %s, change state to FOLLOWER", req.GetLeaderId(),
		time.Now().String())
	n.LatestHeartbeatAt = time.Now()
	n.LatestHeartbeatFrom = req.GetLeaderId()
	n.State = FOLLOWER

	return &raft.AppendEntriesResponse{Term: n.CurrentTerm, Success: true}, nil
}

func (n *Node) RequestVote(ctx context.Context, req *raft.RequestVoteRequest) (*raft.RequestVoteResponse, error) {
	log.Printf("RequestVote received from %d, term %d", req.GetCandidateId(), req.GetTerm())
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if req.GetTerm() < n.CurrentTerm {
		log.Printf("RequestVote from %d, term %d is staler than me, term %d", req.GetCandidateId(),
			req.GetTerm(), n.CurrentTerm)
		return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
	}

	if req.GetTerm() == n.CurrentTerm {
		if n.State == CANDIDATE || n.State == LEADER {
			log.Printf("RequestVote from %d, we have the same term %d, reject it", req.GetCandidateId(),
				n.CurrentTerm)
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
		}
	}

	if n.State == FOLLOWER {
		// 没vote过或者vote给了同样的candidate，则继续vote
		if n.VotedFor < 0 || n.VotedFor == req.GetCandidateId() {
			n.VotedFor = req.GetCandidateId()
			log.Printf("RequestVote from %d, I haven't voted for anyone, so vote for it",
				req.GetCandidateId())
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: true}, nil
		} else {
			log.Printf("RequestVote from %d, I have voted for %d, reject it",
				req.GetCandidateId(), n.VotedFor)
			return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
		}
	} else {
		log.Printf("RequestVote from %d, it has larger term %d than mine %d, I'm FOLLOWER now",
			req.GetCandidateId(), req.GetTerm(), n.CurrentTerm)
		n.State = FOLLOWER
		n.VotedFor = req.CandidateId
	}

	return &raft.RequestVoteResponse{Term: n.CurrentTerm, VoteGranted: false}, nil
}

func (n *Node) HeartbeatMonitor() {
	t := time.NewTicker(n.HeartbeatCheckInterval)
	log.Printf("start HeartbeatMonitor, check interval %d s", n.HeartbeatCheckInterval/1e9)
	for range t.C {
		n.Mux.Lock()
		if n.State != LEADER && n.LatestHeartbeatAt.Add(n.ElectionTimeout).After(time.Now()) {
			n.State = CANDIDATE
			n.CurrentTerm++
			log.Printf("HeartbeatMonitor election timeout, begin new election phase, term %d",
				n.CurrentTerm)
			n.BeginElection()
		}
		n.Mux.Unlock()
	}
}
