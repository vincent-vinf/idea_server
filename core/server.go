package core

import (
	"fmt"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/initialize"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	// 初始化 redis 服务
	initialize.Redis()
	// 定时任务
	go initialize.Cron()

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.IDEA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.IDEA_LOG.Info("server run success on ", zap.String("address", address))

	global.IDEA_LOG.Error(s.ListenAndServe().Error())
}
