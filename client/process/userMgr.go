package process

import (
	"chatRoomProject/common/message"
	"fmt"
)

// 客户端需维护的map，即当前在线人数
var onlineUsers  = make(map[int]*message.User, 10)

// 在客户端显示当前在线的客户
func outputOnlineUser() {
	// 遍历一把 onlineUsers
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

// 处理 mes.Data 反序列化后 返回的 NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { // 如果没有，则创建实例
		user = &message.User{
			UserId:     notifyUserStatusMes.UserId,
		}
	}

	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserId] = user

	// 在客户端显示当前在线的客户
	outputOnlineUser()
}