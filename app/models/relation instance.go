package models

import (
	"github.com/jinzhu/gorm"
)

//RelationInstance model
type RelationInstance struct {
	gorm.Model
	ParentEntityInstance   EntityInstance `gorm:"foreignkey:ParentEntityInstanceID;association_autoupdate:false;association_autocreate:false"`
	ChildEntityInstance    EntityInstance `gorm:"foreignkey:ChildEntityInstanceID;association_autoupdate:false;association_autocreate:false"`
	ParentEntityInstanceID uint           //Gorm's variable for auto Preload function
	ChildEntityInstanceID  uint           //Gorm's variable for auto Preload function
}
