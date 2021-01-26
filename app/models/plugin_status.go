package models

//PluginStatus model
type PluginStatus struct {
	ID          uint   `gorm:"primary_key"`
	Description string `gorm:"size:45"`
}
