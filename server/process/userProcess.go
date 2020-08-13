package process2

import (
	"chatRoomProject/common/message"
	"chatRoomProject/server/model"
	"chatRoomProject/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	// 表示该Conn属于哪个用户
	UserId int

}

// 编写通知所有在线用户的方法 ，userId 通知其它用户自己上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int)  {

	// 遍历 onlineUsers, 然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		// 通知其他人
		up.NotifyMeOnline(userId)
	}
}

// 通知其它用户 userId 上线了
func (this *UserProcess) NotifyMeOnline(userId int) {

	// 组装我们的 NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// 将notifyUserStatusMes 序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal(notifyUserStatusMes) err =", err)
		return
	}

	mes.Data = string(data)

	// 再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err =", err)
		return
	}

	// 创建Transfer实例，发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err =", err)
		return
	}

}

// 处理注册信息---> 与redis中数据进行对比 ---> 发送验证于客户端
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	var registerMes message.RegisterMes

	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err =", err)
		return
	}

	// 1. 先声明一个resMes（服务器登录返回信息，包括 消息类型 和 消息内容<即状态码> ）
	var resMes message.Message
	resMes.Type = message.RegisterResMesType // 因为本身就是登录处理函数，因此类型直接赋值

	// 2. 再声明一个registerResMes（服务器登录返回信息，显示是否登录成功） ，并完成赋值
	var registerResMes message.RegisterResMes

	// 3. 用户注册校验 + 信息入库
	err = model.MyUserDao.Resgiter(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS { // 用户id已经存在
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail err =", err)
		return
	}

	// 4. 将data赋值于resMes
	resMes.Data = string(data)

	// 5. 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err =", err)
		return
	}

	// 6. 发送data，将其封装到writePkg 函数中，向客户端发送消息
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}


// 处理登录信息---> 与redis中数据进行对比 ---> 发送验证于客户端
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码
	// 1. 先从mes 中取出mes.Data ，并反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err =", err)
		return
	}

	// 1. 先声明一个resMes（服务器登录返回信息，包括 消息类型 和 消息内容<即状态码> ）
	var resMes message.Message
	resMes.Type = message.LoginResMesType // 因为本身就是登录处理函数，因此类型直接赋值

	// 2. 再声明一个LoginResMes（服务器登录返回信息，显示是否登录成功） ，并完成赋值
	var loginResMes message.LoginResMes

	// **********************到redis数据库验证**********************

	// 2.1 使用 model.MyUserDao 到 redis 验证
	// 将客户端发送到服务器 的 用户名和密码 与数据库中的 用户名密码 进行对比
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS { // 用户不存在，错误在Login函数中已经赋值
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD { // 密码错误
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}

	} else {
		loginResMes.Code = 200

		// 这里把登录成功的用户放入到userMgr ,将登录成功的用户的UserId赋值给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)

		// 通知其它在线用户， UserId 上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)

		// 将当前在线用户的Id放入loginResMes.UsersId
		// 遍历userMgr.onlinUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}

		fmt.Println(user, "登录成功")
	}

	// 3. 将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail err =", err)
		return
	}

	// 4. 将data赋值于resMes
	resMes.Data = string(data)

	// 5. 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err =", err)
		return
	}

	// 6. 发送data，将其封装到writePkg 函数中
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
