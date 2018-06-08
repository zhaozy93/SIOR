package raft

import (
	"fmt"
	"global"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func (client *Raft) startLeaderWork() {
	go client.startTTL()
}

func (client *Raft) startTTL() {
	circle := time.Millisecond * time.Duration(client.Property.Heartbeat)
	leaderTtl_timer := time.NewTimer(circle)
	td := time.Duration(time.Duration(client.Property.HeartbeatTimeout) * time.Millisecond)
	for {
		select {
		case <-leaderTtl_timer.C:
			if !client.IsLeader() {
				return
			}
			fmt.Printf("leader worker: send ttl to cluster with term: %d\n", client.GetTerm())
			cnt := 0
			waitingKeys := ""
			client.DataLocker.Lock()
			for key, _ := range client.WaitingData {
				waitingKeys = waitingKeys + "-" + key
			}
			justaddKeys := ""
			for key, value := range client.JustADD {
				justaddKeys = justaddKeys + "-" + key + ":" + value
			}
			client.DataLocker.Unlock()
			if len(waitingKeys) > 1 {
				waitingKeys = string(waitingKeys[1:])
			} else {
				waitingKeys = ""
			}
			if len(justaddKeys) > 1 {
				justaddKeys = string(justaddKeys[1:])
			} else {
				justaddKeys = ""
			}
			for _, v := range client.Clusters {
				if strings.Index(v, global.FinalPort) != -1 {
					continue
				}
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
				url := fmt.Sprintf("http://%s/ttl?term=%d&leader=%s&waitingKeys=%s&justaddKeys=%s", v, client.GetTerm(), client.GetHost(), waitingKeys, justaddKeys)
				fmt.Println(url)
				resp, err := httpClient.Get(url)
				if err != nil {
					continue
				} else {
					robots, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						continue
					}
					fmt.Println("resp:", string(robots))
					if string(robots) == "true" {
						cnt = cnt + 1
					}
				}
			}
			if float64(cnt)/float64(len(client.Clusters)) > 0.5 {
				waitingKeysSlice := strings.Split(waitingKeys, "-")
				client.DataLocker.Lock()
				for _, key := range waitingKeysSlice {
					if value, ok := client.WaitingData[key]; ok {
						delete(client.WaitingData, key)
						client.JustADD[key] = value
					}
				}
				justaddKeysSlice := strings.Split(justaddKeys, "-")
				for _, kvstr := range justaddKeysSlice {
					kvslice := strings.Split(kvstr, ":")
					if value, ok := client.JustADD[kvslice[0]]; ok {
						delete(client.JustADD, kvslice[0])
						client.Data[kvslice[0]] = value
					}
				}
				client.DataLocker.Unlock()
			} else {
				// 心跳回复一直失败，那么就失败
				// waitingdata不允许写进去
			}
			leaderTtl_timer.Reset(circle)
		}
	}
}

func (client *Raft) SetKey(key, value string) bool {
	client.DataLocker.Lock()
	client.WaitingData[key] = value
	fmt.Println("add", client.WaitingData[key])
	client.DataLocker.Unlock()
	for {
		if client.IsLeader() {
			if _, ok := client.WaitingData[key]; !ok {
				return true
			}
		} else {
			return false
		}
		time.Sleep(time.Millisecond * 400)
	}

}
