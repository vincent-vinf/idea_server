package router

import "idea_server/router/user"

type RouterGroup struct {
	User user.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
