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

func (n *Node) LeaderAppendEntries(term int32) {
	f := func(node *Node) {
		req := &raft.AppendEntriesRequest{
			Term:         node.CurrentTerm,
			LeaderId:     node.ID,
			PrevLogIndex: 0,
			PrevLogTerm:  0,
			LeaderCommit: 0,
		}
		for id := range node.Clients {
			log.Printf("LeaderAppendEntries to %d", id)
			go node.AppendEntriesToSingle(id, node.CurrentTerm, req)
		}
	}
	if !n.EnsureAndDo(term, LEADER, f) {
		log.Printf("abort LeaderAppendEntries, i'm not leader anymore")
	}
}

func (n *Node) AppendEntriesToSingle(id, term int32, req *raft.AppendEntriesRequest) {
	tries := 0
	for {
		if !n.EnsureAndDo(term, LEADER) {
			log.Printf("abort LeaderAppendEntriesToSingle, i'm not leader anymore")
			return
		}

		log.Printf("AppendEntries to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].C.AppendEntries(context.TODO(), req)
		if err != nil {
			log.Printf("AppendEntries to %d, try no.%d, err: %s", id, tries, err)
		}

		shouldReturn := false
		if resp.GetSuccess() {
			log.Printf("AppendEntries to %d, try no.%d, success", id, tries)
			shouldReturn = true
		} else {
			log.Printf("AppendEntries to %d, try no.%d, fail", id, tries)
			if resp.GetTerm() > n.CurrentTerm {
				log.Printf("AppendEntries to %d, try no.%d, target has larger term %d than mine %d, demote to follower",
					id, tries, resp.GetTerm(), n.CurrentTerm)
				n.DemoteToFollower(resp.GetTerm())
				shouldReturn = true
			}
		}

		if !n.EnsureAndDo(term, LEADER) {
			shouldReturn = true
		}
		if shouldReturn {
			return
		}

		tries++
		time.Sleep(300 * time.Millisecond)
	}
}

func (n *Node) BeginElection(term int32) {
	// vote to myself
	if !n.EnsureAndDo(term, CANDIDATE, func(node *Node) {
		log.Printf("BeginElection inconsistent state %d term %d", n.State, n.CurrentTerm)
		node.VoteFor(node.ID)
		node.ReceiveVote()
	}) {
		return
	}

	n.CandidateRequestVote(term)

	select {
	case r := <-n.ElectionResultChannel:
		if !r {
			log.Printf("BeginElection term %d, I losed", term)
		} else {
			f := func(node *Node) {
				node.PromoteToLeader()
			}
			if n.EnsureAndDo(term, CANDIDATE, f) {
				log.Printf("BeginElection term %d I win, annouce my leadership", term)
				n.LeaderAppendEntries(term)
			} else {
				log.Printf("BeginElection inconsistent state, state %d, term %d",
					n.State, n.CurrentTerm)
			}
		}
	}
}

func (n *Node) CandidateRequestVote(term int32) {
	log.Printf("CandidateRequestVote start sending requests")
	f := func(node *Node) {
		req := &raft.RequestVoteRequest{
			Term:         node.CurrentTerm,
			CandidateId:  node.ID,
			LastLogIndex: 0,
			LastLogTerm:  0,
		}
		for id := range node.Clients {
			log.Printf("CandidateRequestVote to %d", id)
			go n.RequestVoteToSingle(id, term, req)
		}
	}
	if !n.EnsureAndDo(term, CANDIDATE, f) {
		log.Printf("abort CandidateRequestVote, inconsistent state, state %d, CurentTerm %d, expected term %d",
			n.State, n.CurrentTerm, term)
	}
}

func (n *Node) RequestVoteToSingle(id, term int32, req *raft.RequestVoteRequest) {
	tries := 0
	for {
		if !n.EnsureAndDo(term, CANDIDATE) {
			log.Printf("abort CandidateRequestVote, inconsistent state, state %d, CurentTerm %d, expected term %d",
				n.State, n.CurrentTerm, term)
			return
		}

		log.Printf("RequestVoteToSingle to %d, try no.%d", id, tries)
		resp, err := n.Clients[id].C.RequestVote(context.TODO(), req)
		if err != nil {
			log.Printf("RequestVoteToSingle to %d, try no.%d, err: %s", id, tries, err)
		}

		shouldReturn := false
		f := func(node *Node) {
			if resp.GetVoteGranted() {
				log.Printf("RequestVoteToSingle to %d, try no.%d, vote granted", id, tries)
				node.ReceiveVote()
				shouldReturn = true
				return
			} else {
				log.Printf("RequestVoteToSingle to %d, try no.%d, fail", id, tries)
				if resp.GetTerm() > node.CurrentTerm {
					log.Printf("RequestVoteToSingle to %d, try no.%d, target has larger term %d, my term %d",
						id, tries, resp.GetTerm(), node.CurrentTerm)
					node.DemoteToFollower(resp.GetTerm())
					shouldReturn = true
					return
				}
			}
		}
		if !n.EnsureAndDo(term, CANDIDATE, f) {
			shouldReturn = true
		}
		if shouldReturn {
			return
		}

		tries++
		time.Sleep(300 * time.Millisecond)
	}
}

func (n *Node) LeaderHeartbeater() {
	log.Printf("start LeaderHeartbeater")
	for {
		select {
		case <-time.After(n.HeartbeatInterval):
			if n.State == LEADER {
				log.Printf("LeaderHeartbeater send heartbeat")
				go n.LeaderAppendEntries(n.CurrentTerm)
			}
		case <-n.DoneC:
			return
		}
	}
}
