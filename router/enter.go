package router

import (
	"idea_server/router/idea"
	"idea_server/router/user"
)

type RouterGroup struct {
	User user.RouterGroup
	Idea idea.IdeaRouter
}

var RouterGroupApp = new(RouterGroup)
