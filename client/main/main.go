package main

import (
	"chatRoomProject/client/process"
	"fmt"
	"os"
)

// 定义两个全局变量，一个表示用户id，一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {

	// 接受用户的选择
	var key int
	// 判断是否继续显示菜单
	// var loop = true

	for true {
		fmt.Println("----------------欢迎登录-----------------")
		fmt.Println("\t\t\t 1. 登录聊天室")
		fmt.Println("\t\t\t 2. 注册用户")
		fmt.Println("\t\t\t 3. 退出系统")
		fmt.Println("\t\t\t 请选择（1-3）")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的ID号：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码：")
			fmt.Scanf("%s\n", &userPwd)

			// 登录
			// 创建一个实例
			up := process.UserProcess{}
			up.Login(userId, userPwd)

		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的id：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户的昵称：")
			fmt.Scanf("%s\n", &userName)

			// 创建实例，完成注册
			up := process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")

		}
	}

}
