package controllers

import (
	"github.com/revel/revel"
	"myapp/app/table"
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

func (c Auth) DoLogin(email, pwd string, captcha string) revel.Result {
	sessionId := c.Session.Id()
	var msg = ""

	// > 5次需要验证码, 直到登录成功
	if session.LoginTimesIsOver(sessionId) && session.GetCaptcha(sessionId) != captcha {
		msg = "captchaError"
	} else {
		usertable := auth.Login(email, pwd)
		if usertable.Email != "" {
			c.SetSession(usertable)
			session.ClearLoginTimes(sessionId)
			return c.RenderJson(table.Re{Ok: true})
		} else {
			// 登录错误, 则错误次数++
			msg = "wrongUsernameOrPassword"
			session.IncrLoginTimes(sessionId)
		}
	}

	return c.RenderJson(table.Re{Ok: false, Item: session.LoginTimesIsOver(sessionId), Msg: c.Message(msg)})
}
