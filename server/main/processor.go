package main

import (
	"chatRoomProject/common/message"
	"chatRoomProject/server/process"
	"chatRoomProject/server/utils"
	"fmt"
	"io"
	"net"
)

// 创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

// 根据客户端发送的 消息种类（登录/注册） ，调用对应函数
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:					// 登录处理

		// 创建UserProcess实例
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case message.RegisterMesType:				// 注册处理
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:					// 群发消息
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)


	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

// 读取客户端发送的信息
func (this *Processor) process2() (err error) {

	for {

		// 创建一个Tranfer实例完成读包（conn）任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也退出...")
				return err
			} else {
				return err
			}
		}
		fmt.Println("发送的来的数据为：")
		fmt.Println(mes) // json格式

		// 根据客户端发送的 消息种类（登录/注册） ，调用对应函数
		err = this.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("this.serverProcessMes(&mes) err =", err)
			// return
		}

	}
}
