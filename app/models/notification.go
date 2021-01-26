package models

import (
	"github.com/jinzhu/gorm"
)

// Notification model
type Notification struct {
	gorm.Model
	User    User   `gorm:"foreignkey:UserID;association_autoupdate:false;association_autocreate:false"`
	UserID  uint   //Gorm's variable for auto Preload function
	Role    Role   `gorm:"foreignkey:RoleID;association_autoupdate:false;association_autocreate:false"`
	RoleID  uint   //Gorm's variable for auto Preload function
	Viewed  string `gorm:"size:5"`
	Text    string `gorm:"type:text"`
	Subject string `gorm:"type:text"`
}
