package utils

import (
"chatRoomProject/common/message"
"encoding/binary"
"encoding/json"
"fmt"
"net"
)

// 将方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte	// 传输时使用的缓冲

}

// 读取客户端发送给服务器的数据
func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	// buf := make([]byte, 1024*4)
	fmt.Println("正在读取客户端发送给服务器的数据....")

	// conn.Read 在conn没有被关闭的情况下，才会阻塞
	// 如果客户端关闭conn关闭了conn，则不会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		// err = errors.New("read pkg header error")
		// fmt.Println("conn.Read buf err =", err)
		return
	}

	// 根据buf[:4} 转成一个uint32类型---> pkgLen
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	// 根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen]) // 读取conn前pkgLen个数据到buf
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		// fmt.Println("conn.Read pkg err =", err)
		return
	}

	// 反序列化为 -> message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal XXXX err =", err)
		return
	}
	return
}

// 服务器发送数据给客户端
func (this *Transfer) WritePkg(data []byte) (err error) {

	// 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen) // 将uint32 ---> []byte

	// 发送长度
	_, err = this.Conn.Write(this.Buf[:4])
	if err != nil {
		fmt.Println("conn.Write(buf) fail err =", err)
		return
	}

	// 发送data本身
	n, err := this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(buf) fail err =", err)
		return
	}
	return
}