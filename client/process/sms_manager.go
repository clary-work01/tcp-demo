package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

// 處理接收到的群發消息 顯示即可
func ShowGroupMes(mes *message.Message){
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!=nil{
	fmt.Println("ShowGroupMes()失敗",err)
	}
	
	info := fmt.Sprintf("用戶id:%d對大家說：%s",smsMes.UserId,smsMes.Content)
	fmt.Println(info)
 }