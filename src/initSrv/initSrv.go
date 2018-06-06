package initSrv

import (
	gcfg "code.google.com/p/gcfg"
	logger "github.com/shengkehua/xlog4go"
	"global"
)

func InitConf() error {
	return gcfg.ReadFileInto(&global.Cfg, global.CONF_FILE)
}

func InitLogger() error {
	return logger.SetupLogWithConf(global.Cfg.Service.LogFile)
}
