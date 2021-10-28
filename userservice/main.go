package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"idea_server/redisdb"
	"idea_server/route"
	"idea_server/util"
	"log"
	"math/rand"
	"time"
)

const (
	codeExpiration = time.Minute * 10 // 验证码有效期
	codeForbidden  = time.Minute      // 验证码重复发送时间
)

var (
	emailRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := route.New()
	r.AddPostRoute("/register", registerHandler)
	r.AddGetRoute("/email/code", emailCodeHandler)
	r.Run(":8000")
}

func registerHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	code := c.PostForm("code")
	if email == "" || password == "" || code == "" || !util.CheckEmail(email) {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}

	if !isAvailableEmailCode(email, code) {
		c.JSON(400, gin.H{
			"error": "The verification code does not exist or has expired",
		})
		return
	}
	c.JSON(200, gin.H{
		"email": "ok",
	})
}

func emailCodeHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" || !util.CheckEmail(email) {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}
	if !isAllowedIP(c.ClientIP()) {
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
	insertEmailCode(email, code, c.ClientIP())
	log.Println("code: ", code)
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func insertEmailCode(email, code, ip string) {
	rdb := redisdb.GetInstance()
	ctx := context.Background()
	// 插入验证码和请求的ip，防止过多请求
	rdb.Set(ctx, email, code, codeExpiration)
	rdb.Set(ctx, ip, ip, codeForbidden)
}

func isAvailableEmailCode(email, code string) bool {
	rdb := redisdb.GetInstance()
	ctx := context.Background()
	re, err := rdb.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	if re != code {
		return false
	}
	return true
}

func isAllowedIP(ip string) bool {
	rdb := redisdb.GetInstance()
	ctx := context.Background()
	_, err := rdb.Get(ctx, ip).Result()
	if err == redis.Nil {
		return true
	} else {
		return false
	}
}

//func devicesHandler(c *gin.Context) {
//	//claims := jwt.ExtractClaims(c)
//	//userid := claims[util.IdentityKey].(string)
//
//	t, _ := c.Get(util.IdentityKey)
//	user := t.(*util.User)
//	c.JSON(200, gin.H{
//		"email": user.Email,
//	})
//}
