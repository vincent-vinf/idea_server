package im

import (
	"idea_server/db/mysql"
	"log"
	"sync"
)

type Bucket struct {
	groupsMap map[string]*Group
	nodeMap   map[string]*Node
	lock      sync.RWMutex
}

func NewBucket() *Bucket {
	return &Bucket{
		groupsMap: make(map[string]*Group),
		nodeMap:   make(map[string]*Node),
	}
}

func (b *Bucket) Put(node *Node) {
	groupList := mysql.AllUserGroup(node.Id)
	if groupList == nil {
		log.Println("id -> groups error")
		return
	}
	b.lock.Lock()
	for _, gid := range groupList {
		if group, ok := b.groupsMap[gid]; ok {
			group.Put(node)
		} else {
			group = NewGroup(gid)
			b.groupsMap[gid] = group
			group.Put(node)
		}
	}
	b.nodeMap[node.Id] = node
	b.lock.Unlock()
}

func (b *Bucket) Node(uid string) (node *Node) {
	b.lock.RLock()
	node = b.nodeMap[uid]
	b.lock.RUnlock()
	return
}
