package controllers

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Forced Mysql driver
	"github.com/revel/revel"
)

//DB ...
var DB *gorm.DB

//InitDB ...
func InitDB() {
	dbInfo, _ := revel.Config.String("db.info")
	logmode, _ := revel.Config.Bool("db.logmode")

	db, err := gorm.Open("mysql", dbInfo)

	if err != nil {
		log.Panicf("Failed gorm.Open: %v\n", err)
	}

	db.SingularTable(true)
	db.LogMode(logmode)
	db.DB()

	DB = db
}
