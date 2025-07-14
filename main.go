package main

import (
	"fmt"
)

var localAddress string

var counter int

var userNum int

var boss dhtNode

const total = 100

func main() {
	var op string
	var username string
	var filename string
	for {
		fmt.Scanln(&op)
		if op == "login" { // 用户登录
			fmt.Scanln(&username)
			Login(username, &boss)
			fmt.Printf("User %s successfully [login]", username)
		} else if op == "upload" { //用户上传文件
			fmt.Scanln(&username, &filename)
			var err error
			err = Upload(UploadStruct{InputPath: filename, OutputPath: filename}, &boss)
			if err != nil {
				fmt.Printf("User %s successfully [upload] %s", username, filename)
			} else {
				fmt.Printf("Failed to upload.Please retry later.")
			}
		} else if op == "download" { //用户下载文件
			fmt.Scanln(&username, &filename)
			var err error
			err = Download(DownloadStruct{InputPath: filename, OutputPath: filename}, &boss)
			if err != nil {
				fmt.Printf("User %s successfully download %s", username, filename)
			} else {
				fmt.Printf("Failed to download.Please retry later.")
			}
		} else if op == "exit" {
			boss.ForceQuit()
			fmt.Println("Bye. Bet it has been a terrible experience.")
			break
		}
	}
}
