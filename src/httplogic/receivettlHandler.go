package httplogic

import (
	"fmt"
	logger "github.com/shengkehua/xlog4go"
	"helper/context"
	"net/http"
	"raft"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func ReceivettlHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	fmt.Println("收到leader ttl")
	logger.Info("收到leader ttl")
	fmt.Fprint(w, "ok")
	client.HeartbeatChan <- true
	// fmt.Printf("%v", client)
	return 0
}
