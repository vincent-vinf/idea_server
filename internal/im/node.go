package im

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
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

func NewNode(id string, conn *websocket.Conn, bucket *Bucket) *Node {
	return &Node{
		Id:     id,
		conn:   conn,
		Ch:     make(chan *message.Msg, channelSize),
		bucket: bucket,
		Pre:    nil,
		Next:   nil,
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
		msg := &message.Msg{}
		err = proto.Unmarshal(bytes, msg)
		if err != nil {
			log.Println(err)
			continue
		}

		if int2str(msg.Uid) != n.Id {
			log.Println("id?")
			continue
		}

		switch msg.Op {
		case message.SendMsg:
			dispatch(msg, n.bucket)
		}
	}
}

func (n *Node) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

func dispatch(msg *message.Msg, bucket *Bucket) {
	if msg.IsGroup {

	} else {
		log.Println(int2str(msg.Did))
		log.Println(msg)
		node := bucket.Node(strconv.FormatInt(int64(msg.Did), 10))
		if node != nil {
			node.Ch <- msg
		}
	}
}

func int2str(in int32) string {
	return strconv.FormatInt(int64(in), 10)
}
