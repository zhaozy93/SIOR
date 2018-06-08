package httplogic

import (
	"encoding/json"
	"fmt"
	"helper/context"
	"net/http"
	"raft"
)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

func GetInfoHandler(w http.ResponseWriter, r *http.Request, c *context.Context) int {
	client := raft.GetRaftClient()
	str, err := json.Marshal(client)
	if err != nil {
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, string(str))
	}
	return 0
}
