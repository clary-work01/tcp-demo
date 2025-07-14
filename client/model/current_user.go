package model

import (
	"chatroom/common"
	"net"
)

type CurrentUser struct{
	Conn net.Conn
	common.User     // 嵌入一個匿名結構體 
}
