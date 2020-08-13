package model

// 定义用户结构体

type User struct {
	// 序列化与反序列化成功，需保证用户json字符串key与结构体对应的tag名字对应一致
	UserId int	`json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
}