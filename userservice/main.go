package main

import (
	"github.com/gin-gonic/gin"
	"idea_server/redisdb"
	"idea_server/util"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	redisdb.ExampleClient()
	//r := route.New()
	//r.AddPostRoute("/register", registerHandler)
	//r.AddAuthGetRoute("/device", devicesHandler)
	//r.Run(":8000")
}

func registerHandler(c *gin.Context) {
	//email := c.PostForm("email")
	//password := c.PostForm("password")
	code := c.PostForm("password")
	if code == "" {
		c.JSON(400, gin.H{
			"email": "ok",
		})
		return
	}
	c.JSON(200, gin.H{
		"email": "ok",
	})
}

func devicesHandler(c *gin.Context) {
	//claims := jwt.ExtractClaims(c)
	//userid := claims[util.IdentityKey].(string)

	t, _ := c.Get(util.IdentityKey)
	user := t.(*util.User)
	c.JSON(200, gin.H{
		"email": user.Email,
	})
}
