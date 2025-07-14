package process

import "fmt"

// 因為userMananger在服務端有且只有一個
// 在很多地方都會用到 因此定義在全局變量
var (
	userManager *UserManager
)

// 維護在線用戶列表
type UserManager struct {
	onlineUsers map[int]*UserProcess
}	 

func init(){
	userManager = &UserManager{
		onlineUsers: make(map[int]*UserProcess,1024),
	}
}
 
func (u *UserManager)AddOnlineUser(userProcess *UserProcess){
	u.onlineUsers[userProcess.UserId] = userProcess
}

func (u *UserManager)DeleteOnlineUser(userId int){
	delete(u.onlineUsers,userId)
}

// 返回當前所有在線用戶
func (u *UserManager)GetAllOnlineUsers() map[int]*UserProcess{
	return u.onlineUsers
}
func (u *UserManager)GetOnlineUserById(userId int) (userProcess *UserProcess,err error){
	userProcess,ok := u.onlineUsers[userId]
	if !ok {
		fmt.Errorf("用戶%d不在線",userId)
		return 
	}
	return 
} 