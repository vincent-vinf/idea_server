package initialize

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"idea_server/global"
	"idea_server/model/user/request"
	"idea_server/service"
	"idea_server/utils"
	"idea_server/utils/constant"
	"time"
)

// TokenUserInfo 结构体中的数据将会编码进 token
type TokenUserInfo struct {
	ID string
}

func JWTAuth() (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	jwtCfg := global.IDEA_CONFIG.JWT
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       constant.AppRealm,
		Key:         []byte(jwtCfg.SigningKey),
		Timeout:     time.Duration(jwtCfg.Timeout),
		MaxRefresh:  time.Duration(jwtCfg.MaxRefresh),
		IdentityKey: constant.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*TokenUserInfo); ok {
				return jwt.MapClaims{
					constant.IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &TokenUserInfo{
				ID: claims[constant.IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginInfo request.Login
			if err := c.ShouldBind(&loginInfo); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginInfo.Email
			password := loginInfo.Passwd
			code := loginInfo.Code

			if !utils.IsEmail(email) {
				return nil, jwt.ErrFailedAuthentication
			}

			var userBaseService = service.ServiceGroupApp.UserServiceGroup.UserBaseService

			if (code != "" && userBaseService.IsCorrectEmailCode(email, code)) || (password != "" && userBaseService.Login(email, password)) {
				id := userBaseService.GetID(email)
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
		return nil, errors.New("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		return nil, errors.New("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	return authMiddleware, nil
}
