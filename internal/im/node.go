package im

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"idea_server/protolcol/message"
	"log"
)

const (
	channelSize = 10
)

// Node only send/receive chat data
type Node struct {
	conn *websocket.Conn
	ch   chan *message.Msg

	Pre  *Node
	Next *Node
}

func NewNode(conn *websocket.Conn) *Node {
	return &Node{
		conn: conn,
		ch:   make(chan *message.Msg, channelSize),
		Pre:  nil,
		Next: nil,
	}
}

func (n *Node) Send() {
	for {
		msg := <-n.ch
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
		msg := &message.Msg{}
		err = proto.Unmarshal(bytes, msg)
		if err != nil {
			log.Println(err)
			continue
		}

	}
}

func (n *Node) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
