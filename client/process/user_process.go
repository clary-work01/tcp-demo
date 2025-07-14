package process

import (
	"chatroom/common"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)


type UserProcess struct {

}

func ConnectToServer()(conn net.Conn){
	// 1. 連接到服務器
	conn,err:=  net.Dial("tcp","localhost:8889") // 之後會去讀配置文件
	if err!=nil{
		fmt.Println("net.Dial err=",err)
		return 
	}
	// 拿到conn後 應該馬上寫一個延時關閉!!!!!!!
	defer conn.Close()

	return conn
}

func (u *UserProcess)Login(userId int, userPwd string)(err error){
	// 開始訂協議

	// 1. 連接到服務器
	conn,err:=  net.Dial("tcp","localhost:8889") // 之後會去讀配置文件
	if err!=nil{
		fmt.Println("net.Dial err=",err)
		return 
	}
	// 拿到conn後 應該馬上寫一個延時關閉!!!!!!!
	defer conn.Close()

	// 2. 準備通過conn發送消息給服務
	var mes message.Message
	mes.Type = message.LoginMesType
	
	// 3. 創建一個LoginMes結構體實例
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

    // 4. 將loginMes序列化
	data,err := json.Marshal(loginMes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return 
	}
	// 5. 把data賦給mes.Data
	mes.Data = string(data)

	// 6.將mes序列化
	data,err = json.Marshal(mes)
	if err!=nil {
		fmt.Println("json.Marshal err=",err)
		return 
	}

	// 7. data就是我們要發送的消息
	// 因為使用分層模式(mvc) 先創建一個Transfer實例 然後調用WritePkg方法
	transfer := &message.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(data)
	if err!= nil{
		fmt.Println("登入發送信息錯誤write pkg fail",err)
		return
	}
	
	// 處理服務器返回的消息	...
	mes, err = transfer.ReadPkg()
	if err!=nil{
		fmt.Println("readPkg() fail",err)
		return 
	}
	// 將 mes的Data反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data),&loginResMes)

	
	if loginResMes.Code == 200 {
		// 可以顯示當前在線用戶列表
		for _,v := range loginResMes.UsersId{
			// 在線用戶列表排除掉自己的id
			if v == userId {
				continue
			}
			user := &common.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		showOnlineUsers()
		// 這裡我們還需要在客戶端啟動一個協程
		// 該協程保持和服務器端的通訊 如果服務器有數據突送給客戶端 則接收並顯示在客戶端的終端	
		go ProcessServerMes(conn)
		for{
		   // 1. 顯示我們的登入成功菜單
			ShowMenu()
		}
	}else{
		fmt.Println(loginResMes.Error)
	} 
	
	return  
}	

func (u *UserProcess)Register(userId int, userPwd string,userName string)(err error){
	// 1. 連接到服務器
	conn,err:=  net.Dial("tcp","localhost:8889") // 之後會去讀配置文件
	if err!=nil{
		fmt.Println("net.Dial err=",err)
		return 
	}
	// 拿到conn後 應該馬上寫一個延時關閉!!!!!!!
	defer conn.Close()

	// 2. 準備通過conn發送消息給服務
	var mes message.Message
	mes.Type = message.RegisterMesType
	
	// 3. 創建一個RegisterMes結構體實例
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 將registerMes序列化
	data,err := json.Marshal(registerMes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return 
	}

	// 5. 把data賦給mes.Data
	mes.Data = string(data)

	// 6.將mes序列化
	data,err = json.Marshal(mes)
	if err!=nil {
		fmt.Println("json.Marshal err=",err)
		return 
	}

	// 7. data就是我們要發送的消息
	// 因為使用分層模式(mvc) 先創建一個Transfer實例 然後調用WritePkg方法
	transfer := &message.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(data)
	if err!= nil{
		fmt.Println("註冊發送信息錯誤 write pkg fail",err)
		return
	}

	// 處理服務器返回的消息	...
	mes, err = transfer.ReadPkg()
	if err!=nil{
		fmt.Println("readPkg() fail",err)
		return 
	}
	// 將 mes的Data反序列化成LoginResMes
	var regitserResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data),&regitserResMes)
	if regitserResMes.Code == 200 {
		fmt.Println("註冊成功,請重新登入")
		os.Exit(0)
	}else{
		fmt.Println(regitserResMes.Error)
		os.Exit(0)
	} 
	
	return  	
}