package process

import (
	"chatRoomProject/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) { // 传入mes类型为SmsMes
	// 显示服务器发来的消息
	// 1. 反序列化
	var smsMes message.SmMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &smsMes) err =", err)
		return
	}

	info := fmt.Sprintf("用户id：\t%d 对大家说：\t%s", smsMes.UserId, smsMes.Content)

	fmt.Println(info)
	fmt.Println()
}