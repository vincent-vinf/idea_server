// @Title  im cmd
// @Description  提供im服务
// @Author  Vincent
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"idea_server/db/mysql"
	"idea_server/internal/im"
	"idea_server/myjwt"
	"idea_server/route"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	serve = im.NewServe()
)

func connectHandle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	t, _ := c.Get(myjwt.IdentityKey)
	user := t.(*myjwt.TokenUserInfo)
	id := user.ID
	node := im.NewNode(id, conn)

	log.Printf("id:%s ws connect", id)

	serve.Bucket(id).Put(node)

	go node.Send()
	go node.Receive()
}

func main() {
	mysql.Close()
	go serve.Run()
	//gin.SetMode(gin.ReleaseMode)
	r := route.New(":8001", false)
	r.AddGetAuthRoute("/ws", connectHandle)
	r.Run()
}
