package main

import (
	"fmt"
)

var localAddress string

var counter int

var userNum int

var boss *clientNode

func main() {
	var op string
	var username string
	var filename string
	for {
		fmt.Scanln(&op)
		// 用户注册
		if op == "register" {
			fmt.Scanln(&username)
			Register()
			fmt.Printf("User %s successfully [register]", username)
		} else if op == "login" { // 用户登录
			fmt.Scanln(&username)
			boss.Login(username)
			fmt.Printf("User %s successfully [login]", username)
		} else if op == "upload" { //用户上传文件
			fmt.Scanln(&username, &filename)
			var err error
			err = boss.upload(UploadStruct{InputPath: filename, OutputPath: filename}, boss.ip)
			if err != nil {
				fmt.Printf("User %s successfully [upload] %s", username, filename)
			} else {
				fmt.Printf("Failed to upload.Please retry later.")
			}
		} else if op == "download" { //用户下载文件
			fmt.Scanln(&username, &filename)
			var err error
			err = boss.download(DownloadStruct{InputPath: filename, OutputPath: filename}, boss.ip)
			if err != nil {
				fmt.Printf("User %s successfully download %s", username, filename)
			} else {
				fmt.Printf("Failed to download.Please retry later.")
			}
		} else if op == "exit" {
			for i := 0; i < counter; i++ {
				var addr = portToAddr(localAddress, i)
				if addr == boss.ip {
					continue
				}
				boss.RemoteCall(addr, "ClientNode.QuitAll", "", nil)
			}
			boss.QuitAll("", nil)
			fmt.Println("Bye. Bet it has been a wonderful experience.")
			break
		}
	}
}
