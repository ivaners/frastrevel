package controllers

import (
	"github.com/revel/revel"
)

type Note struct {
	BaseController
}

func (c Note) Index() revel.Result {

	return c.Render()
}
