package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
)

// 編寫一個函數 專門處理登錄請求
func ServerProcessLogin(conn net.Conn, mes *message.Message){
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

	// 如果用戶id=100 pwd=123456 認為合法 否則不合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456"{
		// 合法
		loginResMes.Code = 200
	}else{
		// 不合法
		loginResMes.Code = 500 // 500狀態碼 表示該用戶不存在
		loginResMes.Error = "該用戶不存在 請註冊再使用..."
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
	err = message.WritePkg(conn,data)
	if err != nil{
		fmt.Println("writePkg fail",err)
		return 
	}
}