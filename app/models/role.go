package models

import (
	"github.com/jinzhu/gorm"
)

//Role model
type Role struct {
	gorm.Model
	Description string `gorm:"size:45"`
}
