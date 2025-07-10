package main

import (
	"math/rand"
	"time"
)

type clientNode struct {
	node     dhtNode
	ip       string
	username string
	password string
}

func init() {
	// localAddress = GetLocalAddress()
	rand.Seed(time.Now().UnixNano())
	localAddress = "127.0.0.1"
	counter = 0
}

func (client *clientNode) create(username string, password string) bool {
	client.node = NewNode(counter)
	client.ip = portToAddr(localAddress, counter)
	counter++
	client.username = username
	client.password = password
	return true
}

func (client *clientNode) login(username string, password string) bool {

	return true
}
