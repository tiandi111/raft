package raft

import (
	"context"
	"fmt"
	"github.com/tiandi111/raft/pkg/rpc/raft"
	"google.golang.org/grpc"
	"log"
	"time"
)

func (n *Node) InitClient() error {
	f := func(id int32, addr string) (*Client, error) {
		var ferr error
		for i := 0; ; i++ {
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Printf("init client%d failed, try no.%d, err: %s", id, i, err)
				ferr = err
				time.Sleep(2 * time.Second)
				continue
			}
			return &Client{C: raft.NewRaftClient(conn), conn: conn}, nil
		}
		return nil, fmt.Errorf("InitClient to node%d failed, err: %s", id, ferr)
	}
	for id, addr := range n.Config.Others {
		c, err := f(id, addr)
		if err != nil {
			log.Printf("InitClient failed err: %s", err)
			return err
		}
		n.Clients[id] = c
	}
	return nil
}

func (n *Node) LeaderAppendEntries() {
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
		if n.State != LEADER {
			log.Printf("abort LeaderAppendEntriesToSingle, i'm not leader anymore")
			return
		}
		log.Printf("AppendEntries to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].C.AppendEntries(context.TODO(), req)
		if err != nil {
			log.Printf("AppendEntries to %d, try no.%d, err: %s", id, tries, err)
		}
		if resp.GetSuccess() {
			log.Printf("AppendEntries to %d, try no.%d, success", id, tries)
			return
		} else {
			log.Printf("AppendEntries to %d, try no.%d, fail", id, tries)
			if resp.GetTerm() >= n.CurrentTerm {
				log.Printf("AppendEntries to %d, try no.%d, target has larger term %d than mine %d",
					id, tries, resp.GetTerm(), n.CurrentTerm)
				n.State = FOLLOWER
				n.CurrentTerm = resp.GetTerm()
				return
			}
		}
		tries++
	}
}

// todo: this part
func (n *Node) BeginElection() {
	// vote to myself
	n.StoreReceivedVotes(1)

	n.RLockStatus()
	term := n.CurrentTerm
	n.RUnlockStatus()

	n.CandidateRequestVote(term)

	for {
		if n.State != CANDIDATE || n.CurrentTerm != term {
			return
		}

		//todo: use condition
		if n.ReceivedVotes >= int32((len(n.Clients)+1)/2+1) {
			n.State = LEADER
			n.ReceivedVotes = 0
			n.LeaderAppendEntries()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (n *Node) CandidateRequestVote(term int32) {
	if n.State != CANDIDATE || n.CurrentTerm != term {
		log.Printf("abort CandidateRequestVote, inconsistent state, state %d, CurentTerm %d, expected term %d",
			n.State, n.CurrentTerm, term)
		n.ReceivedVotes = 0
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
		go n.RequestVoteToSingle(id, term, req)
	}
}

func (n *Node) RequestVoteToSingle(id, term int32, req *raft.RequestVoteRequest) {
	tries := 0
	for {
		if n.State != CANDIDATE || n.CurrentTerm != term {
			log.Printf("abort CandidateRequestVote, inconsistent state, state %d, CurentTerm %d, expected term %d",
				n.State, n.CurrentTerm, term)
			n.ReceivedVotes = 0
			return
		}
		log.Printf("RequestVoteToSingle to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].C.RequestVote(context.TODO(), req)
		if err != nil {
			log.Printf("RequestVoteToSingle to %d, try no.%d, err: %s", id, tries, err)
		}
		if resp.GetVoteGranted() {
			log.Printf("RequestVoteToSingle to %d, try no.%d, vote granted", id, tries)
			n.ReceivedVotes++
			return
		} else {
			log.Printf("RequestVoteToSingle to %d, try no.%d, fail", id, tries)
			if resp.GetTerm() >= n.CurrentTerm {
				log.Printf("RequestVoteToSingle to %d, try no.%d, target has larger term %d, my term %d",
					id, tries, resp.GetTerm(), n.CurrentTerm)
				n.State = FOLLOWER
				n.CurrentTerm = resp.GetTerm()
				return
			}
		}
		tries++
	}
}

func (n *Node) LeaderHeartbeater() {
	ticker := time.NewTicker(time.Millisecond * 1500)
	log.Printf("start LeaderHeartbeater")
	for range ticker.C {
		if n.State == LEADER {
			log.Printf("LeaderHeartbeater send heartbeat")
			n.LeaderAppendEntries()
		}
	}
}
