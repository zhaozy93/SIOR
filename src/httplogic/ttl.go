package httplogic

import (
	"fmt"
	"helper/context"
	"net/http"
	"raft"
	"strconv"
	"strings"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func ReceivettlHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
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
	leader := r.FormValue("leader")
	if leader == "" {
		fmt.Fprint(w, "false")
		return 0
	}
	addKV := r.FormValue("justaddKeys")
	if len(addKV) > 0 {
		fmt.Println("receive ", addKV)
		justaddKeysSlice := strings.Split(addKV, "-")
		for _, kvstr := range justaddKeysSlice {
			kvslice := strings.Split(kvstr, ":")
			client.WriteData(kvslice[0], kvslice[1])
		}
	}
	fmt.Fprint(w, "true")
	client.HeartbeatChan <- &raft.TTLer{term, leader}
	return 0
}
