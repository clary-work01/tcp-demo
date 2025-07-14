package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// 顯示登入成功後的介面...

func ShowMenu(){
	fmt.Println("------恭喜xxx登入成功------")
	fmt.Println("------1. 顯示在線用戶列表")
	fmt.Println("------2. 發送消息")
	fmt.Println("------3. 信息列表")
	fmt.Println("------4. 退出系統")
	fmt.Println("------請選擇(1-4):")

	var key int
	fmt.Scanf("%d\n",&key)
	switch key{
		case 1:
			// 顯示在線用戶列表
			showOnlineUsers()
		case 2:
			fmt.Println("請輸入你想對大家說的話")
			var content string
			fmt.Scanf("%s\n",&content)
			SmsProcess := &SmsProcess{}
			SmsProcess.SendGroupMsg(content)
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("你選擇了退出系統...")
			os.Exit(0)
		default :
			fmt.Println("你的輸入有誤 請重新輸入")
	}
}

// 和服務器端保持通訊
func ProcessServerMes(conn net.Conn){
	// 創建一個transfer實例 不停的讀取服務器發送的消息
	transfer := &message.Transfer{
		Conn: conn,
	}
	for {
		fmt.Printf("客戶端正在等待讀取服務器發送的消息")
		mes,err := transfer.ReadPkg()
		if err != nil{
			fmt.Println("ReadPkg() fail",err)
		}
		// 如果讀到訊息 
		switch mes.Type{
			case message.NotifyUserStatusMesType: // 有人上線了
				//  1. 取出 NotifyUserStatusMes
				var notifyUserStatusMes message.NotifyUserStatusMes
				json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes)

				//  2. 把此用戶的訊息和狀態保存到客戶端map[int]User中  
				updateUserStatus(&notifyUserStatusMes)
			case message.SmsMesType: // 有人群發消息了
				ShowGroupMes(&mes)
			default:
				fmt.Println("服務器端返回一個未知的消息類型 ")
		}
	}
}