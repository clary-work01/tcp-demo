package process

import (
	"chatroom/client/model"
	"chatroom/common"
	"chatroom/common/message"
	"fmt"
	"net"
)
 
var onlineUsers = make(map[int]*common.User,10)
// 在登入成功後 對currentUser初始化 
var currentUser model.CurrentUser 

func InitCurrentUser(conn net.Conn,userId int,userStatus int){
	currentUser.Conn = conn
	currentUser.UserId = userId
	currentUser.UserStatus = userStatus
}

// 在客戶端顯示當前在線的用戶
func showOnlineUsers(){
	fmt.Println("當前在線用戶列表：")
	for id  := range onlineUsers{
		fmt.Println("用戶id:\n",id)
	}
}

func updateUserStatus(notifyUserStatusMes  *message.NotifyUserStatusMes){
	user,ok :=onlineUsers[notifyUserStatusMes.UserId] 
	if !ok { 
		// 若onlineUsers中原本沒有此用戶 才加入此用戶id
		user = &common.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	showOnlineUsers()
}
