package models

//Type model
type Type struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"size:255"`
}
