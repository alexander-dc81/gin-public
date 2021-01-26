package controllers

import (
	"fmt"
	"gin/app/models"
	"gin/app/search"
	"reflect"
	"strings"
	"time"

	"github.com/revel/revel"

	"github.com/google/uuid"
)

//Instance ...
type Instance struct {
	*revel.Controller
}

//Dashboard ...
func (c Instance) Dashboard(ID int) revel.Result {
	entityInstances := &[]models.EntityInstance{}

	currentTime := time.Now()
	currentTime = currentTime.AddDate(0, -1, 0)

	DB.Where("entity_id = ? AND created_at > ?", ID, currentTime).Find(&entityInstances)
	c.ViewArgs["entity_instances_created"] = len(*entityInstances)

	DB.Where("entity_ID = ? AND updated_at > ?", ID, currentTime).Find(&entityInstances)
	c.ViewArgs["entity_instances_updated"] = len(*entityInstances)

	DB.Where("entity_ID = ? AND created_at < ?", ID, currentTime).Find(&entityInstances)
	c.ViewArgs["entity_instances_old"] = len(*entityInstances)

	return c.Render()
}

//Edit ...
func (c Instance) Edit(ID int, entityID int) revel.Result {
	entityInstance := &models.EntityInstance{}
	entity := &models.Entity{}

	if ID != 0 {
		DB.Where("ID = ?", ID).Preload("Entity").Preload("Entity.Nodes").Preload("Entity.Nodes.Type").Preload("NodeInstances.Node").Preload("NodeInstances.Node.Type").Find(&entityInstance)
		generateNodeInstances(entityInstance, &entityInstance.Entity)
	} else {
		DB.Where("ID = ?", entityID).Preload("Nodes").Preload("Nodes.Type").Find(&entity)
		generateNodeInstances(entityInstance, entity)
	}

	relations := []models.Relation{}

	DB.Where("parent_entity_id = ?", entityID).Preload("ChildEntity.Nodes").Find(&relations)

	c.ViewArgs["entity_instance"] = &entityInstance
	c.ViewArgs["relations"] = &relations

	return c.Render()
}

// Generate a node instance based on the Entity type and nodes of the parent Entity
func generateNodeInstances(entityInstance *models.EntityInstance, entity *models.Entity) {
	entityInstance.Entity = *entity

	for _, node := range entity.Nodes {
		nodeInstance := getNodeFromNodeinstanceByID(entityInstance.NodeInstances, node.ID)

		if nodeInstance == nil {
			nodeInstance = &models.NodeInstance{}
			nodeInstance.Node = node
			nodeInstance.NodeID = node.ID
			nodeInstance.EntityInstanceID = entityInstance.ID
			entityInstance.NodeInstances = append(entityInstance.NodeInstances, *nodeInstance)
		}
	}
}

// Search in the slice for the specific node, if exists. nil elsewhere
func getNodeFromNodeinstanceByID(s []models.NodeInstance, ID uint) *models.NodeInstance {
	for _, nodeInstance := range s {
		if nodeInstance.Node.ID == ID {
			return &nodeInstance
		}
	}

	return nil
}

//Save ...
func (c Instance) Save(entityInstance models.EntityInstance) revel.Result {
	if entityInstance.ID != 0 {
		DB.Save(&entityInstance)
	} else {
		entityInstance.UUID = uuid.New()
		DB.Create(&entityInstance)
	}

	search.Index.Index(fmt.Sprint(entityInstance.ID), &entityInstance)

	return c.Redirect(Instance.Dashboard, entityInstance.ID)
}

//Delete ...
func (c Instance) Delete(ID int) revel.Result {
	entityInstance := &models.EntityInstance{}
	DB.Where("ID = ?", ID).Delete(&entityInstance)
	search.Index.Delete(fmt.Sprint(ID))

	return c.Redirect(Instance.Dashboard)
}

//List ...
func (c Instance) List(searchText string, start int, ID int) revel.Result {
	fmt.Println("Ricerca fatta per ID " + fmt.Sprint(ID))
	entity := &models.Entity{}

	searchText = strings.TrimSpace(searchText)

	if len(searchText) > 0 {
		searchText = " +NodeInstances.Value:\"" + searchText + "*\""
	}

	searchText = "+Entity.ID:" + fmt.Sprint(ID) + searchText

	fmt.Println(searchText)
	results := search.Search(searchText, []string{"NodeInstances.Node.ID", "NodeInstances.Node.Name", "NodeInstances.Value"}, start, 0)

	c.ViewArgs["results"] = results.Hits

	DB.Where("ID = ?", ID).Preload("Nodes").Find(&entity)

	c.ViewArgs["ID"] = ID
	c.ViewArgs["entity"] = &entity

	return c.Render()
}

//UUIDBinder for UUIDs in forms
var UUIDBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		if len(val) == 0 {
			return reflect.Zero(typ)
		}
		var uuidValue, err = uuid.Parse(val)

		if err != nil {
			revel.AppLog.Warn("UUIDBinder Conversion Error", "error", err)
			return reflect.Zero(typ)
		}

		return reflect.ValueOf(uuidValue)
	}),

	Unbind: func(output map[string]string, key string, val interface{}) {
		output[key] = fmt.Sprintf("%f", val)
	},
}

// Init the UUID specific binder for bind/unbind the the Revel form
func init() {
	revel.TypeBinders[reflect.TypeOf(uuid.UUID{})] = UUIDBinder
}
