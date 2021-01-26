package models

import (
	"github.com/jinzhu/gorm"
)

//Relation model
type Relation struct {
	gorm.Model
	ParentEntity   Entity `gorm:"foreignkey:ParentEntityID;association_autoupdate:false;association_autocreate:false"`
	ChildEntity    Entity `gorm:"foreignkey:ChildEntityID;association_autoupdate:false;association_autocreate:false"`
	ParentEntityID uint   //Gorm's variable for auto Preload function
	ChildEntityID  uint   //Gorm's variable for auto Preload function
}
