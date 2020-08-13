package model

import (
	"chatRoomProject/common/message"
	"net"
)

//  客户端使用地方较多，将其作为全局变量 ---> smProcess 中
type CerUser struct {
	Conn net.Conn
	message.User
}