package idea

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	ideaReq "idea_server/model/idea/request"
	userRes "idea_server/model/user/response"
	"idea_server/utils"
	"idea_server/utils/constant"
)

type IdeaApi struct {
}

func (e *IdeaApi) CreateIdea(c *gin.Context) {
	rawJson := make(map[string]interface{})
	_ = c.ShouldBindJSON(&rawJson)
	userInfo, _ := c.Get(constant.IdentityKey)
	ok, err := ideaService.CreateIdea(userInfo.(*userRes.UserResponse).ID, rawJson["content"].(string))
	if ok {
		response.OkWithMessage("创建想法成功", c)
	} else {
		global.IDEA_LOG.Error("创建想法失败", zap.Error(err))
		response.FailWithMessage("创建想法失败", c)
	}
}

func (e *IdeaApi) GetIdeaList(c *gin.Context) {
	var info ideaReq.SearchIdeaParams
	_ = c.ShouldBindJSON(&info)

	if err := utils.Verify(info.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err, list, total := ideaService.GetIdeaList(info.Idea, info.PageInfo, info.OrderKey, info.Desc); err != nil {
		global.IDEA_LOG.Error("获取想法列表失败!", zap.Error(err))
		response.FailWithMessage("获取想法列表失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取成功", c)
	}
}

func (e *IdeaApi) GetIdeaInfo(c *gin.Context) {
	var info request.GetById
	_ = c.ShouldBindJSON(&info)
	if ideaInfo, err := ideaService.GetIdeaInfo(&info); err != nil {
		global.IDEA_LOG.Error("查询想法失败", zap.Error(err))
		response.FailWithMessage("查询想法失败", c)
	} else {
		response.OkWithData(ideaInfo, c)
	}
}

func (e *IdeaApi) UpdateIdea(c *gin.Context) {

}

func (e *IdeaApi) DeleteIdea(c *gin.Context) {

}
