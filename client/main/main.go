// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )

// // 客戶端發送單行數據 然後退出
// func sendOneLine(conn net.Conn){
// 	for {
// 	   reader:= bufio.NewReader(os.Stdin) // os.Stdin代表標準輸入[終端]
// 		// 從終端讀出一行用戶輸入
// 	   line,err := reader.ReadString('\n')
// 	   if err!=nil{
// 		   fmt.Println("readString err",err)
// 	   }
// 	   line = strings.TrimSpace(line)
// 	   if line ==  "exit"{
// 		fmt.Println("same")
// 		return
// 	   }

// 	   // 發給服務器
// 	   n,err :=conn.Write([]byte(line))
// 	   if err!=nil{
// 		   fmt.Println("發送數據失敗",err)
// 	   }
// 	   fmt.Printf("客戶端發送了%d字節的數據\n",n)
// 	}
// }

// func main(){
// 	conn,err := net.Dial("tcp","127.0.0.1:8888")

// 	if err!=nil {
// 		fmt.Println("client dial err",err)
// 	}

// 	fmt.Println("連接到server")

// 	sendOneLine(conn)
// }

package main

import (
	"chatroom/client/process"
	"fmt"
)

var(
	userId int
	userPwd string
	userName string
)

func main(){
	// 接收用戶的選擇
	var key int

	for {
		fmt.Println("----------------多人聊天系統-----------------")
		fmt.Println("\t\t 1 登入聊天室")
		fmt.Println("\t\t 2 註冊用戶")
		fmt.Println("\t\t 3 退出系統")
		fmt.Println("\t\t 請選擇(1-3):")

		fmt.Scanf("%d\n",&key)

		switch key{
		case 1: 
			fmt.Println("登入聊天室")

			fmt.Println("請輸入用戶id")
			fmt.Scanf("%d\n",&userId)	
			fmt.Println("請輸入用戶密碼")
			fmt.Scanf("%s\n",&userPwd)	
			// 創建一個UserProcess實例
			userProcess := &process.UserProcess{}
			err := userProcess.Login(userId,userPwd)
			if err != nil{
				fmt.Println("userProcess.Login() fail",err)
			}
		case 2:
			fmt.Println("註冊用戶")

			fmt.Println("請輸入用戶id")
			fmt.Scanf("%d\n",&userId)	
			fmt.Println("請輸入用戶密碼")
			fmt.Scanf("%s\n",&userPwd)	
			fmt.Println("請輸入用戶暱稱")
			fmt.Scanf("%s\n",&userName)	 
			// 創建一個UserProcess實例
			userProcess := &process.UserProcess{}
			err := userProcess.Register(userId,userPwd,userName)
			if err != nil{
				fmt.Println("userProcess.Register() fail",err)
			}
		case 3:
			fmt.Println("退出系統")
		default :
			fmt.Println("你的輸入有誤，請重新輸入")
		}
	}

}