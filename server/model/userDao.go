package model

import (
	"chatRoomProject/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// 在服务器启动后就初始化一个userDao实例，将其作为全局变量，需要redis操作时，直接使用即可
var (
	MyUserDao *UserDao
)

// 定义结构体UserDao，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 工厂模式，创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}



// 根据用户ID 返回一个User实例 + err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	// 通过给定ID去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "user", id))
	if err != nil {
		if err == redis.ErrNil { // * 表示在 user 哈希中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	// 将res反序列化成User实例
	user = &User{} //******将json反序列化到实例中才能使用方法*******
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal ~~~~~~~~~ err =", err)
		return
	}
	return
}

// 完成登录校验
/*
Login 如果用户id与密码都正确，则返回一个User实例
如果id或pwd有错误则返回对应错误信息
*/
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	// 从UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()

	// user为反序列化后的用户数据（用户名/密码/昵称）
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// * 验证用户名与密码是否匹配
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}


// 用户注册校验 + 信息入库
func (this *UserDao) Resgiter(user *message.User) (err error) {
	// 从UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()

	// 读取反序列化后的用户数据（用户名/密码/昵称）
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		// 没有错误，表示用户已经存在，不能注册
		err = ERROR_USER_EXISTS
		return
	}

	// 此时 id 在redis尚不存在，可以完成注册，将注册数据序列化准备入库
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// 注册信息入库
	_, err = conn.Do("HSet", "user", user.UserId, string(data))
	if err != nil {
		fmt.Println("注册信息入库出错， err =", err)
		return
	}

	return
}