package message

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
)

type Message struct {
	Type string `json:"type"`
	Data string	`json:"data"`
}

type LoginMes struct {
	UserId int `json:"userId"` // 用戶id 
	UserPwd string `json:"userPwd"`  // 用戶密碼
	UserName string`json:"userName"` // 用戶名
}

type LoginResMes struct {
	Code int `json:"code"` // 返回狀態碼 500:用戶未註冊 200:登入成功
	Error string  `json:"error"` // 返回錯誤訊息
}