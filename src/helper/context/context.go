package context

import (
	"sync/atomic"
	"time"
)

var ii = time.Now().UnixNano()
var i = &ii

type Context struct {
	Caller string
	LogId  int64
}

func NewContext(caller string) *Context {
	return &Context{
		Caller: caller,
		LogId:  atomic.AddInt64((*int64)(i), 1),
	}
}

func (ctx *Context) GetLogid() int64 {
	return ctx.LogId
}
