package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	BaseController
}

func (c App) Index() revel.Result {
	c.RenderArgs["test"] = "Hello world!"

	return c.Render()
}
