package controllers

import (
	"gin/app/models"

	"github.com/revel/revel"
)

//User ...
type User struct {
	*revel.Controller
}

//Dashboard ...
func (c User) Dashboard() revel.Result {
	users := &[]models.User{}

	DB.Preload("Role").Find(&users)

	c.ViewArgs["results"] = &users

	return c.Render()
}

//List ...
func (c User) List(searchText string) revel.Result {
	users := &[]models.User{}

	DB.Find(&users)

	c.ViewArgs["results"] = &users

	return c.Render()
}

//Edit ...
func (c User) Edit(ID int) revel.Result {
	user := &models.User{}
	roles := &[]models.Role{}

	DB.Where("ID = ?", ID).Preload("Role").Find(&user)
	DB.Find(&roles)

	c.ViewArgs["user"] = &user
	c.ViewArgs["roles"] = &roles

	return c.Render()
}

//Save ...
func (c User) Save(user *models.User) revel.Result {
	validateUser(c, user)

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Flash.Error("User information not valid!")
		// Keep the validation error from above by setting a flash cookie
		c.Validation.Keep()
		// Copies all given parameters (URL, Form, Multipart) to the flash cookie
		c.FlashParams()

		return c.Redirect(User.Edit, user.ID)
	}

	if user.ID == 0 {
		user.SetNewPassword(c.Params.Form["password"][0])
		DB.Create(&user)
	} else {
		userOld := &models.User{}
		DB.Where("ID = ?", user.ID).Find(&userOld)
		user.HashedPassword = userOld.HashedPassword

		DB.Save(&user)
	}

	return c.Redirect(User.Dashboard)

}

func validateUser(c User, user *models.User) {
	c.Validation.Required(user.Email).MessageKey("errors.email.required")
	c.Validation.Email(user.Email).MessageKey("errors.email.valid")

	if user.ID == 0 {
		password1 := c.Params.Form["password"][0]
		password2 := c.Params.Form["confirm_password"][0]
		c.Validation.Required(password1).MessageKey("errors.password.required")
		c.Validation.Required(password1 == password2).MessageKey("errors.password.match")
	}
}

//Delete ...
func (c User) Delete(ID int) revel.Result {
	user := &models.User{}
	DB.Where("ID = ?", ID).Delete(&user)

	return c.Redirect(User.Dashboard)
}
