package main

import (
	network "BT/Network"
	"math/rand"
	"sync"
	"time"
)

type clientNode struct {
	node dhtNode
	ip   string
	network.NetworkStation
}

type UploadStruct struct {
	InputPath  string
	OutputPath string
}

type DownloadStruct struct {
	InputPath  string
	OutputPath string
}

func init() {
	rand.Seed(time.Now().UnixNano())
	localAddress = "127.0.0.1"
	counter = 0
	userNum = 0
	boss = new(clientNode)
	boss.ip = portToAddr(localAddress, counter)
	boss.node = NewNode(counter)
	counter++
	wg := new(sync.WaitGroup)
	boss.node.Run(wg)
	boss.InitRPC(boss, "ClientNode")
	go boss.RunRPCServer(boss.ip, wg)
}

func Register() {
	var client clientNode
	client.ip = portToAddr(localAddress, counter)
	client.node = NewNode(counter)
	counter++
	wg := new(sync.WaitGroup)
	client.node.Run(wg)
	client.InitRPC(&client, "ClientNode")
	go client.RunRPCServer(client.ip, wg)
}

func (client *clientNode) Login(username string) {
	if userNum == 0 {
		client.node.Create()
	} else {
		client.node.Join(username)
	}
	userNum++
}

func (client *clientNode) upload(filename UploadStruct, Addr string) error {
	err := client.RemoteCall(Addr, "ClientNode.Upload", filename, &struct{}{})
	if err != nil {
		return err
	}
	return nil
}

func (client *clientNode) Upload(filename UploadStruct, reply *struct{}) error {
	err := upload(filename.InputPath, filename.OutputPath, &client.node)
	if err != nil {
		return err
	}
	return nil
}

func (client *clientNode) download(filename DownloadStruct, Addr string) error {
	err := client.RemoteCall(Addr, "ClientNode.Download", filename, &struct{}{})
	if err != nil {
		return err
	}
	return nil
}

func (client *clientNode) Download(filename DownloadStruct, reply *struct{}) error {
	err := download(filename.InputPath, filename.OutputPath, &client.node)
	if err != nil {
		return err
	}
	return nil
}

func (client *clientNode) QuitAll(_ string, reply *struct{}) error {
	client.node.ForceQuit()
	client.StopRPCServer()
	return nil
}
