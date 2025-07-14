package process

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	// 增加一個字段 表示該conn是哪個用戶的
	UserId int
}
// 通知所有其他在線用戶我上線了 
func (u *UserProcess) NotifyOtherOnlineUser(userId int){
	// 遍歷 onlineUsers 一個一個的通知
	fmt.Println("onlineUsers",userManager.onlineUsers)
	for id,userProcess := range userManager.onlineUsers{
		if id == userId {
			continue
		}
		// 開始通知 另外寫一個方法
		userProcess.NotifyMyOnlineMsg(userId)
	}
}
func (u *UserProcess) NotifyMyOnlineMsg(userId int){
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// 將notifyUserStatusMes序列化
	data,err := json.Marshal(notifyUserStatusMes)
	if err!=nil{
		fmt.Println("json.Marshal err",err)
		return 
	}
	mes.Data = string(data)
	 
	data,err = json.Marshal(mes)
	if err!=nil{
		fmt.Println("json.Marshal err",err)
		return 
	}
	// 發送
	// 創建一個transfer實例 不停的讀取服務器發送的消息
	transfer := &message.Transfer{
		Conn: u.Conn,
	}
	err = transfer.WritePkg(data)
	if err != nil{
		fmt.Println("notifyMyOnlineMsg WritePkg err",err )
		return 
	}
}
// 處理註冊請求
func (u *UserProcess) ServerProcessRegister(conn net.Conn,mes *message.Message)(err error){
	// 核心代碼
	// 1. 先從mes中取出mes.Data 並直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data),&registerMes)
	if err!=nil{
		fmt.Println("json.Unmarshal fail err=",err)
		return 
	}
	// 1 先聲明一個 resMes
	var resMes message.Message
	resMes.Type = message. RegisterResMesType
	// 2 聲明一個 LoginResMes
	var registerResMes message.RegisterResMes

	// 我們需要到redis去完成用戶註冊
	err = model.MyUserDao.Register(&registerMes.User)

	if err!=nil{
		if err ==model.ERROR_USER_EXISTS{
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error() 
		}else{
			registerResMes.Code = 506
			registerResMes.Error = "註冊發生未知錯誤"
		}
	}else{
		registerResMes.Code = 200
		fmt.Println("註冊成功")
	}

	// 3 將loginResMes 序列化
	data,err := json.Marshal(registerResMes)
	if err!=nil{
		fmt.Println("json.Marshal fail",err)
		return 
	}
	// 4 將data賦值給resMes
	resMes.Data = string(data)

	// 5 對resMes進行序列化 準備發送
	data,err = json.Marshal(resMes)
	if err != nil{
		fmt.Println("json.Marshal fail",err)
		return 
	}
	// 6 發送 將其封裝到一個writePkg函數
	// 因為使用分層模式(mvc) 先創建一個Transfer實例 然後調用WritePkg方法
	transfer := &message.Transfer{
		Conn: u.Conn,
	}
	err = transfer.WritePkg(data)
	if err != nil{
		fmt.Println("writePkg fail",err)
		return 
	}
	return 
}

// 處理登入請求
func (u *UserProcess)ServerProcessLogin(conn net.Conn, mes *message.Message)(err error){
	// 核心代碼
	// 1. 先從mes中取出mes.Data 並直接反序列化成LoginMes
	var loginMes message.LoginMes
	err =json.Unmarshal([]byte(mes.Data),&loginMes)
	if err!=nil{
		fmt.Println("json.Unmarshal fail err=",err)
		return 
	}
	// 1 先聲明一個 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2 聲明一個 LoginResMes
	var loginResMes message.LoginResMes

	// 我們需要到redis去完成用戶驗證
	// 如果用戶id=100 pwd=123456 認為合法 否則不合法
	user,err := model.MyUserDao.Login(loginMes.UserId,loginMes.UserPwd)

	if err!=nil {
		if err == model.ERROR_USER_NOT_EXISTS{
			loginResMes.Code = 500 // 500狀態碼 表示該用戶不存在
			loginResMes.Error = model.ERROR_USER_NOT_EXISTS.Error()
		}else if   err  == model.ERROR_USER_PWD{
			loginResMes.Code = 403 
			loginResMes.Error = model.ERROR_USER_PWD.Error()
		}else{
			loginResMes.Code = 500 
			loginResMes.Error = "服務器內部錯誤"
		}

	}else{
		loginResMes.Code = 200
		// 這邊用戶登入成功 要把該用戶加入到userManager的onlineUsers中
		u.UserId = loginMes.UserId 
		userManager.AddOnlineUser(u)
		// 通知其他在線用戶我上線了
		u.NotifyOtherOnlineUser(loginMes.UserId )
		// 把當前在線的用戶列表放入loginResMes.UserId
		for id := range userManager.onlineUsers{
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user,"登入成功")
	} 

	// 3 將loginResMes 序列化
	data,err := json.Marshal(loginResMes)
	if err!=nil{
		fmt.Println("json.Marshal fail",err)
		return 
	}

	// 4 將data賦值給resMes
	resMes.Data = string(data)

	// 5 對resMes進行序列化 準備發送
	data,err = json.Marshal(resMes)
	if err != nil{
		fmt.Println("json.Marshal fail",err)
		return 
	}
	// 6 發送 將其封裝到一個writePkg函數
	// 因為使用分層模式(mvc) 先創建一個Transfer實例 然後調用WritePkg方法
	transfer := &message.Transfer{
		Conn: u.Conn,
	}
	err = transfer.WritePkg(data)
	if err != nil{
		fmt.Println("writePkg fail",err)
		return 
	}
	return 
}