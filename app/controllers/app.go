package controllers

import (
	"gin/app/models"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

//App ...
type App struct {
	*revel.Controller
}

//AddUser ...
func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

func (c App) getUser(email string) (user *models.User) {
	user = &models.User{}
	c.Session.GetInto("fulluser", user, false)

	if user.Email == email {
		return user
	}

	user.ID = 0
	DB.Where("email = ?", email).Preload("Role").Find(&user)

	c.Session["fulluser"] = user

	return
}

//Index ...
func (c App) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(Dashboard.Dashboard)
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

//Register ...
func (c App) Register() revel.Result {
	return c.Render()
}

//Login ...
func (c App) Login(email string, password string, remember bool) revel.Result {
	user := c.getUser(email)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = email
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}

			return c.Redirect(Dashboard.Dashboard)
		}
	}

	c.Flash.Out["username"] = email
	c.Flash.Error("Login failed")
	return c.Redirect(App.Index)
}

//Logout ...
func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(App.Index)
}
