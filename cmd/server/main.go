package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
)
func readPkg(conn net.Conn)(mes message.Message, err error){
	buf := make([]byte,8096)
	// 讀取客戶端發送的消息
	// 從conn讀 4個字節 到buf去
	_,err = conn.Read(buf[:4])
	if err!=nil{
		err = errors.New("read pkg header error")
		return 
	}
	fmt.Println("讀到的buf=",buf[:4])
	// 把buf[:4] 從[]byte 轉成 int
    pkgLen := int(binary.BigEndian.Uint32(buf[:4]))
    fmt.Printf("pkgLen= %d\n", pkgLen)

	// 從conn讀 pkgLen個字節 到buf去
	n,err := conn.Read(buf[:pkgLen])
	if n != pkgLen || err != nil{
		err = errors.New("read pkg body error")
		return 
	}

	// 把 buf 反序列化成 message.Message
	err = json.Unmarshal(buf[:pkgLen],&mes)
	if err!=nil {
		fmt.Println("json.Unmarshal err= ",err)
		return 
	}
	
	
		// 顯示客戶端輸入內容到終端
		// msg = msg + string(buf[:n])
		// fmt.Println(msg)
		return 
}

// 根據客戶端發送消息種類不同 決定調用哪個函數來處理
func serverProcessMes(conn net.Conn,mes *message.Message)(err error){
	switch mes.Type{
		case message.LoginMesType:
			// 處理登入
		case message.RegisterMesType:
			// 處理註冊
		default : 
			fmt.Println("消息類型不存在 無法處理")
	}
}
// 處理和客戶端的通訊
func process(conn net.Conn){
	// 延時關閉conn !!!
	defer conn.Close()


	fmt.Printf("等待客戶端:%s輸入...\n",conn.RemoteAddr().String())
	// 循環讀取客戶端發送的消息
	for{
		// 把讀取數據包封裝成一個函數readPkg() 返回message.Message,error
		mes ,err :=readPkg(conn)
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

	}
}


func main(){
	test,err := json.Marshal(1)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(test))
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
			go process(conn)
		}

	}
}
