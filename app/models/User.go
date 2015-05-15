package models

import (
	"gopkg.in/mgo.v2/bson"
	"myapp/app/db"
	. "myapp/app/lea"
	"myapp/app/table"
	"strings"
	"time"
)

type User struct {
}

// 自增Usn
// 每次notebook,note添加, 修改, 删除, 都要修改
func (this *User) IncrUsn(userId string) int {
	user := table.User{}
	query := bson.M{"_id": bson.ObjectIdHex(userId)}
	db.GetByQWithFields(db.Users, query, []string{"Usn"}, &user)
	usn := user.Usn
	usn += 1
	Log("inc Usn")
	db.UpdateByQField(db.Users, query, "Usn", usn)
	return usn
	//	return db.Update(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, bson.M{"$inc": bson.M{"ReadNum": 1}})
}

func (this *User) GetUsn(userId string) int {
	user := table.User{}
	query := bson.M{"_id": bson.ObjectIdHex(userId)}
	db.GetByQWithFields(db.Users, query, []string{"Usn"}, &user)
	return user.Usn
}

// 添加用户
func (this *User) AddUser(user table.User) bool {
	if user.UserId == "" {
		user.UserId = bson.NewObjectId()
	}
	user.CreatedTime = time.Now()

	// if user.Email != "" {
	// 	user.Email = strings.ToLower(user.Email)

	// 	// 发送验证邮箱
	// 	go func() {
	// 		emailService.RegisterSendActiveEmail(user, user.Email)
	// 		// 发送给我 life@leanote.com
	// 		emailService.SendEmail("life@leanote.com", "新增用户", "{header}用户名"+user.Email+"{footer}")
	// 	}()
	// }

	return db.Insert(db.Users, user)
}

// 通过email得到userId
func (this *User) GetUserId(email string) string {
	email = strings.ToLower(email)
	user := table.User{}
	db.GetByQ(db.Users, bson.M{"Email": email}, &user)
	return user.UserId.Hex()
}

// 得到用户名
func (this *User) GetUsername(userId string) string {
	user := table.User{}
	db.GetByQWithFields(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, []string{"Username"}, &user)
	return user.Username
}

// 是否存在该用户 email
func (this *User) IsExistsUser(email string) bool {
	if this.GetUserId(email) == "" {
		return false
	}
	return true
}

// 是否存在该用户 username
func (this *User) IsExistsUserByUsername(username string) bool {
	return db.Count(db.Users, bson.M{"Username": username}) >= 1
}

// 得到用户信息, userId, username, email
func (this *User) GetUserInfoByAny(idEmailUsername string) table.User {
	if IsObjectId(idEmailUsername) {
		return this.GetUserInfo(idEmailUsername)
	}

	if strings.Contains(idEmailUsername, "@") {
		return this.GetUserInfoByEmail(idEmailUsername)
	}

	// username
	return this.GetUserInfoByUsername(idEmailUsername)
}

func (this *User) setUserLogo(user *table.User) {
	// Logo路径问题, 有些有http: 有些没有
	if user.Logo == "" {
		user.Logo = "images/blog/default_avatar.png"
	}
	if user.Logo != "" && !strings.HasPrefix(user.Logo, "http") {
		user.Logo = strings.Trim(user.Logo, "/")
		// user.Logo = configService.GetSiteUrl() + "/" + user.Logo
	}
}

