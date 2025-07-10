package main

import (
	"chatroom/common/message"
	"chatroom/server/process"
	"fmt"
	"io"
	"net"
)

// 處理和客戶端的通訊
func processing(conn net.Conn){
	// 延時關閉conn !!!
	defer conn.Close()

	fmt.Printf("等待客戶端:%s輸入...\n",conn.RemoteAddr().String())
	// 循環讀取客戶端發送的消息
	for{
		// 把讀取數據包封裝成一個函數readPkg() 返回message.Message,error
		mes ,err :=message.ReadPkg(conn)
		if err!=nil{
			if err == io.EOF{
				fmt.Println("客戶端已退出")
				return 
			}else{
				fmt.Println("readPkg err=",err)
				return 
			}
		}

		fmt.Println("mes=",mes)

		err = serverProcessMes(conn,&mes) 	
		if err!=nil{
			fmt.Println("server process mes fail",err)
			return
		}
	}
}


// 根據客戶端發送消息種類不同 決定調用哪個函數來處理
func serverProcessMes(conn net.Conn,mes *message.Message)(err error){
	switch mes.Type{
		case message.LoginMesType:
			// 處理登入
			process.ServerProcessLogin(conn,mes)
		case message.RegisterMesType:
			// 處理註冊
		default : 
			fmt.Println("消息類型不存在 無法處理")
	}
	return 
}