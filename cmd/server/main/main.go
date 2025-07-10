package main

import (
	"fmt"
	"net"
)



func main(){
	fmt.Println("服務器在8889端口監聽...")
	listen ,err := net.Listen("tcp","0.0.0.0:8889")

	if err!= nil{
		fmt.Println("listen err",err)
	}
	// 一旦監聽成功 等待客戶端來連接服務器
	fmt.Println("listen success")

	defer listen.Close() // 延時關閉

	for{ // 循環等待客戶端連接我
		conn,err := listen.Accept()
		if err!=nil{
			fmt.Println("Accept error",err)
		}else{
			fmt.Println("Accept success",conn.RemoteAddr()) // 客戶端port會隨機分配
			// 一旦連接成功 啟動一個協程和客戶端保持通訊
			go processing(conn)
		}

	}
}
