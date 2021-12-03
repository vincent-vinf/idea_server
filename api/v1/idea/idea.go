package idea

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	user "idea_server/model/user/response"
	"idea_server/utils/constant"
)

type IdeaApi struct {
}

func (e *IdeaApi) CreateIdea(c *gin.Context) {
	rawJson := make(map[string]interface{})
	_ = c.ShouldBindJSON(&rawJson)
	userInfo, _ := c.Get(constant.IdentityKey)
	ok, err := ideaService.CreateIdea(userInfo.(*user.UserResponse).ID, rawJson["content"].(string))
	if ok {
		response.OkWithMessage("创建想法成功", c)
	} else {
		global.IDEA_LOG.Error("创建想法失败", zap.Error(err))
		response.FailWithMessage("创建想法失败："+err.Error(), c)
	}
}

func (e *IdeaApi) GetIdeaList(c *gin.Context) {

}

func (e *IdeaApi) GetIdeaInfo(c *gin.Context) {
	var info request.GetById
	_ = c.ShouldBindJSON(&info)
	if ideaInfo, err := ideaService.GetIdeaInfo(&info); err != nil {
		global.IDEA_LOG.Error("查询想法失败", zap.Error(err))
		response.FailWithMessage("查询想法失败："+err.Error(), c)
	} else {
		response.OkWithData(ideaInfo, c)
	}
}

func (e *IdeaApi) UpdateIdea(c *gin.Context) {

}

func (e *IdeaApi) DeleteIdea(c *gin.Context) {

}
