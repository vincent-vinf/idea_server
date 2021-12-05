package initialize

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/middleware"
	"idea_server/router"
)

// Routers 初始化总路由
func Routers() *gin.Engine {
	var Router = gin.Default()

	//Router.Use(middleware.LoadTls())
	// 跨域
	Router.Use(middleware.Cors())
	// 获取路由组实例
	userRouter := router.RouterGroupApp.User
	ideaRouter := router.RouterGroupApp.Idea

	// 无需鉴权
	PublicGroup := Router.Group("")
	{
		userRouter.InitUserBaseRouter(PublicGroup)
	}

	// jwt
	authMiddleware, err := JWTAuth()
	if err != nil {
		global.IDEA_LOG.Error("jwt 初始化失败", zap.Error(err))
		panic(err)
	}
	// 404
	Router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		global.IDEA_LOG.Info(fmt.Sprintf("NoRoute claims: %#v\n", claims))
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	// 登录
	PublicGroup.POST("/login", authMiddleware.LoginHandler)
	// 刷新 token
	PublicGroup.GET("/refresh_token", authMiddleware.RefreshHandler)

	// 需鉴权
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(authMiddleware.MiddlewareFunc()).Use(middleware.LimitLogin())
	{
		userRouter.InitUserRouter(PrivateGroup)
		ideaRouter.InitIdeaRouter(PrivateGroup)
		ideaRouter.InitIdeaCommentRouter(PrivateGroup)
	}
	// TODO install plugin see gva

	return Router
}
