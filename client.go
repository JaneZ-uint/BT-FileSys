package main

import (
	"math/rand"
	"sync"
	"time"
)

type clientNode struct {
	node     dhtNode
	ip       string
	username string
}

func init() {
	// localAddress = GetLocalAddress()
	rand.Seed(time.Now().UnixNano())
	localAddress = "127.0.0.1"
	counter = 0
	userNum = 0
}

func (client *clientNode) login(username string) {
	client.node = NewNode(counter)
	client.ip = portToAddr(localAddress, counter)
	counter++
	client.username = username
	wg := new(sync.WaitGroup)
	client.node.Run(wg)
	if userNum == 0 {
		client.node.Create()
	} else {
		client.node.Join(client.ip)
	}
	userNum++
}

func (client *clientNode) logout(username string) {

}

func (client *clientNode) checkout() bool {

}

func (client *clientNode) upload() {

}

func (client *clientNode) download() {

}
