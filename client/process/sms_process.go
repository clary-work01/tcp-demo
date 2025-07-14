package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct{

}

// 發送群聊消息
func (s *SmsProcess) SendGroupMsg(content string)(err error){
	// 創建一個mes實例
	var mes message.Message
	mes.Type = message.SmsMesType
	// 創建一個smsMes實例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = currentUser.UserId
	smsMes.UserStatus = currentUser.UserStatus
	// 序列化 smsMes
	data,err :=  json.Marshal(smsMes)
	if err!=nil{
		fmt.Println("SendGroupMsg() json.Marshal 失敗",err.Error())
		return 
	}
	mes.Data = string(data)
	// 序列化 mes
	data,err = json.Marshal(mes)
	if err!=nil{
		fmt.Println("SendGroupMsg() json.Marshal 失敗",err.Error())
		return 
	}
	// 把mes發送給服務器
	// 創建一個transfer實例 不停的讀取服務器發送的消息
	transfer := &message.Transfer{
		Conn: currentUser.Conn,
	}
	err = transfer.WritePkg(data)
	if err!=nil{
		fmt.Println("SendGroupMsg() 失敗",err.Error())
		return 
	}
	return 
} 