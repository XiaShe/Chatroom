package process2

import (
	"chatRoomProject/common/message"
	"chatRoomProject/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct{

}

// 转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 遍历服务端的onlineUsers map[int]*UserProcess
	// 将消息转发到各服务端

	// 取出mes的内容 SmsMes，取出其中userid，群发消息略过自己
	var smsMes message.SmMes

	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte[mes.Data], &smsMes) err =", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err =", err)
	}

	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}


}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	// Transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err =", err)
		return
	}
}
