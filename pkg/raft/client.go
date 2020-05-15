package raft

import (
	"context"
	"github.com/grpc/grpc-go"
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"log"
	"time"
)

func (n *Node) InitClient() error {
	//todo: close connections
	f := func(id int32, addr string) (raft.RaftClient, error) {
		conn, err := grpc.Dial(addr)
		if err != nil {
			log.Printf("init client%d failed, err: %s", err)
			return nil, err
		}
		return raft.NewRaftClient(conn), nil
	}
	for id, addr := range n.Config.Others {
		c, err := f(id, addr)
		if err != nil {
			log.Printf("InitClient failed err: %s", err)
			return err
		}
		n.Clients[id] = c
	}
}

func (n *Node) LeaderAppendEntries() {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if n.State != LEADER {
		log.Printf("abort LeaderAppendEntries, i'm not leader anymore")
		return
	}
	req := &raft.AppendEntriesRequest{
		Term:         n.CurrentTerm,
		LeaderId:     n.ID,
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		LeaderCommit: 0,
	}
	for id := range n.Clients {
		log.Printf("LeaderAppendEntries to %d", id)
		go n.AppendEntriesToSingle(id, req)
	}
}

func (n *Node) AppendEntriesToSingle(id int32, req *raft.AppendEntriesRequest) {
	tries := 0
	for {
		n.Mux.Lock()
		if n.State != LEADER {
			log.Printf("abort LeaderAppendEntriesToSingle, i'm not leader anymore")
			n.Mux.Unlock()
			return
		}
		log.Printf("AppendEntries to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].AppendEntries(context.TODO(), req)
		if err != nil {
			log.Printf("AppendEntries to %d, try no.%d, err: %s", id, tries, err)
		}
		n.Mux.Lock()
		if resp.GetSuccess() {
			log.Printf("AppendEntries to %d, try no.%d, success", id, tries)
			n.Mux.Unlock()
			return
		} else {
			log.Printf("AppendEntries to %d, try no.%d, fail", id, tries)
			if resp.GetTerm() >= n.CurrentTerm {
				log.Printf("AppendEntries to %d, try no.%d, target has larger term %d, my term %d",
					id, tries, resp.GetTerm(), n.CurrentTerm)
				n.State = FOLLOWER
				n.CurrentTerm = resp.GetTerm()
				n.LatestHeartbeatAt = time.Now()
				n.LatestHeartbeatFrom = id
				n.Mux.Unlock()
				return
			}
		}
		n.Mux.Unlock()
		tries++
	}
}

// todo: this part
//func (n *Node) BeginElection() {
//	n.Mux.Lock()
//	defer n.Mux.Unlock()
//	// vote to myself
//	n.ReceivedVotes = 1
//	n.CandidateRequestVote()
//	for n.State == CANDIDATE && n.ReceivedVotes < (len(n.Clients)+1)/2+1 &&
//		n.LatestHeartbeatAt.Add(n.ElectionTimeOut).After(time.Now()) {
//		n.State = LEADER
//		n.ReceivedVotes = 0
//		n.LeaderAppendEntries()
//	}
//	if n.State != CANDIDATE {
//		return
//	}
//	if n.ReceivedVotes
//}

func (n *Node) CandidateRequestVote() {
	n.Mux.Lock()
	if n.State != CANDIDATE {
		log.Printf("abort CandidateRequestVote, i'm not CANDIDATE anymore")
		n.ReceivedVotes = 0
		n.Mux.Unlock()
		return
	}
	req := &raft.RequestVoteRequest{
		Term:         n.CurrentTerm,
		CandidateId:  n.ID,
		LastLogIndex: 0,
		LastLogTerm:  0,
	}
	for id := range n.Clients {
		log.Printf("CandidateRequestVote to %d", id)
		go n.RequestVoteToSingle(id, req)
	}
}

func (n *Node) RequestVoteToSingle(id int32, req *raft.RequestVoteRequest) {
	tries := 0
	for {
		n.Mux.Lock()
		if n.State != CANDIDATE {
			log.Printf("abort RequestVoteToSingle, i'm not CANDIDATE anymore")
			n.ReceivedVotes = 0
			return
		}
		n.Mux.Unlock()
		log.Printf("RequestVoteToSingle to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].RequestVote(context.TODO(), req)
		if err != nil {
			log.Printf("RequestVoteToSingle to %d, try no.%d, err: %s", id, tries, err)
		}
		n.Mux.Lock()
		if resp.GetVoteGranted() {
			log.Printf("RequestVoteToSingle to %d, try no.%d, vote granted", id, tries)
			n.ReceivedVotes++
			n.Mux.Unlock()
			return
		} else {
			log.Printf("RequestVoteToSingle to %d, try no.%d, fail", id, tries)
			if resp.GetTerm() >= n.CurrentTerm {
				log.Printf("RequestVoteToSingle to %d, try no.%d, target has larger term %d, my term %d",
					id, tries, resp.GetTerm(), n.CurrentTerm)
				n.State = FOLLOWER
				n.CurrentTerm = resp.GetTerm()
				n.LatestHeartbeatAt = time.Now()
				n.LatestHeartbeatFrom = id
				n.Mux.Unlock()
				return
			}
		}
		n.Mux.Unlock()
		tries++
	}
}
