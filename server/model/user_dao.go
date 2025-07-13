package model

import (
	"chatroom/common"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// 在服務器啟動後 就初始化一個UserDao實例
// 把它做成全局變量 在需要和redis操作時 直接使用即可
var (
	MyUserDao *UserDao
)

// 完成對User結構體的各種操作
type UserDao struct {
	rdb *redis.Client
}

func NewUserDao(c *redis.Client)(userDao *UserDao){
	userDao = &UserDao{
		rdb: c,
	}
	return 
}

func(u *UserDao) GetUserById(userId int) (*common.User, error) {
	
	// 執行 HGET 命令
	result, err := u.rdb.HGet( "users", strconv.Itoa(userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("用戶 %d 不存在", userId)
		}
		return nil, fmt.Errorf("redis 錯誤: %s", err)
	}

	var user common.User
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析錯誤: %v", err)
	}

	return &user, nil
}

// Login 完成登入的校驗
func (u *UserDao) Login(userId int, userPwd string) (user *common.User, err error) {
	// 獲取用戶資料
	user, err = u.GetUserById(userId)
	if err != nil {
		err = ERROR_USER_NOT_EXISTS
		return nil, err
	}
	
	// 校驗密碼
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return nil, err
	}
	
	return user, nil
}

// Register 完成註冊的校驗
func (u *UserDao) Register(user *common.User) (err error) {
	// 獲取用戶資料
	_, err = u.GetUserById(user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return 
	}
	// 此時 說明用戶id還未註冊過
	data,err := json.Marshal(user)
	if err!=nil {
		return 
	}
	// 保存到redis庫
    _, err = u.rdb.HSet("users",strconv.Itoa(user.UserId), string(data)).Result()
	if err!=nil {
		fmt.Println("保存用戶註冊錯誤",err)
	}
	return 
}