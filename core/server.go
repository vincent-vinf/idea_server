package core

import (
	"fmt"
	"idea_server/global"
	"idea_server/initialize"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.IDEA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	//global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Println("server run success on", address)
	//global.IDEA_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Println(s.ListenAndServe().Error())
	//global.IDEA_LOG.Error(s.ListenAndServe().Error())
}
