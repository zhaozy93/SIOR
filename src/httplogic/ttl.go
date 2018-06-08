package httplogic

import (
	"fmt"
	// logger "github.com/shengkehua/xlog4go"
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
	leader := r.FormValue("leader")
	if leader == "" {
		fmt.Fprint(w, "false")
		// fmt.Println("term error, vote false。 term: null")
		return 0
	}
	addKV := r.FormValue("justaddKeys")
	if len(addKV) > 0 {
		fmt.Println("receive ", addKV)
		client.DataLocker.Lock()
		justaddKeysSlice := strings.Split(addKV, "-")
		for _, kvstr := range justaddKeysSlice {
			kvslice := strings.Split(kvstr, ":")
			client.Data[kvslice[0]] = kvslice[1]
		}
		client.DataLocker.Unlock()
	}
	fmt.Fprint(w, "true")
	client.HeartbeatChan <- &raft.TTLer{term, leader}
	// fmt.Printf("%v", client)
	return 0
}
