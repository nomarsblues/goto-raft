package main

import "fmt"

type Rpc interface {
	CallVote(addr string, req VoteReq) VoteRes
	CallHeartbeat(addr string, req HeartbeatMsg)
}

type DummyRpc struct {
}

func (rpc DummyRpc) CallVote(addr string, req VoteReq) VoteRes {
	fmt.Printf("call vote")
	return VoteRes{}
}

func (rpc DummyRpc) CallHeartbeat(addr string, req HeartbeatMsg) {
	fmt.Println("call heartbeat")
}
