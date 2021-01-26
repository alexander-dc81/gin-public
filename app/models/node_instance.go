package models

import (
	"github.com/jinzhu/gorm"
)

//NodeInstance model
type NodeInstance struct {
	gorm.Model
	Node             Node   `gorm:"foreignkey:NodeID;association_autoupdate:false;association_autocreate:false"`
	Value            string `gorm:"type:text"`
	NodeID           uint   //Gorm's variable for auto Preload function
	EntityInstanceID uint
}
