package main

import (
	"fmt"
	"sync"
)

var localAddress string

var counter int

var userNum int

var boss dhtNode

const total = 100

var nodes [total + 1]dhtNode
var nodeAddresses [total + 1]string
var nameMap map[string]int
var isExist map[string]bool

func main() {
	nameMap = make(map[string]int)
	isExist = make(map[string]bool)
	userNum = 0
	wg := new(sync.WaitGroup)
	for i := 0; i <= total; i++ {
		nodes[i] = NewNode(20000 + i)
		nodeAddresses[i] = fmt.Sprintf("%s:%d", localAddress, 20000+i)
		wg.Add(1)
		go nodes[i].Run(wg)
	}
	wg.Wait()
	boss = nodes[0]

	var op string
	var username string
	var inputName string
	var outputName string
	fmt.Print("Welcome to JaneZ's DHT File Sharing System!\n")
	fmt.Print("Please input your operation:\n")
	for {
		op = ""
		username = ""
		inputName = ""
		outputName = ""
		fmt.Scanln(&op, &username, &inputName, &outputName)
		fmt.Print("\n")
		if op == "login" { // 用户登录
			Login(username)
			fmt.Printf("User %s successfully [login]", username)
		} else if op == "upload" { //用户上传文件
			var err error
			err = Upload(UploadStruct{InputPath: inputName, OutputPath: outputName}, &boss)
			if err == nil {
				fmt.Printf("User %s successfully [upload] %s", username, outputName)
			} else {
				fmt.Printf("Because of:", err)
				fmt.Printf("Failed to upload.Please retry later.")
			}
		} else if op == "download" { //用户下载文件
			var err error
			err = Download(DownloadStruct{InputPath: inputName, OutputPath: outputName}, &boss)
			if err == nil {
				fmt.Printf("User %s successfully download %s", username, outputName)
			} else {
				fmt.Printf("Because of:", err)
				fmt.Printf("Failed to download.Please retry later.")
			}
		} else if op == "exit" {
			for i := 0; i <= total; i++ {
				nodes[i].ForceQuit()
			}
			fmt.Println("Bye. Wish u a good day.")
			break
		}
		fmt.Print("\n")
	}
}
