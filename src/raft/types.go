package raft

import (
	"fmt"
	"global"
	"sync"
	"sync/atomic"
)

var raftClient *Raft

const (
	leader     = 1
	candidater = 2
	follower   = 3
)

type Raft struct {
	Property        *RaftProperty     `json:"property"`
	Status          int               `json:"status"`
	HeartbeatChan   chan *TTLer       `json:"-"`
	Term            uint64            `json:"currentTerm"`
	Clusters        []string          `json:"clusters"` // 集群中的iplist
	LastVote        *VoteProperty     `json:"-"`
	LastTtl         *TtlProperty      `json:"ttlInfo"`
	IsTry2Leadering bool              `json:"-"`
	locker          sync.Mutex        `json:"-"`
	Host            string            `json:"host"`
	Data            map[string]string `json:"data"`
	WaitingData     map[string]string `json:"-"`
	JustADD         map[string]string `json:"-"`
	DataLocker      sync.Mutex        `json:"-"`
}

type RaftProperty struct {
	Heartbeat        int `json:"heartbeatTimeout"` //Leader的心跳的频率  毫秒数
	HeartbeatTimeout int `json:"-"`                //Leader的心跳的反馈timeout  毫秒数
	ElectionTimeout  int `json:"ElectionTimeout"`  //Foller收不到心跳后变为Candidate 毫秒数  此数值应该大于Heartbeat
}

type VoteProperty struct {
	Term uint64
	Time int64
}
type TtlProperty struct {
	Term   uint64 `json:"term"`
	Time   int64  `json:"time"`
	Leader string `json:"leader"`
}

func (client *Raft) UpdateTerm(term uint64) {
	client.Term = term
}

// 角色判定
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

// 每个term只允许投一次票
// 我们记录上次投票的term
// 如果上次投票的term大于这次投票的term那我们就直接返回ture表示已经投过票了
// 投票是不带主信息的 也就是说投票只看term不看其他信息
func (client *Raft) HasVote(term uint64) bool {
	if client.IsLeader() {
		fmt.Println("leader not vote")
		return true
	}
	lastVoteTerm := client.LastVote.Term
	if lastVoteTerm >= term {
		return true
	}
	client.LastVote.Term = term
	client.UpdateTerm(term)
	return false
}

func (client *Raft) AddTerm() {
	atomic.AddUint64(&client.Term, uint64(1))
}
func (client *Raft) GetTerm() uint64 {
	return client.Term
}
func (client *Raft) GetHost() string {
	if client.Host == "" {
		client.Host = global.IPPort
	}
	return client.Host
}
func (client *Raft) GetLeader() string {
	if client.IsLeader() {
		return client.GetHost()
	} else {
		return client.LastTtl.Leader
	}
}
func (client *Raft) GetKey(key string) (string, bool) {
	if val, ok := client.Data[key]; ok {
		return val, true
	}
	return "", false
}
func (client *Raft) WriteData(key, value string) {
	client.DataLocker.Lock()
	client.Data[key] = value
	defer client.DataLocker.Unlock()
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

type TTLer struct {
	Term   uint64
	Leader string
}
