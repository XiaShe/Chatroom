package process2

import "fmt"

// 因为UserMgr实例在服务器端有且只有一个
// 因为在很多地方都会使用到，因此将其定义为全局变量
var (
	userMgr	*UserMgr
)

// 在线用户信息
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 对于 userMgr 进行初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 增加在线用户
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除在线用户
func (this *UserMgr) DelOnlineUser(UserId int) {
	delete(this.onlineUsers, UserId)
}

// 返回当前所有的在线用户
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

// 返回id对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	// 从map中取出值
	up, ok := this.onlineUsers[userId]
	if !ok { // 查找用户当前不在线
		err = fmt.Errorf("用户ID不存在", userId)
		return
	}
	return
}