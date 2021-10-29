package route

import (
	"github.com/gin-gonic/gin"
	"idea_server/myjwt"
	"log"

	"github.com/unrolled/secure"
)

type Route struct {
	router    *gin.Engine
	authRoute *gin.RouterGroup
	host      string
	isTLS     bool
}

func (r *Route) Run() {
	if r.isTLS {
		log.Println("use TLS")
		log.Fatal(r.router.RunTLS(r.host, "./cert/idea.vinf.top.cer", "./cert/idea.vinf.top.key"))
	} else {
		log.Fatal(r.router.Run(r.host))
	}
}

func New(host string, isTLS bool) *Route {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if isTLS {
		router.Use(TlsHandler(host))
	}

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
		host:      host,
		isTLS: isTLS,
	}
}

func TlsHandler(host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     host,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
			return
		}
		c.Next()
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
