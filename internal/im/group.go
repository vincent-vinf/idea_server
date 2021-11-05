package im

import (
	"idea_server/protolcol/message"
	"sync"
)

type Group struct {
	gid    string
	Online int32
	head   *Node
	lock   sync.RWMutex
}

func NewGroup(gid string) *Group {
	return &Group{
		gid: gid,
	}
}

func (g *Group) Put(node *Node) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.head == nil {
		g.head = node
		node.Pre, node.Next = nil, nil
	} else {
		head := g.head
		g.head = node
		node.Pre, node.Next = nil, head
	}
}

func (g *Group) Del(node *Node) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.head == node {
		g.head = node.Next
		node.Next = nil
		g.head.Pre = nil
	} else {
		node.Pre.Next = node.Next
		if node.Next != nil {
			node.Next.Pre = node.Pre
			node.Next = nil
		}
		node.Pre = nil
	}
}

func (g *Group) SendMsg(msg *message.Msg) {
	for h := g.head; h != nil; h = h.Next {
		// Don't forward to source
		if h.Id != int2str(msg.Uid) {
			h.Ch <- msg
		}
	}
}
