package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model
	Name           string `gorm:"size:255"`
	LastName       string `gorm:"size:255"`
	Email          string `gorm:"type:varchar(100);unique_index"`
	HashedPassword []byte
	Active         string `gorm:"size:5"`
	Role           Role   `gorm:"foreignkey:RoleID;association_autoupdate:false;association_autocreate:false"`
	RoleID         uint   //Gorm's variable for auto Preload function
}

// SetNewPassword set a new hashsed password to user
func (user *User) SetNewPassword(passwordString string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	user.HashedPassword = bcryptPassword
}
