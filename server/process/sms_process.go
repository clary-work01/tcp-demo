package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
)


type SmsProcess struct {

}
// 轉發消息方法
func (s *SmsProcess)SendGroupMes (mes *message.Message ){
	// 把mes拿出來看一下
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!=nil{
		fmt.Println("SendGroupMes()反序列化失敗",err.Error() )
		return 
	}
	// 原封不動去轉發
	data,err := json.Marshal(mes)
	if err!=nil{
		fmt.Println("SendGroupMes()序列化失敗",err.Error() )
		return 
	}

	// 遍歷userManager的onlineUsers 把消息轉發出去
	for id,userProcess := range userManager.onlineUsers{
		// 這邊要從onlineUsers過濾掉 發消息的客戶端 
		if id == smsMes.UserId{
			continue
		}
		s.SendMsgToEachOnlineUser(data,userProcess.Conn)
	}
}

func (s *SmsProcess)SendMsgToEachOnlineUser(data []byte,conn net.Conn){
	// 發送
	// 創建一個transfer實例 不停的讀取服務器發送的消息
	transfer := &message.Transfer{
		Conn: conn,
	}
	err := transfer.WritePkg(data)
	if err!=nil{
		fmt.Println("轉發消息失敗 err=",err)
	}
}