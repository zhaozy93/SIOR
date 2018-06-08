package raft

import (
	"fmt"
	"global"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func NewRaftClient() *Raft {
	rand.Seed(time.Now().UnixNano())
	randomTime := int(rand.Int31n(int32(global.Cfg.Raft.ElectionTimeoutMax - global.Cfg.Raft.ElectionTimeoutMin)))
	property := &RaftProperty{
		Heartbeat:        global.Cfg.Raft.Heartbeat,
		HeartbeatTimeout: global.Cfg.Raft.HeartbeatTimeout,
		ElectionTimeout:  randomTime + global.Cfg.Raft.ElectionTimeoutMin,
	}
	client := &Raft{
		Property:      property,
		Status:        follower,
		HeartbeatChan: make(chan *TTLer),
		Term:          uint64(0),
		Clusters:      global.Cfg.Hosts.Cluster,
		Host:          global.IPPort,
		LastTtl: &TtlProperty{
			Time:   0,
			Term:   uint64(0),
			Leader: "",
		},
		LastVote: &VoteProperty{
			Time: 0,
			Term: uint64(0),
		},
		Data:        make(map[string]string),
		WaitingData: make(map[string]string),
		JustADD:     make(map[string]string),
	}
	return client
}

func InitRaftClient() {
	raftClient = NewRaftClient()
	go raftClient.candidateChecker()
	// go raftClient.httpLogic()
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
			// Leader不允许尝试变为Candidate
			if client.IsLeader() || client.IsTry2Leadering {
				break
			}
			fmt.Println("try to swicth candidate")
			client.switch2Candidate()
			client.AddTerm()
			go client.try2Leader()
		case ttler := <-client.HeartbeatChan:
			try2Candidate_timer.Reset(circle)
			if client.IsLeader() {
				break
			}
			client.LastTtl = &TtlProperty{
				Time:   time.Now().UnixNano(),
				Term:   ttler.Term,
				Leader: ttler.Leader,
			}
			fmt.Println("receive ttl from leader, still stay in follow status")
			client.switch2Follower()
			client.UpdateTerm(ttler.Term)
		}
	}
}

// 投票环节
// 拿到多数票就变为Leader 否则继续作为Candidate
func (client *Raft) try2Leader() {
	// 正在进行vote或者不是Candidate不能进行try2Leader
	if client.IsTry2Leadering || !client.IsCandidater() {
		return
	}
	client.IsTry2Leadering = true
	td := time.Duration(time.Duration(global.Cfg.Raft.ElectionVoteTimeout) * time.Millisecond)
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
			// fmt.Println(err.Error())
			continue
		}
		defer resp.Body.Close()
		body, berr := ioutil.ReadAll(resp.Body)
		if berr != nil {
			// fmt.Println(berr.Error())
			continue
		}
		result := string(body)
		if result == "true" {
			voteCnt = voteCnt + 1
		}
	}
	if float64(voteCnt)/float64(len(client.Clusters)) > 0.5 {
		client.switch2Leader()
		go client.startLeaderWork()
		fmt.Printf("swicth to leader success voteCnt: %d of %d\n", voteCnt, len(client.Clusters))
	} else {
		fmt.Printf("still in candidate status voteCnt: %d of %d\n", voteCnt, len(client.Clusters))
	}
	client.IsTry2Leadering = false
}

// 启动一个协程来处理来自http的请求，例如投票
// 解耦http接口
// func (client *Raft) httpLogic() {
// 	for {
// 		select {
// 		case term := <-client.HeartbeatChan:
// 			try2Candidate_timer.Reset(circle)
// 			if client.IsLeader() {
// 				break
// 			}
// 			client.LastTtl = &TtlProperty{
// 				Time: time.Now().UnixNano(),
// 			}
// 			fmt.Println("receive ttl from leader, still stay in follow status")
// 			client.switch2Follower()
// 			client.UpdateTerm(term)
// 		}
// 	}
// }