// 得到用户信息 userId
func (this *User) GetUserInfo(userId string) table.User {
	user := table.User{}
	db.Get(db.Users, userId, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

// 得到用户信息 email
func (this *User) GetUserInfoByEmail(email string) table.User {
	user := table.User{}
	db.GetByQ(db.Users, bson.M{"Email": email}, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

// 得到用户信息 username
func (this *User) GetUserInfoByUsername(username string) table.User {
	user := table.User{}
	username = strings.ToLower(username)
	db.GetByQ(db.Users, bson.M{"Username": username}, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

func (this *User) GetUserInfoByThirdUserId(thirdUserId string) table.User {
	user := table.User{}
	db.GetByQ(db.Users, bson.M{"ThirdUserId": thirdUserId}, &user)
	this.setUserLogo(&user)
	return user
}
func (this *User) ListUserInfosByUserIds(userIds []bson.ObjectId) []table.User {
	users := []table.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	return users
}
func (this *User) ListUserInfosByEmails(emails []string) []table.User {
	users := []table.User{}
	db.ListByQ(db.Users, bson.M{"Email": bson.M{"$in": emails}}, &users)
	return users
}

// 用户信息即可
func (this *User) MapUserInfoByUserIds(userIds []bson.ObjectId) map[bson.ObjectId]table.User {
	users := []table.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)

	userMap := make(map[bson.ObjectId]table.User, len(users))
	for _, user := range users {
		this.setUserLogo(&user)
		userMap[user.UserId] = user
	}
	return userMap
}

// 用户信息和博客设置信息
func (this *User) MapUserInfoAndBlogInfosByUserIds(userIds []bson.ObjectId) map[bson.ObjectId]table.User {
	return this.MapUserInfoByUserIds(userIds)
}

// 通过ids得到users, 按id的顺序组织users
func (this *User) GetUserInfosOrderBySeq(userIds []bson.ObjectId) []table.User {
	users := []table.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)

	usersMap := map[bson.ObjectId]table.User{}
	for _, user := range users {
		usersMap[user.UserId] = user
	}

	hasAppend := map[bson.ObjectId]bool{} // 为了防止userIds有重复的
	users2 := []table.User{}
	for _, userId := range userIds {
		if user, ok := usersMap[userId]; ok && !hasAppend[userId] {
			hasAppend[userId] = true
			users2 = append(users2, user)
		}
	}
	return users2
}

// 使用email(username), pwd得到用户信息
func (this *User) LoginGetUserInfo(emailOrUsername, md5Pwd string) table.User {
	emailOrUsername = strings.ToLower(emailOrUsername)

	user := table.User{}

	if strings.Contains(emailOrUsername, "@") {
		db.GetByQ(db.Users, bson.M{"Email": emailOrUsername, "Pwd": md5Pwd}, &user)
	} else {
		db.GetByQ(db.Users, bson.M{"Username": emailOrUsername, "Pwd": md5Pwd}, &user)
	}
	this.setUserLogo(&user)
	return user
}

// 更新username
func (this *User) UpdateUsername(userId, username string) (bool, string) {
	if userId == "" || username == "" || username == "admin" { // admin用户是内置的, 不能设置
		return false, "usernameIsExisted"
	}
	usernameRaw := username // 原先的, 可能是同一个, 但有大小写
	username = strings.ToLower(username)

	// 先判断是否存在
	userIdO := bson.ObjectIdHex(userId)
	if db.Has(db.Users, bson.M{"Username": username, "_id": bson.M{"$ne": userIdO}}) {
		return false, "usernameIsExisted"
	}

	ok := db.UpdateByQMap(db.Users, bson.M{"_id": userIdO}, bson.M{"Username": username, "UsernameRaw": usernameRaw})
	return ok, ""
}

// 修改头像
func (this *User) UpdateAvatar(userId, avatarPath string) bool {
	userIdO := bson.ObjectIdHex(userId)
	return db.UpdateByQField(db.Users, bson.M{"_id": userIdO}, "Logo", avatarPath)
}

//----------------------
// 已经登录了的用户修改密码
func (this *User) UpdatePwd(userId, oldPwd, pwd string) (bool, string) {
	userInfo := this.GetUserInfo(userId)
	if userInfo.Pwd != Md5(oldPwd) {
		return false, "oldPasswordError"
	}
	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Pwd", Md5(pwd))
	return ok, ""
}

// 管理员重置密码
// func (this *User) ResetPwd(adminUserId, userId, pwd string) (ok bool, msg string) {
// 	if configService.GetAdminUserId() != adminUserId {
// 		return
// 	}
// 	ok = db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Pwd", Md5(pwd))
// 	return
// }

// 修改主题
func (this *User) UpdateTheme(userId, theme string) bool {
	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Theme", theme)
	return ok
}

//---------------
// 修改email

// 注册后验证邮箱
// func (this *User) ActiveEmail(token string) (ok bool, msg, email string) {
// 	tokenInfo := table.Token{}
// 	if ok, msg, tokenInfo = tokenService.VerifyToken(token, table.TokenActiveEmail); ok {
// 		// 修改之后的邮箱
// 		email = tokenInfo.Email
// 		userInfo := this.GetUserInfoByEmail(email)
// 		if userInfo.UserId == "" {
// 			ok = false
// 			msg = "不存在该用户"
// 			return
// 		}

// 		// 修改之, 并将verified = true
// 		ok = db.UpdateByQMap(db.Users, bson.M{"_id": userInfo.UserId}, bson.M{"Verified": true})
// 		return
// 	}

// 	ok = false
// 	msg = "该链接已过期"
// 	return
// }

// 修改邮箱
// 在此之前, 验证token是否过期
// 验证email是否有人注册了
// func (this *User) UpdateEmail(token string) (ok bool, msg, email string) {
// 	tokenInfo := table.Token{}
// 	if ok, msg, tokenInfo = tokenService.VerifyToken(token, table.TokenUpdateEmail); ok {
// 		// 修改之后的邮箱
// 		email = tokenInfo.Email
// 		// 先验证该email是否被注册了
// 		if userService.IsExistsUser(email) {
// 			ok = false
// 			msg = "该邮箱已注册"
// 			return
// 		}

// 		// 修改之, 并将verified = true
// 		ok = db.UpdateByQMap(db.Users, bson.M{"_id": tokenInfo.UserId}, bson.M{"Email": email, "Verified": true})
// 		return
// 	}

// 	ok = false
// 	msg = "该链接已过期"
// 	return
// }

//---------
// 第三方添加账号
func (this *User) ThirdAddUser(userId, email, pwd string) (ok bool, msg string) {
	// 判断该用户是否已有了帐户
	userInfo := this.GetUserInfo(userId)
	if userInfo.Email != "" {
		msg = "你已有帐户"
		return
	}

	// 判断email是否存在
	if this.IsExistsUser(email) {
		msg = "该用户已存在"
		return
	}

	ok = db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"Email": email, "Pwd": Md5(pwd)})
	return
}

// 统计
func (this *User) CountUser() int {
	return db.Count(db.Users, bson.M{})
}
