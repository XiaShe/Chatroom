package message

// 消息类型
const (

	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType = "SmsMes"
)

// 定义用户状态常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

// 消息内容，为json格式
type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息的内容
}

// 定义消息（可增加）
type LoginMes struct {
	UserId   int    `json:"userId"`   // 用户名
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

// 服务器返回的消息内容
type LoginResMes struct {
	Code  int    `json:"code"`  // 返回状态码 500表示用户未注册， 200表示登录成功
	UsersId []int				   // 增加字段，保存用户id的切片
	Error string `json:"error"` // 返回错误信息
}

// 注册信息
type RegisterMes struct {
	User  User `json:"user"`
}

// 注册时服务器返回消息内容
type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回状态码 400用户注册成功， 200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

// 发送消息内容
type SmMes struct {
	Content string `json:"content"`
	User	// 匿名结构体 <继承>
}