package httplogic

import (
	"fmt"
	// logger "github.com/shengkehua/xlog4go"
	"helper/context"
	"net/http"
	"raft"
	"strconv"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func ReceivettlHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	fmt.Fprint(w, "ok")
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
		// fmt.Println("no term, vote false。 term: null")
		return 0
	}
	rawValue := r.FormValue("term")
	term, err := strconv.ParseUint(rawValue, 10, 64)
	if err != nil {
		fmt.Fprint(w, "false")
		// fmt.Println("term error, vote false。 term: null")
		return 0
	}
	client.HeartbeatChan <- term
	// fmt.Printf("%v", client)
	return 0
}
