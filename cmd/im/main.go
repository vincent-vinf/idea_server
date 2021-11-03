// @Title  im cmd
// @Description  提供im服务
// @Author  Vincent
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"idea_server/internal/im"
	"idea_server/route"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func connectHandle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	node := im.NewNode(conn)

}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := route.New(":8001", false)
	r.AddGetAuthRoute("/im", connectHandle)
	r.Run()
}
