package main

import (
	"chatRoomProject/server/model"
	"fmt"
	"net"
	"time"
)

// 处理和客户端的通讯
func process(conn net.Conn) {
	// 延时关闭
	defer conn.Close()

	// 实例化结构体
	processor := &Processor{
		Conn: conn,
	}

	// 读取客户端发送的信息
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端与服务器端通讯出现问题 err =", err)
		return
	}
}

// 初始化UserDao
func initUserDao()  {
	// pool是全局变量
	// 初始化顺序：initPool ---> initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

// 服务器启动时，初始化 redis 连接池 与 redis 控制器（数据访问接口）
func init() {
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()

}

// 主线程： + 监听 + 处理通讯
func main() {
	// 监听
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}

	// 延时关闭
	defer listen.Close()

	// 监听成功后，等待客户端连接服务器，处理和客户端的通讯
	for {
		fmt.Println("等待客户端来连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err =", err)
		}

		// 连接成功，启动一个协程和客户端保持通讯
		go process(conn)
	}
}
