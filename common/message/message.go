package message

import (
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
	Error string `json:"error"` // 返回錯誤訊息
}

type RegisterMes struct {
	// ..s

}

func ReadPkg(conn net.Conn)(mes Message, err error){
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
func WritePkg(conn net.Conn,data []byte)(err error){
	// 先發送一個長度給對方
	buf := make([]byte, 4) // int32 需要4 byte
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	fmt.Println("buf=",buf[:4])
	// 發送data長度
	n,err := conn.Write(buf)
	if n != 4 || err!=nil{
		fmt.Println("conn.Write失敗 err=",err)
		return 
	}
	
	// 發送data本身
	_,err = conn.Write(data)
	if err !=nil {
		fmt.Println("conn.Write() fail",err)
		return 
	}

	fmt.Printf("客戶端成功發送\n長度=%d\n內容%s\n",len(data),string(data))
	return 
}