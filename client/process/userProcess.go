package process

import (
	"chatRoomProject/client/utils"
	"chatRoomProject/common/message"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type UserProcess struct {

}

// 用户登录 - 两次序列化
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	// 1. 连接到服务器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err =", err)
		return
	}

	// 延时关闭
	defer conn.Close()

	// 2. 准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType // 传入类型

	// 3. 创建一个 LoginMes 结构体 <包含用户信息：id-密码-昵称>
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4. 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("loginMes json.Marshal err =", err)
		return
	}

	// 5. 把data赋给 mes.Data 字段
	mes.Data = string(data)

	// 6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes json.Marshal err =", err)
		return
	}

	// 实例化utils包内结构体
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 7. 服务器发送数据给客户端
	err = tf.WritePkg(data)

	// #########################################################################//
	//							 处理服务器返回信息								//
	// #########################################################################//

	// 8. 处理服务器端返回的消息，判断是否登录成功
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err =", err)
	}
	// 将mes的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)

	if loginResMes.Code == 200 {
		// 初始化CurUser，登录后发送消息时使用
		CurUser.Conn = conn
		CurUser. UserId = userId
		CurUser.UserStatus = message.UserOnline

		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			fmt.Println("用户id：\t", v)

			// 完成客户端的 onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user

		}
		fmt.Print("\n\n")

		// 保持客户端与服务器的通讯
		go serverProcessMes(conn)
		// 1. 显示 登录成功后的菜单
		for {
			ShowMenu()
		}
	} else {
		// 2. 显示 登录失败的原因
		log.Println(loginResMes.Error)
		// fmt.Println(loginResMes.Error)
	}
	return
}


// 用户注册
func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 1. 连接到服务器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err =", err)
		return
	}

	// 延时关闭
	defer conn.Close()

	// 2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType // 传入类型

	// 3. 创建一个RegisterMes 结构体（赋值）
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("registerMes json.Marshal err =", err)
		return
	}

	// 5. 把data赋给 mes.Data 字段
	mes.Data = string(data)

	// 6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes json.Marshal err =", err)
		return
	}

	// 实例化utils包内结构体
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 7. 服务器发送数据给客户端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("将注册信息发送至服务器时出现错误， err =", err)
	}

	// #########################################################################//
	//							 处理服务器返回信息								//
	// #########################################################################//

	// 8. 处理服务器端返回的消息，判断是否登录成功
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err =", err)
	}
	// 将mes的data部分反序列化成RegisterResMes
	var RegisterResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterResMes)

	if RegisterResMes.Code == 200 {
		fmt.Println("注册成功！！！请重新登录...")
		// 退出
		os.Exit(0)
	} else {
		// 2. 显示 登录失败的原因
		fmt.Println(RegisterResMes.Error)
		// 退出
		os.Exit(0)
	}
	return

}