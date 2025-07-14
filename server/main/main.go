package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)


func main(){
	initRedis("localhost:6379",16,0,300*time.Second)

	model.InitUserDao(rdb)

	start()
}

func start(){
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
			go handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn){
	// 延時關閉conn !!!
	defer conn.Close()
	// 這裡創建一個總控 
	processer := &Processor{
		Conn: conn,
	}
	err := processer.processing()
	if err!=nil{
		fmt.Println("客戶端和服務端通訊協程錯誤",err)
		return 
	}
}