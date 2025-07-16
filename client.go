package main

import (
	"math/rand"
	"sync"
	"time"
)

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
	wg := new(sync.WaitGroup)
	boss = NewNode(20000 + counter)
	counter++
	wg.Add(1)
	go boss.Run(wg)
	wg.Wait()
	/*for i := 0; i < total; i++ {
		wg.Add(1)
		if i == 0 {
			boss = NewNode(counter)
			counter++
			go boss.Run(wg)
		} else {
			var client dhtNode
			client = NewNode(counter)
			counter++
			go client.Run(wg)
		}
	}
	wg.Wait()*/
}

func Login(username string, client *dhtNode) {
	if userNum == 0 {
		(*client).Create()
	} else {
		(*client).Join(username)
	}
	userNum++
}

func Upload(filename UploadStruct, client *dhtNode) error {
	err := upload(filename.InputPath, filename.OutputPath, client)
	if err != nil {
		return err
	}
	return nil
}

func Download(filename DownloadStruct, client *dhtNode) error {
	err := download(filename.InputPath, filename.OutputPath, client)
	if err != nil {
		return err
	}
	return nil
}
