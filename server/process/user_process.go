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
}

// 編寫一個函數 專門處理登錄請求
func (u *UserProcess)ServerProcessLogin(conn net.Conn, mes *message.Message){
	// 核心代碼
	// 1. 先從mes中取出mes.Data 並直接反序列化成LoginMes
	var loginMes message.LoginMes
	err :=json.Unmarshal([]byte(mes.Data),&loginMes)
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
		if err.Error()== model.ERROR_USER_NOT_EXISTS{
			loginResMes.Code = 500 // 500狀態碼 表示該用戶不存在
			loginResMes.Error = model.ERROR_USER_NOT_EXISTS
		}else if   err.Error() == model.ERROR_USER_PWD{
			loginResMes.Code = 403 
			loginResMes.Error = model.ERROR_USER_PWD
		}else{
			loginResMes.Code = 500 
			loginResMes.Error = "服務器內部錯誤"
		}

	}else{
		loginResMes.Code = 200
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
}