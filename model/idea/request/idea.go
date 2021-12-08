package request

import (
	"idea_server/model/common/request"
	"idea_server/model/idea"
)

type GetByIdeaId struct {
	IdeaId float64 `json:"ideaId"`
}

func (r *GetByIdeaId) Uint() uint {
	return uint(r.IdeaId)
}

type GetSimilarIdeasReq struct {
	Text string
}

type SearchIdeaParams struct {
	idea.Idea
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
