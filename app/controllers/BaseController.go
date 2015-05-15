package controllers

import (
	// "bytes"
	// "encoding/json"
	// "github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	"myapp/app/table"
	// "gopkg.in/mgo.v2/bson"
	// "math"
	"strconv"
	// "strings"
)

// 公用Controller, 其它Controller继承它
type BaseController struct {
	*revel.Controller
}

func (c BaseController) SetSession(userInfo table.User) {
	if userInfo.UserId.Hex() != "" {
		c.Session["UserId"] = userInfo.UserId.Hex()
		c.Session["Email"] = userInfo.Email
		c.Session["Username"] = userInfo.Username
		c.Session["UsernameRaw"] = userInfo.UsernameRaw
		c.Session["Theme"] = userInfo.Theme
		c.Session["Logo"] = userInfo.Logo

		c.Session["NotebookWidth"] = strconv.Itoa(userInfo.NotebookWidth)
		c.Session["NoteListWidth"] = strconv.Itoa(userInfo.NoteListWidth)

		if userInfo.Verified {
			c.Session["Verified"] = "1"
		} else {
			c.Session["Verified"] = "0"
		}

		if userInfo.LeftIsMin {
			c.Session["LeftIsMin"] = "1"
		} else {
			c.Session["LeftIsMin"] = "0"
		}
	}
}
