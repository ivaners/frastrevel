package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	. "myapp/app/lea"
	"myapp/app/table"
	// "strconv"
	"strings"
)

// 登录与权限

type Auth struct {
}

// pwd已md5了
func (this *Auth) Login(emailOrUsername, pwd string) table.User {
	emailOrUsername = strings.Trim(emailOrUsername, " ")
	//	pwd = strings.Trim(pwd, " ")
	userInfo := Users.LoginGetUserInfo(emailOrUsername, Md5(pwd))
	return userInfo
}

// 注册
/*
注册 leanote@leanote.com userId = "5368c1aa99c37b029d000001"
添加 在博客上添加一篇欢迎note, note1 5368c1b919807a6f95000000

将nk1(只读), nk2(可写) 分享给该用户
将note1 复制到用户的生活nk上
*/
// 1. 添加用户
// 2. 将leanote共享给我
// [ok]
func (this *Auth) Register(email, pwd, fromUserId string) (bool, string) {
	// 用户是否已存在
	if Users.IsExistsUser(email) {
		return false, "userHasBeenRegistered-" + email
	}
	user := table.User{UserId: bson.NewObjectId(), Email: email, Username: email, Pwd: Md5(pwd)}
	if fromUserId != "" && IsObjectId(fromUserId) {
		user.FromUserId = bson.ObjectIdHex(fromUserId)
	}
	LogJ(user)
	return this.register(user)
}

func (this *Auth) register(user table.User) (bool, string) {
	if Users.AddUser(user) {
	}

	return true, ""
}

//--------------
// 第三方注册

// 第三方得到用户名, 可能需要多次判断
func (this *Auth) getUsername(thirdType, thirdUsername string) (username string) {
	username = thirdType + "-" + thirdUsername
	i := 1
	for {
		if !Users.IsExistsUserByUsername(username) {
			return
		}
		username = fmt.Sprintf("%v%v", username, i)
	}
}

func (this *Auth) ThirdRegister(thirdType, thirdUserId, thirdUsername string) (exists bool, userInfo table.User) {
	userInfo = Users.GetUserInfoByThirdUserId(thirdUserId)
	if userInfo.UserId != "" {
		exists = true
		return
	}

	username := this.getUsername(thirdType, thirdUsername)
	userInfo = table.User{UserId: bson.NewObjectId(),
		Username:      username,
		ThirdUserId:   thirdUserId,
		ThirdUsername: thirdUsername,
	}
	_, _ = this.register(userInfo)
	return
}
