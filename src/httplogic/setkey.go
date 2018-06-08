package httplogic

import (
	"fmt"
	// logger "github.com/shengkehua/xlog4go"
	"helper/context"
	"io/ioutil"
	"net/http"
	"raft"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func SetKeyHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
		// fmt.Println("no term, vote false。 term: null")
		return 0
	}
	key := r.FormValue("key")
	value := r.FormValue("value")
	if key == "" || value == "" {
		fmt.Fprint(w, "false")
		return 0
	}
	if client.IsLeader() {
		result := client.SetKey(key, value)
		if result {
			fmt.Fprint(w, "true")
		} else {
			fmt.Fprint(w, "false")
		}
	} else {
		leaderHost := client.GetLeader()
		url := "http://" + leaderHost + "/setKey?key=" + key + "&value=" + value
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprint(w, "false")
		} else {
			robots, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Fprint(w, "false")
				return 0
			}
			fmt.Fprint(w, string(robots))
		}
	}
	return 0
}
