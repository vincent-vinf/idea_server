package user

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	"idea_server/model/user/request"
	"idea_server/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserBaseApi struct {
}

func (u *UserBaseApi) Register(c *gin.Context) {
	var msg request.Register
	_ = c.ShouldBindJSON(&msg)
	if err := utils.Verify(msg, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ok, err := userBaseService.Register(msg)
	if ok {
		loginData := "{\"email\": \"" + msg.Email + "\", \"passwd\": \"" + msg.Passwd + "\"}"
		res, _ := http.Post("http://127.0.0.1:"+strconv.Itoa(global.IDEA_CONFIG.System.Addr)+"/login", "application/json", bytes.NewBuffer([]byte(loginData)))
		data, _ := ioutil.ReadAll(res.Body)
		m := make(map[string]interface{})
		_ = json.Unmarshal(data, &m)
		response.OkWithData(m, c)
	} else {
		global.IDEA_LOG.Error("注册失败", zap.Error(err))
		response.FailWithMessage("注册失败", c)
	}
}

// GetEmailCode 生成邮箱验证码
func (u *UserBaseApi) GetEmailCode(c *gin.Context) {
	rawJson := make(map[string]interface{})
	_ = c.ShouldBindJSON(&rawJson)
	if code, err := userBaseService.GetEmailCode(rawJson["email"].(string), c.ClientIP()); err != nil {
		global.IDEA_LOG.Error("生成邮箱验证码失败", zap.Error(err))
		response.FailWithMessage("生成邮箱验证码失败", c)
	} else {
		response.OkWithData(gin.H{"email_code": code}, c)
	}

}
