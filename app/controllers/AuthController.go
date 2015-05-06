package controllers

import (
	"github.com/revel/revel"
)

// 用户登录/注销/找回密码

type Auth struct {
	BaseController
}

//--------
// 登录
func (c Auth) Login(email, from string) revel.Result {
	c.RenderArgs["title"] = c.Message("login")
	c.RenderArgs["subTitle"] = c.Message("login")
	c.RenderArgs["email"] = email
	c.RenderArgs["from"] = from
	return c.RenderTemplate("Home/login.html")
}
