package raft

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var raftClient *Raft

const (
	leader     = 1
	candidater = 2
	follower   = 3
)

type Raft struct {
	Property      *RaftProperty
	Status        int
	HeartbeatChan chan bool
	Term          uint64
	Clusters      []string // 集群中的iplist
	LastVote      *VoteProperty
	LastTtl       *TtlProperty
	IsTry2Leader  bool
	locker        sync.Mutex
}

type RaftProperty struct {
	Heartbeat        int //Leader的心跳的频率  毫秒数
	HeartbeatTimeout int //Leader的心跳的反馈timeout  毫秒数
	ElectionTimeout  int //Foller收不到心跳后变为Candidate 毫秒数  此数值应该大于Heartbeat
}

type VoteProperty struct {
	Term uint64
	Time int64
}
type TtlProperty struct {
	Term uint64
	Time int64
}

func (client *Raft) IsFollower() bool {
	client.locker.Lock()
	defer client.locker.Unlock()
	if client.Status == follower {
		return true
	}
	return false
}

func (client *Raft) IsLeader() bool {
	client.locker.Lock()
	defer client.locker.Unlock()
	if client.Status == leader {
		return true
	}
	return false
}
func (client *Raft) IsCandidater() bool {
	client.locker.Lock()
	defer client.locker.Unlock()
	if client.Status == candidater {
		return true
	}
	return false
}

func (client *Raft) HasVote(term uint64) bool {
	if client.IsLeader() {
		fmt.Println("leader not vote")
		return true
	}
	lastTime := client.LastVote.Time
	nowTime := time.Now().UnixNano()
	lastTtlTime := client.LastTtl.Time
	// if nowTime-lastTime > 5000 {
	// 	fmt.Println("leader not vote")
	// }
	// if nowTime-lastTtlTime > 100000 {

	// }
	if nowTime-lastTime > 500000 && nowTime-lastTtlTime > 100000 {
		client.LastVote = &VoteProperty{
			Term: term,
			Time: time.Now().UnixNano(),
		}
		return false
	} else {
		return true
	}

}

func (client *Raft) AddTerm() {
	atomic.AddUint64(&client.Term, uint64(1))
}
func (client *Raft) GetTerm() uint64 {
	return client.Term
}

// candidate --> follower
// leader不能变为follower
// 并且变为follower之后还应该开始变为candidate的检测
func (client *Raft) switch2Follower() {
	if client.IsLeader() || client.IsFollower() {
		return
	}
	client.locker.Lock()
	client.Status = follower
	defer client.locker.Unlock()
	go raftClient.candidateChecker()
	return
}

// follower --> candidate
func (client *Raft) switch2Candidate() {
	if client.IsFollower() {
		client.locker.Lock()
		client.Status = candidater
		defer client.locker.Unlock()
	}
}

// candidate --> leader
func (client *Raft) switch2Leader() {
	if !client.IsCandidater() {
		return
	}
	client.locker.Lock()
	client.Status = leader
	defer client.locker.Unlock()
	return
}
