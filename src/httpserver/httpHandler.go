package httpserver

import (
	"fmt"
	logger "github.com/shengkehua/xlog4go"
	"helper/context"
	"net/http"
	"runtime/debug"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, c *context.Context) int

type HttpRequestHandler struct {
	Name     string
	CallFunc HandlerFunc
}

func NewHttpRequestHandler(name string, callFunc HandlerFunc) *HttpRequestHandler {
	return &HttpRequestHandler{
		Name:     name,
		CallFunc: callFunc,
	}
}

func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := context.NewContext(h.Name)
	logid := context.GetLogid()
	defer HandleError(h.Name, logid)

	status := h.CallFunc(w, r, context)
	if status != 0 {
		logger.Error("Logid: %d is error", logid)
	}
}
func HandleError(funcname string, logid int64) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprintf("funcname:%s is err, reason:%v", funcname, err)
		logger.Info("HandleError# recover.......errMsg:%s, logId:%d, stack:%s", errMsg, logid, string(debug.Stack()))
	}
}
