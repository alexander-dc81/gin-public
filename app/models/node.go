package models

import (
	"github.com/jinzhu/gorm"
)

//Node model
type Node struct {
	gorm.Model
	Name   string `gorm:"size:255"`
	Type   Type   `gorm:"foreignkey:TypeID;association_autoupdate:false;association_autocreate:false"`
	TypeID uint   //Gorm's variable for auto Preload function
}
