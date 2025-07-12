package main

import (
	"fmt"
)

var localAddress string

var counter int

var userNum int

func main() {
	var op string
	var username string
	var filename string
	for {
		fmt.Scanln(&op)
		// 用户登录
		if op == "login" {
			fmt.Scanln(&username)
			fmt.Printf("User %s successfully login", username)
		} else if op == "logout" { //用户登出
			fmt.Scanln(&username)
			fmt.Printf("User %s successfully logout", username)
		} else if op == "upload" { //用户上传文件
			fmt.Scanln(&username, &filename)
			var err error
			if err != nil {
				fmt.Printf("User %s successfully upload %s", username, filename)
			} else {
				fmt.Printf("Failed to upload.Please retry later.")
			}
		} else if op == "download" { //用户下载文件
			fmt.Scanln(&username, &filename)
			var err error
			if err != nil {
				fmt.Printf("User %s successfully download %s", username, filename)
			} else {
				fmt.Printf("Failed to download.Please retry later.")
			}
		} else if op == "check" { //用户权限查看
			fmt.Scanln(&username)
			var IsAvailable bool
			if IsAvailable {
				fmt.Printf("You can do whatever you want.")
			} else {
				fmt.Printf("You've already downloaded too much files without uploading your own. No more downloading!")
			}
		} else if op == "exit" {
			//TODO
			fmt.Println("Bye. Bet it has been a wonderful experience.")
			break
		}
	}
}
