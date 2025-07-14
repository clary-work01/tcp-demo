package message

import (
	"chatroom/common"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
)
// 定義幾個用戶的狀態
const (
	UserOffline = iota
	UserOnline 
 	UserBusy 
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	UserId   int    `json:"userId"`   // 用戶id
	UserPwd  string `json:"userPwd"`  // 用戶密碼
	UserName string `json:"userName"` // 用戶名
}

type LoginResMes struct {
	Code  int    `json:"code"`  // 返回狀態碼 500:用戶未註冊 200:登入成功
	UsersId []int 
	Error string `json:"error"` // 返回錯誤訊息
}

type RegisterMes struct {
	User common.User `json:"user"`
 }

type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回狀態碼 400:用戶已經存在ㄑ 200:註冊成功
	Error string `json:"error"` // 返回錯誤訊息
}
// 為了配合服務器端推送用戶狀態變化的消息
type NotifyUserStatusMes struct{
	UserId int `json:"userId"`
	Status int `json:"status"`
}

// 把函數放到結構體中
type Transfer struct{
	Conn net.Conn
	Buf [8096]byte
}
func (t *Transfer) ReadPkg()(mes Message, err error){
	// 讀取客戶端發送的消息
	// 從conn讀 4個字節 到buf去
	_,err = t.Conn.Read(t.Buf[:4])
	if err!=nil{
		err = errors.New("read pkg header error")
		return 
	}
	// 把buf[:4] 從[]byte 轉成 int
    pkgLen := int(binary.BigEndian.Uint32(t.Buf[:4]))

	// 從conn讀 pkgLen個字節 到buf去
	n,err := t.Conn.Read(t.Buf[:pkgLen])
	if n != pkgLen || err != nil{
		err = errors.New("read pkg body error")
		return 
	}

	// 把 buf 反序列化成 message.Message
	err = json.Unmarshal(t.Buf[:pkgLen],&mes)
	if err!=nil {
		fmt.Println("json.Unmarshal err= ",err)
		return 
	}
	
	fmt.Printf("\n\nReadPkg:長度=%d\t內容%s\n\n",pkgLen,mes)
	return 
}
func (t *Transfer) WritePkg(data []byte)(err error){
	// 先發送一個長度給對方
	buf := make([]byte, 4) // int32 需要4 byte
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	// 發送data長度
	n,err := t.Conn.Write(buf)
	if n != 4 || err!=nil{
		fmt.Println("conn.Write失敗 err=",err)
		return 
	}
	
	// 發送data本身
	_,err = t.Conn.Write(data)
	if err !=nil {
		fmt.Println("conn.Write() fail",err)
		return 
	}

	fmt.Printf("\n\nWritePkg:長度=%d\t內容%s\n\n",len(data),string(data))
	return 
}