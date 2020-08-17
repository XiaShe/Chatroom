package process

import (
	"chatRoomProject/client/utils"
	"chatRoomProject/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// 显示登录成功后的界面
func ShowMenu() {

	fmt.Println("------恭喜 xxx 登录成功---------")
	fmt.Println("------1. 显示在线用户列表--------")
	fmt.Println("------2. 发送消息---------------")
	fmt.Println("------3. 信息列表---------------")
	fmt.Println("------4. 退出系统---------------")
	fmt.Println("请选择（1-4）")

	var key int
	var content string // 输入消息
	// 定义实例，方便再switch内部调用
	smsProcess := &SmsProcess{}

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		// fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("请输入群发消息：")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content) // 群发消息
	case 3:
		fmt.Println("查看信息列表")
	case 4:
		fmt.Println("你选择退出了系统")
		os.Exit(0)
	default:
		fmt.Println("输入错误，请重新输入...")

	}
}

// 和服务器保持通讯
func serverProcessMes(conn net.Conn) {
	// 创建一个tranfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Printf("客户端正在读取服务器发送消息...\n")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err =", err)
		}

		// 成功读取到消息
		// 根据消息类型处理，此处只增加了上线功能
		switch mes.Type {
		case message.NotifyUserStatusMesType: // 有人上线了
			// 1. 取出 NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes // 用户ID与当前状态

			// 反序列化 ---> notifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 2. 把该用户的信息、状态保存到客户端map中
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:   // 有人群发消息了
			// 取出消息，反序列化
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")

		}

	}
}