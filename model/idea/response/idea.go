package response

import (
	"idea_server/model/common/response"
	"idea_server/model/idea"
)

type IdeaListResponse struct {
	idea.Idea
	IsLike bool `json:"isLike"`
}

type IdeaInfoResponse struct {
	idea.Idea
	Comments response.PageResult `json:"comments"`
}