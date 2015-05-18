package controllers

import (
	"github.com/revel/revel"
	"myapp/app/models"
	"myapp/app/table"
	"strings"
)

var note *models.Note
var session *models.Session
var auth *models.Auth
var user *models.User

// 拦截器
// 不需要拦截的url
// Index 除了Note之外都不需要
var commonUrl = map[string]map[string]bool{
// "App": map[string]bool{"Index": true},
}

func needValidate(controller, method string) bool {
	// 在里面
	if v, ok := commonUrl[controller]; ok {
		// 在commonUrl里
		if _, ok2 := v[method]; ok2 {
			return false
		}
		return true
	} else {
		// controller不在这里的, 肯定要验证
		return true
	}
}

func AuthInterceptor(c *revel.Controller) revel.Result {
	// 全部变成首字大写
	var controller = strings.Title(c.Name)
	var method = strings.Title(c.MethodName)

	// 是否需要验证?
	if !needValidate(controller, method) {
		return nil
	}
	// 验证是否已登录
	if userId, ok := c.Session["UserId"]; ok && userId != "" {
		return nil // 已登录
	}

	// 没有登录, 判断是否是ajax操作
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		re := table.NewRe()
		re.Msg = "NOTLOGIN"
		return c.RenderJson(re)
	}

	return c.Redirect("/login")
}

func Init() {
	note = models.NoteS
	session = models.SessionS
	auth = models.AuthS
	user = models.Users
}

func init() {
	// interceptor
	// revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Index{}) // Index.Note自己校验
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &App{})

}
