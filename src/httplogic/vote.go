package httplogic

import (
	"fmt"
	"helper/context"
	"net/http"
	"raft"
	"strconv"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func VoteHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	// fmt.Println("receive vote request")
	client := raft.GetRaftClient()
	// 如果他自己就是leader则不允许给其他candidate投票
	// 按理来说这种现象是不会出现的
	if client.IsLeader() {
		fmt.Fprint(w, "false")
		return 0
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
		return 0
	}
	rawValue := r.FormValue("term")
	term, err := strconv.ParseUint(rawValue, 10, 64)
	if err != nil {
		fmt.Fprint(w, "false")
		return 0
	}
	if client.HasVote(term) {
		fmt.Fprint(w, "false")
		return 0
	}
	fmt.Fprint(w, "true")
	// 如果他已经投过票了也不允许投票
	return 0
}
