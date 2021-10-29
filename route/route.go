package route

import (
	"github.com/gin-gonic/gin"
	"idea_server/myjwt"
	"log"
)

type Route struct {
	router    *gin.Engine
	authRoute *gin.RouterGroup
}

func (r *Route) Run(port string) {
	log.Fatal(r.router.Run(port))
}

func New() *Route {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// 初始化JWT中间件
	authMiddleware, err := myjwt.Init()
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/login", authMiddleware.LoginHandler)

	// 一组需要验证的路由
	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	return &Route{
		router:    router,
		authRoute: auth,
	}
}

func (r *Route) AddGetRoute(path string, handlers ...gin.HandlerFunc) {
	r.router.GET(path, handlers...)
}

func (r *Route) AddPostRoute(path string, handlers ...gin.HandlerFunc) {
	r.router.POST(path, handlers...)
}

func (r *Route) AddAuthGetRoute(path string, handlers ...gin.HandlerFunc) {
	r.authRoute.GET(path, handlers...)
}

func (r *Route) AddAuthPostRoute(path string, handlers ...gin.HandlerFunc) {
	r.authRoute.POST(path, handlers...)
}
