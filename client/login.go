package client

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func Login(userId int, userPwd string)(err error){

	// 開始訂協議
	// return nil

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
	// 7.1 先把data的長度發送給服務器 
	// 使用 BigEndian 把data長度從 int轉[]byte
	buf := make([]byte, 4) // int32 需要4 byte
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	fmt.Printf("BigEndian: %v\n", buf)
	// 發送data長度
	n,err := conn.Write(buf)
	if n != 4 || err!=nil{
		fmt.Println("conn.Write失敗 err=",err)
		return 
	}
	fmt.Printf("客戶端成功發送長度=%d 內容%s",len(data),string(data))

	// 發送消息本身
	_,err = conn.Write(data)
	if err !=nil {
		fmt.Println("conn.Write() fail",err)
		return 
	}
	// 處理服務器返回的消息	...
	time.Sleep(20*time.Second)
	fmt.Println("休眠了20...")
    
	return 
}	