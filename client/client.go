package client

import "BT/core"

type clientNode struct {
	node      *core.KademliaNode
	username  string
	password  string
	previlege int
}

func (client *clientNode) create(username string, password string) bool {
	return true
}

func (client *clientNode) login(username string, password string) bool {
	return true
}
