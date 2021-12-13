package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	"idea_server/utils"
	"idea_server/utils/constant"
)

type UserApi struct {
}

func (u *UserApi) GetMyInfo(c *gin.Context) {
	userInfo, _ := c.Get(constant.IdentityKey)
	if userInfo == nil {
		global.IDEA_LOG.Error("获取个人信息失败")
		response.Fail(c)
		return
	}
	response.OkWithDetailed(userInfo, "获取个人信息成功", c)
}

func (u *UserApi) GetUserInfo(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if info, err := userService.GetUserInfo(ids.Ids, utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
	} else {
		response.OkWithDetailed(info, "获取用户信息成功", c)
	}
}

func (u *UserApi) Notice(c *gin.Context) {
	var info request.PageInfo
	_ = c.ShouldBindJSON(&info)

	if err := utils.Verify(info, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, num := userService.Notice(info, utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("获取用户通知失败", zap.Error(err))
		response.FailWithMessage("获取用户通知失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Num:      num,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取用户通知成功", c)
	}
}
