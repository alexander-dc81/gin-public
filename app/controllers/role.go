package controllers

import (
	"gin/app/models"

	"github.com/revel/revel"
)

//Role ...
type Role struct {
	*revel.Controller
}

//List ...
func (c Role) List() revel.Result {
	roles := &[]models.Role{}

	DB.Find(&roles)

	c.ViewArgs["roles"] = &roles

	return c.Render()
}

//Delete ...
func (c Role) Delete(ID int) revel.Result {
	role := &models.Role{}
	DB.Where("ID = ?", ID).Delete(&role)

	return c.Redirect(Role.List)
}

//New ...
func (c Role) New(role models.Role) revel.Result {

	DB.Create(&role)

	return c.Redirect(Role.List)
}
