package response

import (
	"idea_server/model/idea"
)

type IdeaListResponse struct {
	idea.Idea
	IsLike bool `json:"isLike"`
	LikeCount int64 `json:"likeCount"`
}

type IdeaCommentResponse struct {
	idea.IdeaComment
	Replys []idea.IdeaComment `json:"replys"`
}

type SimilarIdea struct {
	idea.Idea
	Similarity float64 `json:"similarity"`
}

type IdeaInfoResponse struct {
	//idea.Idea
	Comments []IdeaCommentResponse `json:"comments"`
	IsLike bool `json:"isLike"`
}