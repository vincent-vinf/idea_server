package request

import (
	"idea_server/model/common/request"
	"idea_server/model/idea"
)

type GetIdeaCommentListReq struct {
	idea.IdeaComment
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
