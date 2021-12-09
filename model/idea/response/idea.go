package response

import (
	"idea_server/model/idea"
)

type IdeaListResponse struct {
	idea.Idea
	TypeName  string `json:"typeName"`
	IsLike    bool   `json:"isLike"`
	LikeCount int64  `json:"likeCount"`
}

type IdeaCommentResponse struct {
	idea.IdeaComment
	Replys []idea.IdeaComment `json:"replys"`
}

type SimilarIdea struct {
	idea.Idea
	Similarity float64 `json:"similarity"`
	TypeName  string `json:"typeName"`
}

type IdeaInfoResponse struct {
	//idea.Idea
	Comments []IdeaCommentResponse `json:"comments"`
	IsLike   bool                  `json:"isLike"`
}

type SimilarModelResponse struct {
	IdeaId uint `json:"ideaId"`
	Similarity float64 `json:"similarity"`
}