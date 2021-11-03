package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	"idea_server/model/user/request"
)

type UserBaseApi struct {
}

func (u *UserBaseApi) Register(c *gin.Context) {
	var msg request.Register
	_ = c.ShouldBindJSON(&msg)
	// TODO 请求体字段验证（是否为空、强密码）

	ok, err := userBaseService.Register(msg)
	if ok {
		response.OkWithMessage("注册成功", c)
	} else {
		global.IDEA_LOG.Error("注册失败", zap.Error(err))
		response.FailWithMessage("注册失败："+err.Error(), c)
	}
}

func (u *UserBaseApi) Login(c *gin.Context) {

}

// GetEmailCode 生成邮箱验证码
func (u *UserBaseApi) GetEmailCode(c *gin.Context) {
	fmt.Println("parma", c.Param("email"))
	fmt.Println("keys", c.Keys)
	//if code, err := userBaseService.GetEmailCode(c.PostForm("email"), c.ClientIP()); err != nil {
	//	global.IDEA_LOG.Error("生成邮箱验证码失败", zap.Error(err))
	//	response.FailWithMessage("生成邮箱验证码失败："+err.Error(), c)
	//} else {
	//	response.OkWithData(gin.H{"email_code": code}, c)
	//}

}
