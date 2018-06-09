package httplogic

import (
	"fmt"
	"helper/context"
	"net/http"
	"raft"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func GetKeyHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "false")
		return 0
	}
	key := r.FormValue("key")
	value, ok := client.GetKey(key)
	fmt.Fprint(w, fmt.Sprintf("{err: %t, value: %s}", ok, value))
	return 0
}
