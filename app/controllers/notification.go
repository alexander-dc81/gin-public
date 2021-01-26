package controllers

import (
	"gin/app/models"

	"github.com/revel/revel"
)

//Notification ...
type Notification struct {
	*revel.Controller
}

//All ...
func (c Notification) All(userID uint) revel.Result {
	notifications := &[]models.Notification{}
	DB.Limit(20).Where("user_id = ?", userID).Find(&notifications)

	response := models.JSONResponse{}
	response.Data = *notifications

	return c.RenderJSON(response)
}

//Save ...
func (c Notification) Save(notification *models.Notification) revel.Result {
	go notifyUsersByRole(notification) // Go routine (thread) for async notification creation

	return c.Redirect(Notification.List)
}

//Edit ...
func (c Notification) Edit() revel.Result {
	roles := &[]models.Role{}

	DB.Find(&roles)
	c.ViewArgs["roles"] = *roles

	return c.Render()
}

func notifyUsersByRole(notification *models.Notification) {
	users := &[]models.User{}

	DB.Where("role_id = ?", notification.Role.ID).Find(&users)

	for _, user := range *users {
		notification.User = user
		notification.UserID = user.ID

		DB.Save(&notification)
	}
}

//List ...
func (c Notification) List() revel.Result {
	notifications := &[]models.Notification{}

	DB.Limit(20).Preload("Role").Find(&notifications)

	c.ViewArgs["notifications"] = *notifications

	return c.Render()
}

//Delete ...
func (c Notification) Delete(ID int) revel.Result {
	notification := &models.Notification{}

	DB.Where("ID = ?", ID).Delete(&notification)

	return c.Redirect(Notification.List)
}

//View ...
func (c Notification) View(ID int) revel.Result {
	notification := &models.Notification{}

	DB.Where("ID = ?", ID).Find(&notification)

	notification.Viewed = "true"

	DB.Save(&notification)

	response := models.JSONResponse{}
	response.Data = *notification

	return c.RenderJSON(response)
}
