package main

import (
	"math/rand"
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
	counter++
}

func Login(username string) {
	if userNum == 0 {
		nodes[0].Create()
		nameMap[username] = 0
		boss = nodes[0]
		userNum++
		isExist[username] = true
	} else {
		if isExist[username] {
			userNum = nameMap[username]
			boss = nodes[userNum]
			return
		}
		//fmt.Print(userNum, "   ")
		isExist[username] = true
		//fmt.Println(nodeAddresses[userNum])
		nodes[userNum].Join(nodeAddresses[userNum])
		nameMap[username] = userNum
		boss = nodes[userNum]
		userNum++
	}
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
