package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

//EntityInstance model
type EntityInstance struct {
	gorm.Model
	UUID          uuid.UUID `gorm:"size:36"`
	Entity        Entity    `gorm:"foreignkey:EntityID;association_autoupdate:false;association_autocreate:false"`
	NodeInstances []NodeInstance
	EntityID      uint //Gorm's variable for auto Preload function
}
