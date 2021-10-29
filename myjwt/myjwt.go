package myjwt

import (
	"errors"
	"idea_server/db"
	"idea_server/redisdb"
	"idea_server/util"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 登录所需的数据
type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     string `form:"code" json:"code"`
}

// 结构体中的数据将会编码进token
type TokenUserInfo struct {
	ID string
}

const (
	IdentityKey     = "id"
	appRealm        = "idea"
	tokenTimeout    = time.Hour * 1
	tokenMaxRefresh = time.Hour * 2
)

func Init() (*jwt.GinJWTMiddleware, error) {
	cfg := util.LoadJWTCfg()
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       appRealm,
		Key:         cfg.SecretKey,
		Timeout:     tokenTimeout,
		MaxRefresh:  tokenMaxRefresh,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*TokenUserInfo); ok {
				return jwt.MapClaims{
					IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &TokenUserInfo{
				ID: claims[IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginInfo login
			if err := c.ShouldBind(&loginInfo); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginInfo.Email
			password := loginInfo.Password
			code := loginInfo.Code

			if !util.IsEmail(email) {
				return nil, jwt.ErrFailedAuthentication
			}

			if (code != "" && redisdb.IsCorrectEmailCode(email, code)) || (password != "" && db.Login(email, password)) {
				id, err := db.GetID(email)
				if err != nil {
					log.Println(err)
					return nil, jwt.ErrFailedAuthentication
				}
				log.Println(id)
				u := &TokenUserInfo{
					ID: id,
				}
				return u, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})
	if err != nil {
		return nil, err
	}
	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return nil, errors.New("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}
	return authMiddleware, nil
}
