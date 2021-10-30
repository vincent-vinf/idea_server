// @Title  connect service
// @Description  连接层，保持连接，维持在线状态
// @Author  Vincent
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"idea_server/redisdb"
	"idea_server/route"
	"log"
)

var upgrader = websocket.Upgrader{}

func connectHandle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		//读取ws中的数据
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = conn.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func main() {
	defer redisdb.Close()
	//gin.SetMode(gin.ReleaseMode)
	r := route.New(":8001", true)
	r.AddGetRoute("/connect", connectHandle)
	r.Run()

}
