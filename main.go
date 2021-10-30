package main

import (
	"idea_server/core"
	"idea_server/global"
	"idea_server/initialize"
)

func main() {
	global.IDEA_VP = core.Viper() // 初始化 viper 配置管理
	//global.GVA_LOG = core.Zap()       // 初始化 zap 日志库
	global.IDEA_DB = initialize.Gorm() // 连接 mysql
	if global.IDEA_DB != nil {
		db, _ := global.IDEA_DB.DB()
		defer db.Close()
	}
	core.RunServer()
}
