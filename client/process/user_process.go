package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
)


type UserProcess struct {

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
		fmt.Println("write pkg fail",err)
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
		fmt.Println("登入成功")
	}else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	} 
	
	return 
}	