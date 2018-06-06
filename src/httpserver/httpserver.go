package httpserver

import (
	"fmt"
	logger "github.com/shengkehua/xlog4go"
	"global"
	"httplogic"
	"net/http"
	"strconv"
)

type SServer struct {
	port int
	//key: uri, value: handler
	urlHandler map[string]*HttpRequestHandler
}

func NewServer(p int) *SServer {
	return &SServer{
		port:       p,
		urlHandler: make(map[string]*HttpRequestHandler, 0),
	}
}

func (s *SServer) Register(name, uri string, hf HandlerFunc) {
	handler := NewHttpRequestHandler(name, hf)
	s.urlHandler[uri] = handler
}

func (s *SServer) Run() {
	mux := http.NewServeMux()

	addr := ":" + strconv.Itoa(s.port)
	for k, v := range s.urlHandler {
		mux.Handle(k, v)
	}

	logger.Info("run http server with addr:%s", addr)

	for i := 0; i < 5; i++ {
		addr = ":" + strconv.Itoa(s.port+i)
		global.FinalPort = addr
		if err := http.ListenAndServe(addr, mux); err != nil {
			//端口占用导致监听失败，sleep 2秒重试一次
			// logger.Error("run http server fail:%s", err.Error())
			fmt.Printf("listening port error %s\n", addr)
		} else {
			fmt.Printf("listening port %s\n", addr)
			return
		}
	}
	// if err := http.ListenAndServe(addr, mux); err != nil {
	// 	//端口占用导致监听失败，sleep 2秒重试一次
	// 	logger.Error("run http server fail:%s", err.Error())
	// }
}

func RunHttpServer() {
	s := NewServer(global.Cfg.Service.Port)

	// s.Register(global.IF_ASYNCBYAREA, GetDriverByAreaAsyncHandler)

	// s.Register("testHttp", "/test", test)
	s.Register("receiveHeartbeatChan", "/ttl", httplogic.ReceivettlHandler)
	s.Register("vote", "/getVote", httplogic.VoteHandler)

	s.Run()
}
