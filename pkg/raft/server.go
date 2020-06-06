package raft

import (
	"context"
	"fmt"
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"github.com/tiandi111/raft/util"
	"log"
	"time"
)

func (n *Node) AppendEntries(ctx context.Context, req *raft.AppendEntriesRequest) (*raft.AppendEntriesResponse, error) {
	log.Printf("AppendEntries received from %d, term %d", req.GetLeaderId(), req.GetTerm())

	n.LockStatus()
	defer n.UnlockStatus()

	term := n.CurrentTerm

	if req.GetTerm() < term {
		log.Printf("AppendEntries from %d, term %d is staler than me, term %d", req.GetLeaderId(),
			req.GetTerm(), term)
		return &raft.AppendEntriesResponse{Term: term, Success: false}, nil
	}

	log.Printf("AppendEntries record heartbeat from %d at %s, change state to FOLLOWER", req.GetLeaderId(),
		time.Now().String())

	n.RLockHeartbeat()
	n.LatestHeartbeatAt = time.Now()
	n.LatestHeartbeatFrom = req.GetLeaderId()
	n.DemoteToFollower(req.GetTerm())
	n.RUnlockHeartbeat()

	if req.Entry != nil {
		if len(n.Logs)-1 < int(req.PrevLogIndex) {
			log.Printf("AppendEntries I don't have log entry %d in my local, reject it", req.PrevLogIndex)
			return &raft.AppendEntriesResponse{Term: term, Success: false}, nil
		} else {
			localPrevLog := n.Logs[req.PrevLogIndex]
			if localPrevLog.Term != req.PrevLogTerm {
				log.Printf("AppendEntries localPrevLog has different term with leader's, reject it. local %d, leader: %d",
					localPrevLog.Term, req.PrevLogTerm)
				return &raft.AppendEntriesResponse{Term: term, Success: false}, nil
			} else {
				log.Printf("AppendEntries accept new entry %v", *req.Entry)
				n.Logs = n.Logs[:req.PrevLogIndex+1]
				n.Logs = append(n.Logs, &LogEntry{
					Term:    req.Entry.Term,
					Command: req.Entry.Command,
				})
				if req.LeaderCommit > n.CommitIndex {
					n.CommitIndex = util.MinInt32(req.LeaderCommit, req.PrevLogIndex+1)
					log.Printf("AppendEntries update commit index to %d", n.CommitIndex)
				}
			}
		}
	}

	return &raft.AppendEntriesResponse{Term: term, Success: true}, nil
}

func (n *Node) RequestVote(ctx context.Context, req *raft.RequestVoteRequest) (*raft.RequestVoteResponse, error) {
	log.Printf("RequestVote received from %d, term %d", req.GetCandidateId(), req.GetTerm())

	n.LockStatus()
	defer n.UnlockStatus()

	term := n.CurrentTerm

	if req.GetTerm() < term {
		log.Printf("RequestVote from %d, term %d is staler than me, term %d", req.GetCandidateId(),
			req.GetTerm(), term)
		return &raft.RequestVoteResponse{Term: term, VoteGranted: false}, nil
	}

	if req.GetTerm() == term {
		if n.State == CANDIDATE || n.State == LEADER {
			log.Printf("RequestVote from %d, we have the same term %d, reject it", req.GetCandidateId(),
				term)
			return &raft.RequestVoteResponse{Term: term, VoteGranted: false}, nil
		}
	}

	if n.State == FOLLOWER {
		// 没vote过或者vote给了同样的candidate，且Candidate的log更up-to-date，则继续vote
		if (n.VotedFor <= 0 || n.VotedFor == req.GetCandidateId()) && n.IsCandidateLogMoreUpToDate(req.GetLastLogIndex(), req.GetLastLogTerm()) {
			log.Printf("RequestVote from %d, I haven't voted for anyone, so vote for it",
				req.GetCandidateId())
			n.VoteFor(req.GetCandidateId())
			return &raft.RequestVoteResponse{Term: term, VoteGranted: true}, nil
		} else {
			log.Printf("RequestVote from %d, I have voted %d, reject it",
				req.GetCandidateId(), n.VotedFor)
			return &raft.RequestVoteResponse{Term: term, VoteGranted: false}, nil
		}
	} else {
		log.Printf("RequestVote from %d, it has larger term %d than mine %d, I'm FOLLOWER now",
			req.GetCandidateId(), req.GetTerm(), term)
		n.DemoteToFollower(req.GetTerm())
	}

	return &raft.RequestVoteResponse{Term: term, VoteGranted: false}, nil
}

func (n *Node) HeartbeatMonitor() {
	log.Printf("start HeartbeatMonitor, check interval %ds", n.HeartbeatCheckInterval/1e9)
	for {
		select {
		case <-time.After(n.HeartbeatCheckInterval):
			n.LockStatus()
			n.RLockHeartbeat()

			if n.State != LEADER && n.LatestHeartbeatAt.Add(n.ElectionTimeout).Before(time.Now()) {

				newterm := n.CurrentTerm + 1
				n.PromoteToCandidate(newterm)

				log.Printf("HeartbeatMonitor election timeout, begin new election phase, term %d",
					n.CurrentTerm)

				n.UnlockStatus()
				n.RUnlockHeartbeat()

				go n.BeginElection(newterm)
			} else {
				n.UnlockStatus()
				n.RUnlockHeartbeat()
			}
		case <-n.DoneC:
			return
		}
	}
}

func (n *Node) Reporter() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		n.Report()
	}
}

func (n *Node) Report() {
	var state string
	switch n.State {
	case LEADER:
		state = "LEADER"
	case CANDIDATE:
		state = "CANDIDATE"
	case FOLLOWER:
		state = "FOLLOWER"
	}
	log.Printf("Reporter term[%d] state[%s] commitIndex[%d] logs[%s]", n.CurrentTerm, state,
		n.CommitIndex, SprintLogs(n.Logs))
}

func SprintLogs(logs []*LogEntry) string {
	str := ""
	for _, entry := range logs {
		str = fmt.Sprintf("%s %s", str, entry.Command)
	}
	return str
}
