package request

import (
	"idea_server/model/common/request"
	"idea_server/model/idea"
)

type SearchIdeaParams struct {
	idea.Idea
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

