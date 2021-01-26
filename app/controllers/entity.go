package controllers

import (
	"fmt"
	"gin/app/models"
	"gin/app/search"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/revel/revel"
)

//Entity ...
type Entity struct {
	*revel.Controller
}

//Entity ...
func (c Entity) Entity() revel.Result {
	return c.Render()
}

func (c App) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(App.Index)
	}
	return nil
}

//Index ...
func (c Entity) Index() revel.Result {
	//var entities = []models.Entity{}

	//DB.Find(&entities)
	search, size, page, nextPage := "", 10, 1, 2

	c.ViewArgs["entities"] = List(search, size, page)
	c.ViewArgs["search"] = search
	c.ViewArgs["size"] = size
	c.ViewArgs["page"] = page
	c.ViewArgs["nextPage"] = nextPage

	return c.Render()
}

//List ...
func List(search string, size int, page int) []models.Entity {
	if page == 0 {
		page = 1
	}

	search = strings.TrimSpace(search)

	var entities = []models.Entity{}

	where := ""

	if search != "" {
		where = "%" + strings.ToLower(search) + "%"
		where, _, _ = squirrel.Or{squirrel.Expr("lower(Name) like ?", search), squirrel.Expr("lower(Description) like ?", search)}.ToSql()
	}

	DB.Where(where).Find(&entities).Offset((page - 1) * size).Limit(size)

	return entities
}

//View ...
func (c Entity) View(ID int) revel.Result {
	entity := &models.Entity{}

	DB.Where("ID = ?", ID).Find(&entity)
	c.ViewArgs["entity"] = &entity

	return c.Render()
}

//Edit ...
func (c Entity) Edit(ID int) revel.Result {
	entity := &models.Entity{}

	if ID != 0 {
		DB.Where("ID = ?", ID).Find(&entity)
	}

	c.ViewArgs["entity"] = &entity

	return c.Render()
}

//SaveEntity ...
func (c Entity) SaveEntity(entity *models.Entity) revel.Result {
	if entity.ID != 0 {
		DB.Save(&entity)
	} else {
		DB.Create(&entity)
	}

	search.Index.Index(fmt.Sprint(&entity.ID), entity)

	return c.Redirect(Entity.Index)
}

//Nodes ...
func (c Entity) Nodes(ID int) revel.Result {
	entity := &models.Entity{}
	DB.Where("ID = ?", ID).Preload("Nodes").Preload("Nodes.Type").Find(&entity)

	c.ViewArgs["nodes"] = entity.Nodes
	c.ViewArgs["ID"] = ID

	types := []models.Type{}
	DB.Find(&types)
	c.ViewArgs["types"] = types

	return c.Render()
}

//AddNode ...
func (c Entity) AddNode(ID int, node models.Node) revel.Result {
	entity := &models.Entity{}
	DB.Where("ID = ?", ID).Preload("Nodes").Find(&entity)
	entity.Nodes = append(entity.Nodes, node)

	DB.Save(&entity)

	return c.Redirect(Entity.Nodes, ID)
}

//DeleteNode ...
func (c Entity) DeleteNode(ID int, entityID int) revel.Result {
	node := &models.Node{}
	DB.Where("ID = ?", ID).Delete(&node)

	return c.Redirect(Entity.Nodes, entityID)
}

//Delete ...
func (c Entity) Delete(ID int) revel.Result {
	entity := &models.Entity{}
	DB.Where("ID = ?", ID).Delete(&entity)
	search.Index.Delete(fmt.Sprint(ID))

	return c.Redirect(Entity.Index)
}

//All ...
func (c Entity) All() revel.Result {
	entities := []models.Entity{}
	DB.Find(&entities)

	response := models.JSONResponse{}
	response.Data = entities

	return c.RenderJSON(response)
}

//SearchEntity ...
func (c Entity) SearchEntity(searchText string, from int) revel.Result {
	//results := search.Search(searchText+"*", []string{"Entity.Name", "Entity.Description"}, from, 0)
	entities := &[]models.Entity{}
	DB.Offset(from).Find(&entities)

	c.ViewArgs["results"] = &entities

	return c.RenderTemplate("entity/search.html")
}

//Search ...
func (c Entity) Search() revel.Result {
	return c.Render()
}
