package im

import (
	"context"
	"github.com/golang/protobuf/proto"
	"idea_server/mq"
	"idea_server/protolcol/message"
	"log"
)

type Serve struct {
	buckets []*Bucket
	mq      *mq.Manager
}

func NewServe() *Serve {
	return &Serve{
		buckets: []*Bucket{NewBucket()},
		mq:      nil,
	}
}

func (s *Serve) Bucket(uid string) *Bucket {
	return s.buckets[0]
}

func (s *Serve) Run() {
	mq.GetInstance().Consume(context.TODO(), func(bytes []byte) {
		msg := &message.Msg{}
		err := proto.Unmarshal(bytes, msg)
		if err != nil {
			log.Println(err)
		}
		did := int2str(msg.Did)
		uid := int2str(msg.Uid)
		bucket := s.Bucket(uid)
		if msg.IsGroup {
			group := bucket.Group(did)
			if group != nil {
				group.SendMsg(msg)
			} else {
				log.Println("error group id")
			}
		} else {
			node := bucket.Node(did)
			if node != nil {
				node.Ch <- msg
			}
		}
	})
}
