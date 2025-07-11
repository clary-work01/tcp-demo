package process

import (
	"chatroom/common/message"
	"fmt"
	"net"
	"os"
)

// 顯示登入成功後的介面...

func ShowMenu(){
	fmt.Println("------恭喜xxx登入成功------")
	fmt.Println("------1. 顯示在線用戶列表------")
	fmt.Println("------2. 發送消息------")
	fmt.Println("------3. 信息列表------")
	fmt.Println("------4. 退出系統------")
	fmt.Println("------請選擇(1-4):------")

	var key int
	fmt.Scanf("%d\n",&key)
	switch key{
	case 1:
		fmt.Println("顯示在線用戶列表")
	case 2:
		fmt.Println("發送消息")
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
		 
		fmt.Printf("mes=%v\n",mes)
	}
	
}