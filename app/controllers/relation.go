package controllers

import (
	"fmt"
	"gin/app/models"
	"gin/app/search"
	"strings"

	"github.com/revel/revel"
)

//Relation ...
type Relation struct {
	*revel.Controller
}

//List ...
func (c Relation) List(entityID int) revel.Result {
	relations := []models.Relation{}
	DB.Where("parent_entity_id = ?", entityID).Find(&relations)

	c.ViewArgs["relations"] = relations
	c.ViewArgs["ID"] = entityID

	entities := []models.Entity{}
	DB.Where("ID <> ?", entityID).Find(&entities)
	c.ViewArgs["entities"] = entities

	fmt.Println(len(entities))

	return c.Render()
}

//Add ...
func (c Relation) Add(ID int, relation models.Relation) revel.Result {
	fmt.Println("Inizio Salvataggio Relation")

	parentEntity := &models.Entity{}
	DB.Where("ID = ?", ID).Find(&parentEntity)

	childEntity := &models.Entity{}
	DB.Where("ID = ?", relation.ChildEntityID).Find(&childEntity)

	relation.ParentEntity = *parentEntity
	relation.ChildEntity = *childEntity

	DB.Save(&relation)

	fmt.Println("Fine Salvataggio Relation")
	return c.Redirect(Relation.List, ID)
}

//Delete ...
func (c Relation) Delete(ID int, entityID int) revel.Result {
	relation := &models.Relation{}
	DB.Where("ID = ?", ID).Delete(&relation)

	return c.Redirect(Relation.List, entityID)
}

//EntityInstanceRelation ...
func EntityInstanceRelation(relationID uint) []models.RelationInstance {
	relationInstances := []models.RelationInstance{}
	DB.Where("relation_id = ?", relationID).Preload("ChildEntityInstance.NodeInstances").Preload("ChildEntityInstance.NodeInstances.Node").Find(&relationInstances)

	return relationInstances
}

//Search ...
func (c Relation) Search(searchText string, from int, entityID int) revel.Result {
	fmt.Println("Testo di ricerca: ", searchText)

	//entityInstances := []models.EntityInstance{}

	searchText = strings.TrimSpace(searchText)

	if len(searchText) > 0 {
		searchText = " +NodeInstances.Value:\"" + searchText + "*\""
	}

	searchText = "+Entity.ID:" + fmt.Sprint(entityID) + searchText

	fmt.Println(searchText)
	results := search.Search(searchText, []string{"NodeInstances.Node.ID", "NodeInstances.Node.Name", "NodeInstances.Value"}, from, 0)

	//c.ViewArgs["results"] = results.Hits

	response := models.JSONResponse{}
	response.Data = results.Hits

	return c.RenderJSON(response)
}
