package initialize

import (
	"github.com/gin-gonic/gin"
	"idea_server/middleware"
	"idea_server/router"
)

// Routers 初始化总路由
func Routers() *gin.Engine {
	var Router = gin.Default()

	// 跨域
	Router.Use(middleware.Cors())
	// 获取路由组实例
	userRouter := router.RouterGroupApp.User

	// 无需鉴权
	PublicGroup := Router.Group("")
	{
		userRouter.InitUserBaseRouter(PublicGroup)
	}

	// 需鉴权
	//PrivateGroup := Router.Group("")

	// TODO install plugin see gva

	return Router
}