package main

import (
	"fmt"
	"math/rand"
	"time"
)

var rpc Rpc

type RaftNode struct {
	id        int
	addr      string
	role      Role
	others    map[int]*RaftNode
	term      int
	voteNode  int
	leader    int
	gotVotes  int
	heartbeat chan bool
}

func (node *RaftNode) Start() {
	rpc = new(DummyRpc)
	for {
		switch node.role {
		case Leader:
			node.Heartbeat()
		case Candidate:
			node.term++
			node.voteNode = node.id
			node.gotVotes++
			if node.BeginVote() {
				node.role = Leader
			}
		case Follower:
			select {
			case <-node.heartbeat:
				fmt.Println("get heartbeat")
			case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
				node.role = Candidate
			}
		}
	}
}

func (node *RaftNode) BeginVote() bool {
	req := VoteReq{
		term: node.term,
		id:   node.id,
	}
	// todo 可优化成快速失败
	for e := range node.others {
		res := rpc.CallVote(node.others[e].addr, req)
		if res.vote {
			node.gotVotes++
			if node.gotVotes > len(node.others)/2 {
				return true
			}
		}
		if res.term > node.term {
			node.role = Follower
			node.term = res.term
			node.voteNode = -1
			return false
		}
	}
	return false
}

func (node *RaftNode) HandleVoteReq(req *VoteReq) VoteRes {
	if node.term < req.term && node.voteNode == -1 {
		return VoteRes{
			term: node.term,
			vote: true,
		}
	}
	return VoteRes{
		term: node.term,
		vote: false,
	}
}

func (node *RaftNode) Heartbeat() {
	if node.role != Leader {
		return
	}
	msg := HeartbeatMsg{
		id:   node.id,
		term: node.term,
	}
	for e := range node.others {
		if e != node.id {
			go rpc.CallHeartbeat(node.others[e].addr, msg)
		}
	}
}

func (node *RaftNode) HandleHeartbeat() {
	node.heartbeat <- true
}

type HeartbeatMsg struct {
	id   int
	term int
}

type VoteReq struct {
	term int
	id   int
}

type VoteRes struct {
	term int
	vote bool
}

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)
