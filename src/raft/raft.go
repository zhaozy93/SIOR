package raft

import (
	"fmt"
	"global"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

func NewRaftClient() *Raft {
	rand.Seed(time.Now().UnixNano())
	randomTime := int(rand.Int31n(int32(global.Cfg.Raft.ElectionTimeoutMax - global.Cfg.Raft.ElectionTimeoutMin)))
	fmt.Println("random", randomTime)
	property := &RaftProperty{
		Heartbeat:        global.Cfg.Raft.Heartbeat,
		HeartbeatTimeout: global.Cfg.Raft.HeartbeatTimeout,
		ElectionTimeout:  randomTime + global.Cfg.Raft.ElectionTimeoutMin,
	}
	client := &Raft{
		Property:      property,
		Status:        follower,
		HeartbeatChan: make(chan bool),
		Term:          uint64(0),
		Clusters:      global.Cfg.Hosts.Cluster,
		LastTtl: &TtlProperty{
			Time: 0,
		},
		LastVote: &VoteProperty{
			Time: 0,
		},
	}
	return client
}

func InitRaftClient() {
	raftClient = NewRaftClient()
	go raftClient.candidateChecker()
}

func GetRaftClient() *Raft {
	return raftClient
}

// 定时检测
// 收到来自于leader的ttl就重置时间
// 一段时间收不到ttl则尝试变成candidate
func (client *Raft) candidateChecker() {
	circle := time.Millisecond * time.Duration(client.Property.ElectionTimeout)
	try2Candidate_timer := time.NewTimer(circle)
	for {
		select {
		// follower状态，长时间未收到ttl，则尝试变为Candidate
		case <-try2Candidate_timer.C:
			try2Candidate_timer.Reset(circle)
			// 只有Follower允许尝试变Candidate
			if !client.IsFollower() || client.IsTry2Leader {
				break
			}
			fmt.Println("try to swicth candidate")
			client.switch2Candidate()
			client.AddTerm()
			go client.try2Leader()
		case <-client.HeartbeatChan:
			try2Candidate_timer.Reset(circle)
			if client.IsLeader() {
				break
			}
			client.LastTtl = &TtlProperty{
				Time: time.Now().UnixNano(),
			}
			fmt.Println("receive ttl from leader, still stay in follow status")
			client.switch2Follower()
		}
	}
}

func (client *Raft) try2Leader() {
	// 正在进行vote或者不是Candidate不能进行try2Leader
	if client.IsTry2Leader || !client.IsCandidater() {
		return
	}
	client.IsTry2Leader = true
	td := time.Duration(time.Duration(500) * time.Millisecond)
	voteCnt := 0
	for _, v := range client.Clusters {
		httpClient := &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, td)
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(td))
					return conn, nil
				},
				ResponseHeaderTimeout: td,
			},
		}
		url := fmt.Sprintf("http://%s/getVote?term=%d", v, client.GetTerm())
		resp, err := httpClient.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, berr := ioutil.ReadAll(resp.Body)
		if berr != nil {
			return
		}
		result := string(body)
		if result == "true" {
			voteCnt = voteCnt + 1
		}
	}
	if float64(voteCnt)/float64(len(client.Clusters)) > 0.5 {
		client.switch2Leader()
		go client.startLeaderWork()
		client.IsTry2Leader = false
		return
	} else {
		fmt.Println("still try to candidate")
		// time.Sleep(time.Second * 1)
		client.IsTry2Leader = false
		client.try2Leader()
	}
}

func (client *Raft) startLeaderWork() {
	circle := time.Millisecond * time.Duration(client.Property.Heartbeat)
	leaderTtl_timer := time.NewTimer(circle)
	td := time.Duration(time.Duration(client.Property.HeartbeatTimeout) * time.Millisecond)
	for {
		select {
		case <-leaderTtl_timer.C:
			if !client.IsLeader() {
				return
			}
			fmt.Println("leader worker: send ttl to cluster")
			for _, v := range client.Clusters {
				fmt.Println(v, global.FinalPort)
				if strings.Index(v, global.FinalPort) != -1 {
					continue
				}
				// fmt.Printf("send ttl to %s\n", v)
				httpClient := &http.Client{
					Transport: &http.Transport{
						Dial: func(netw, addr string) (net.Conn, error) {
							conn, err := net.DialTimeout(netw, addr, td)
							if err != nil {
								return nil, err
							}
							conn.SetDeadline(time.Now().Add(td))
							return conn, nil
						},
						ResponseHeaderTimeout: td,
					},
				}
				url := fmt.Sprintf("http://%s/ttl?term=%d", v, client.GetTerm())
				httpClient.Get(url)
			}
			leaderTtl_timer.Reset(circle)
		}
	}
}
