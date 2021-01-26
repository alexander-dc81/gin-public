package models

import (
	"github.com/jinzhu/gorm"
)

//Entity model
type Entity struct {
	gorm.Model
	Name        string `gorm:"size:255"`
	Description string `gorm:"type:longtext"`
	Nodes       []Node `gorm:"many2many:entity_node"`
}
