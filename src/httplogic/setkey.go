package httplogic

import (
	"fmt"
	"helper/context"
	"io/ioutil"
	"net/http"
	"raft"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

// 这个接口可能会hang住，如果key一直不被大多数follower接受
func SetKeyHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
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
