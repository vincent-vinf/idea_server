package im

import (
	"github.com/gorilla/websocket"
	"idea_server/mq"
	"idea_server/protolcol/message"
	"log"
	"strconv"
)

const (
	channelSize = 10
)

// Node only send/receive chat data
type Node struct {
	Id   string
	conn *websocket.Conn
	Ch   chan *message.Msg

	bucket *Bucket

	Pre  *Node
	Next *Node
}

func NewNode(id string, conn *websocket.Conn) *Node {
	return &Node{
		Id:   id,
		conn: conn,
		Ch:   make(chan *message.Msg, channelSize),
		Pre:  nil,
		Next: nil,
	}
}

func (n *Node) Send() {
	for {
		msg := <-n.Ch
		if msg == nil {
			log.Println("node channel error")
			return
		}
		bytes, err := msg.Marshal()
		if err != nil {
			log.Println(err)
			continue
		}
		err = n.conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (n *Node) Receive() {
	for {
		_, bytes, err := n.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		err = mq.GetInstance().Product(bytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (n *Node) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

func int2str(in int32) string {
	return strconv.FormatInt(int64(in), 10)
}
