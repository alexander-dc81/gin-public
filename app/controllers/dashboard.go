package controllers

import (
	"github.com/revel/revel"
)

//Dashboard ...
type Dashboard struct {
	App
}

//Dashboard ...
func (c Dashboard) Dashboard() revel.Result {
	return c.Render()
}
