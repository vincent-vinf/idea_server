package initialize

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/service"
	"os"
)

func Cron() (err error) {
	c := cron.New(cron.WithSeconds())

	spec := "00 00 */1 * * ?"
	_, err = c.AddFunc(spec, service.ServiceGroupApp.IdeaServiceGroup.LifeCronFunc)
	if err != nil {
		global.IDEA_LOG.Error("启动定时任务失败", zap.Error(err))
		os.Exit(1)
	}

	c.Start()
	select {}
}