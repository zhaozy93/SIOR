package main

/*
Main模块:
   1、配置初始化
   2、Log初始化
   3、启动http服务
*/
import (
	"fmt"
	logger "github.com/shengkehua/xlog4go"
	"global"
	"httpserver"
	"initSrv"
	"os"
	"raft"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//init conf
	if err := initSrv.InitConf(); err != nil {
		fmt.Println(fmt.Sprintf("msg=[service init fail] detail=[init config fail] err=[%s]", err.Error()))
		os.Exit(-1)
	}

	if err := initSrv.InitLogger(); err != nil {
		fmt.Println(fmt.Sprintf("msg=[service init fail] detail=[init log fail] err=[%s]", err.Error()))
		os.Exit(-1)
	}
	defer logger.Close()

	helloStr := fmt.Sprintf("Hello SIOR service %v", global.Version)
	fmt.Println(helloStr)
	logger.Info(helloStr)

	raft.InitRaftClient()
	httpserver.RunHttpServer()

}
