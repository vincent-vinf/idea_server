// @Title  user cmd
// @Description  提供用户服务，包括登陆注册，获取用户信息，刷新token等
// @Author  Vincent

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"idea_server/db/mysql"
	"idea_server/db/redisdb"
	"idea_server/myjwt"
	"idea_server/route"
	"idea_server/util"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	emailRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	defer mysql.Close()
	defer redisdb.Close()
	//gin.SetMode(gin.ReleaseMode)
	r := route.New(":8000", false)
	r.AddPostRoute("/register", registerHandler)
	r.AddGetRoute("/email/code", emailCodeHandler)
	r.AddGetAuthRoute("/userinfo", userinfoHandler)

	r.Run()
}

func registerHandler(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	code := c.PostForm("code")
	if strings.TrimSpace(username) == "" || email == "" || password == "" || code == "" || !util.IsEmail(email) {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}

	if !util.IsStrongPasswd(password) {
		c.JSON(400, gin.H{
			"error": "Weak password",
		})
		return
	}

	if !redisdb.IsCorrectEmailCode(email, code) {
		c.JSON(400, gin.H{
			"error": "The verification code does not exist or has expired",
		})
		return
	}

	isExist, err := mysql.IsExistEmail(email)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Server internal error",
		})
		return
	}
	if isExist {
		c.JSON(400, gin.H{
			"error": "ID already exists",
		})
		return
	}

	err = mysql.Register(username, email, password)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Server internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"email": "ok",
	})
}

func emailCodeHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" || !util.IsEmail(email) {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}
	if !redisdb.IsAllowedIP(c.ClientIP()) {
		c.JSON(403, gin.H{
			"error": "Request too frequent",
		})
		return
	}
	// 生成验证码
	code := fmt.Sprintf("%06v", emailRand.Int31n(1000000))

	err := util.SendMail(email, "Idea email verification code", "Your email verification code:"+code)
	if err != nil {
		c.JSON(422, gin.H{
			"error": "Failed to send mail",
		})
		return
	}
	// 插入验证码到redis
	redisdb.InsertEmailCode(email, code, c.ClientIP())
	log.Println("code: ", code)
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func userinfoHandler(c *gin.Context) {
	t, _ := c.Get(myjwt.IdentityKey)
	user := t.(*myjwt.TokenUserInfo)

	c.JSON(200, gin.H{
		"id": user.ID,
	})
}
