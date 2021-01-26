package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Plugin model entity
type Plugin struct {
	gorm.Model
	Uuid             string `gorm:"size:36"`
	Name             string `gorm:"size:45"`
	CleanName        string `gorm:"size:45"`
	Description      string `gorm:"type:longtext"`
	Config           string `gorm:"type:longtext"`
	InstallationDate time.Time
	Status           PluginStatus `gorm:"foreignkey:StatusID;association_autoupdate:false;association_autocreate:false"`
	StatusID         uint         //Gorm's variable for auto Preload function
}
