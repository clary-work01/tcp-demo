package process

import "fmt"

// 顯示登入成功後的介面...

func ShowMenu(){
	fmt.Println("------恭喜xxx登入成功------")
	fmt.Println("------1. 顯示在線用戶列表------")
	fmt.Println("------2. 發送消息------")
	fmt.Println("------3. 信息列表------")
	fmt.Println("------4. 退出系統------")
	fmt.Println("------請選擇(1-4):------")
}