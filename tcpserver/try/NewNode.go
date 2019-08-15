package main

import (
	"github.com/lithammer/shortuuid"
	"net"
)

func NewNode(conn net.Conn, IDS []string) *Node {
	node := &Node{
		ID:      calculateID(IDS),
		country: findCou(),
		region:  findReg(),
		city:    findCity(),
	}

	return node
}

func calculateID(IDS []string) string {

	nodeID := shortuuid.New()
	if !isIDvalid(nodeID, IDS) {
		go calculateID(IDS)
	}
	return nodeID
}

func isIDvalid(ID string, IDS []string) bool {
	i := 0
	for _, id := range IDS {
		if ID != id {
			if i == len(IDS)-1 {
				return true
			}
			continue
		}
		i++
		break
	}
	return false
}
