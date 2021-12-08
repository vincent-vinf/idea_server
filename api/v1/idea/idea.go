package idea

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	ideaReq "idea_server/model/idea/request"
	"idea_server/utils"
)

type IdeaApi struct {
}

func (e *IdeaApi) CreateIdea(c *gin.Context) {
	rawJson := make(map[string]interface{})
	_ = c.ShouldBindJSON(&rawJson)
	ok, err := ideaService.CreateIdea(utils.GetJwtId(c), rawJson["content"].(string))
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

	if err, list, total, num := ideaService.GetIdeaList(info.Idea, info.PageInfo, info.OrderKey, info.Desc, utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("获取想法列表失败!", zap.Error(err))
		response.FailWithMessage("获取想法列表失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Num:      num,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取想法列表成功", c)
	}
}

func (e *IdeaApi) GetIdeaInfo(c *gin.Context) {
	var info request.GetById
	_ = c.ShouldBindJSON(&info)
	if ideaInfo, err := ideaService.GetIdeaInfo(&info, utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("查询想法失败", zap.Error(err))
		response.FailWithMessage("查询想法失败", c)
	} else {
		response.OkWithDetailed(ideaInfo, "查询想法成功", c)
	}
}

func (e *IdeaApi) GetSimilarIdeas(c *gin.Context) {
	var info ideaReq.GetSimilarIdeasReq
	_ = c.ShouldBindJSON(&info)
	if list, err := ideaService.GetSimilarIdeasByText(info.Text); err != nil {
		global.IDEA_LOG.Error("获取相似想法失败", zap.Error(err))
		response.FailWithMessage("获取相似想法失败", c)
	} else {
		response.OkWithDetailed(list, "获取相似想法成功", c)
	}
}

func (e *IdeaApi) UpdateIdea(c *gin.Context) {

}

func (e *IdeaApi) DeleteIdea(c *gin.Context) {

}
