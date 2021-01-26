package controllers

import (
	"fmt"
	"gin/app/models"
	"gin/app/search"
	"time"

	"github.com/revel/revel"
)

func initializeDB() {
	InitDB()
	user := models.User{}

	result := DB.First(&user)

	if result.Error != nil {
		role := models.Role{}
		role.ID = 1
		DB.First(&role)

		firstUser := models.User{Name: "Alexander", LastName: "De Carlo", Email: "alexander.decarlo@gmail.com", Active: "TRUE"}
		firstUser.SetNewPassword("alex")
		firstUser.CreatedAt = time.Now()
		firstUser.UpdatedAt = time.Now()

		firstUser.Role = role

		DB.Create(&firstUser)

		firstUser = models.User{Name: "Alexander", LastName: "De Carlo", Email: "a.decarlo@innovationengineering.eu", Active: "TRUE"}
		firstUser.SetNewPassword("alex")
		firstUser.CreatedAt = time.Now()
		firstUser.UpdatedAt = time.Now()

		role.ID = 2
		DB.First(&role)

		firstUser.Role = role

		DB.Create(&firstUser)
	}
}

func closeDBConnection() {
	//defer DB.Close()
}

func init() {
	revel.OnAppStart(initializeDB)
	revel.OnAppStart(inizializeIndexCore)
	revel.InterceptMethod(Dashboard.checkUser, revel.BEFORE)
	revel.OnAppStop(closeDBConnection)
}

func inizializeIndexCore() {
	indexName, _ := revel.Config.String("bleve.index.name")

	search.CreateOrOpenIndex(indexName)

	actualCount, err := search.Index.DocCount()

	if err != nil {
		fmt.Println(err)
	}

	if actualCount == 0 {
		revel.AppLog.Info("Adding all entities to the main index:", indexName, "!")

		entityInstances := []models.EntityInstance{}
		DB.Preload("Entity").Preload("NodeInstances").Preload("NodeInstances.Node").Preload("NodeInstances.Node.Type").Find(&entityInstances)

		for _, entityInstance := range entityInstances {
			revel.AppLog.Info("Adding new entity instance.", "Adding ID ", entityInstance.ID)
			search.Index.Index(fmt.Sprint(entityInstance.ID), &entityInstance)
		}

	} else {
		revel.AppLog.Info("Search indexes already initialized!")
	}
}
