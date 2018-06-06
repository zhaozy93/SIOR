package httplogic

import (
	"fmt"
	// logger "github.com/shengkehua/xlog4go"
	"helper/context"
	"net/http"
	"raft"
	"strconv"
	"time"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func VoteHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	// fmt.Println("receive vote request")
	client := raft.GetRaftClient()
	// fmt.Println("收到vote请求")
	// 如果他自己就是leader则不允许给其他candidate投票
	if client.IsLeader() {
		fmt.Fprint(w, "false")
		fmt.Println("leader status, vote false。 term: null")
		return 0
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
		fmt.Println("no term, vote false。 term: null")
		return 0
	}
	rawValue := r.FormValue("term")
	term, err := strconv.ParseUint(rawValue, 10, 64)
	if err != nil {
		fmt.Fprint(w, "false")
		fmt.Println("term error, vote false。 term: null")
		return 0
	}
	if client.HasVote(term) {
		fmt.Fprint(w, "false")
		fmt.Printf("just vote, vote false。 term: %d\n", term)
		return 0
	}
	fmt.Fprint(w, "true")
	fmt.Printf("vote true。 term: %d\n", term)
	time.Sleep(2 * time.Second)
	// 如果他已经投过票了也不允许投票
	return 0
}
