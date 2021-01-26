package controllers

import (
	"encoding/json"
	"gin/app/models"
	"gin/app/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/revel/revel"
)

//Plugin ...
type Plugin struct {
	*revel.Controller
}

//KB ...
const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

//Index ...
func (c Plugin) Index() revel.Result {
	plugins := &[]models.Plugin{}

	DB.Preload("Status").Find(&plugins)

	c.ViewArgs["plugins"] = *plugins

	return c.Render()
}

//Upload ...
func (c Plugin) Upload(file []byte) revel.Result {
	// Validation rules.
	c.Validation.Required(file)
	c.Validation.MinSize(file, 1*KB).MessageKey("errors.minfilesize")
	c.Validation.MaxSize(file, 10*MB).MessageKey("errors.maxfilesize")

	// Handle errors.
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Plugin.Index)
	}

	filenameOriginal := c.Params.Files["file"][0].Filename
	filenameAndPath := getRootFolder() + filenameOriginal

	err := ioutil.WriteFile(filenameAndPath, file, 0644)

	if err != nil {
		c.Validation.Error("File Upload Error").Message(err.Error())
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Plugin.Index)
	}

	unzip(filenameAndPath, getRootFolder())
	setup(c, filenameAndPath, filenameOriginal)

	c.Flash.Success("Plugin uploaded successfully!")

	return c.Redirect(Plugin.Index)
}

func unzip(filename string, outputFolder string) {
	files, err := utils.Unzip(filename, outputFolder)

	if err != nil {
		log.Fatal(err)
	}

	revel.AppLog.Info("Unzipped plugin file!", "Files:", files)
}

func getRootFolder() string {
	pluginRootFolder, _ := revel.Config.String("plugin.folder")
	path, err := os.Getwd()

	check(err)

	return path + pluginRootFolder
}

func setup(c Plugin, pluginPath string, pluginName string) {
	pluginPath = strings.Replace(pluginPath, ".zip", "", -1)
	data, err := ioutil.ReadFile(pluginPath + "/plugin.json")

	check(err)

	var configJSONInterface interface{}
	err = json.Unmarshal(data, &configJSONInterface)
	check(err)

	pluginEntity := *convertJSONToEntity(configJSONInterface, string(data))
	pluginStatus := &models.PluginStatus{}

	DB.Where("ID = ?", 3).Find(&pluginStatus)

	pluginEntity.Status = *pluginStatus
	pluginEntity.StatusID = pluginStatus.ID

	jsonMap := configJSONInterface.(map[string]interface{})

	elements := jsonMap["install_elements"].([]interface{})

	//List of possible folder to move: [html css js controller routing configuration]
	for _, element := range elements {
		switch element {
		case "html":
			moveHTMLFiles(pluginPath, pluginEntity.Uuid)
		case "css":
			moveCSSFiles(pluginPath)
		case "js":
			moveJSFiles(pluginPath)
		case "controller":
			moveControllerFiles(pluginPath)
		case "routing":
			pluginStatus = &models.PluginStatus{}
			DB.Where("ID = ?", 4).Find(&pluginStatus)
			pluginEntity.Status = *pluginStatus
			pluginEntity.StatusID = pluginStatus.ID
		case "configuration":
			initializePluginConfiguration(pluginPath, pluginEntity.Uuid)
		default:
		}
	}

	DB.Create(&pluginEntity)
}

func convertJSONToEntity(configJSONInterface interface{}, dataConfig string) *models.Plugin {
	jsonMap := configJSONInterface.(map[string]interface{})

	pluginEntity := models.Plugin{Name: (jsonMap["name"]).(string), Description: (jsonMap["description"]).(string), Config: dataConfig,
		Uuid: jsonMap["uuid"].(string), CleanName: jsonMap["clean_name"].(string), InstallationDate: time.Now()}

	return &pluginEntity
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func createFolderIfNotExists(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	check(err)
}

func moveHTMLFiles(pluginPath string, pluginUUID string) {
	path, err := os.Getwd()

	check(err)

	createFolderIfNotExists(path + "/app/views/plugins/html/" + pluginUUID)

	err = utils.CopyDirectory(pluginPath+"/views/html/", path+"/app/views/plugins/html/")

	check(err)
}

func moveJSFiles(pluginPath string) {
	path, err := os.Getwd()

	check(err)

	err = utils.CopyDirectory(pluginPath+"/views/js/", path+"/app/views/plugins/js/")

	check(err)
}

func moveCSSFiles(pluginPath string) {
	path, err := os.Getwd()

	check(err)

	err = utils.CopyDirectory(pluginPath+"/views/css/", path+"/app/views/plugins/css/")

	check(err)
}

func moveControllerFiles(pluginPath string) {
	path, err := os.Getwd()

	check(err)

	err = utils.CopyDirectory(pluginPath+"/controllers/", path+"/app/controllers/plugins/")

	check(err)
}

func initializePluginConfiguration(pluginPath string, uuid string) {
	pluginConfFile, err := ioutil.ReadFile(pluginPath + "/configuration/plugin.conf")

	check(err)

	path, err := os.Getwd()

	check(err)

	mainConfFile, err := os.OpenFile(path+"/conf/app.conf", os.O_APPEND, 0644)

	check(err)
	_, err = mainConfFile.WriteString("\n###" + uuid + "###\n")
	_, err = mainConfFile.WriteString(string(pluginConfFile) + "\n")
	_, err = mainConfFile.WriteString("###" + uuid + "###")
}

func initializePluginRoutes(pluginPath string, uuid string) {
	pluginRoutesFile, err := ioutil.ReadFile(pluginPath + "/routing/routes")

	check(err)

	path, err := os.Getwd()

	check(err)

	mainRoutesFile, err := os.OpenFile(path+"/conf/routes", os.O_APPEND, 0644)

	check(err)
	_, err = mainRoutesFile.WriteString("\n###" + uuid + "###\n")
	_, err = mainRoutesFile.WriteString(string(pluginRoutesFile) + "\n")
	_, err = mainRoutesFile.WriteString("###" + uuid + "###")

	//Refresh main routing file after installing new pages and controllers
	(*revel.Router).Refresh(revel.MainRouter)
}

//InitRoutes ...
func (c Plugin) InitRoutes(ID uint) revel.Result {
	plugin := &models.Plugin{}

	DB.Where("ID = ?", ID).Find(&plugin)

	pluginPath := getRootFolder() + plugin.CleanName

	initializePluginRoutes(pluginPath, plugin.Uuid)

	pluginStatus := &models.PluginStatus{}

	DB.Where("ID = ?", 3).Find(&pluginStatus)
	plugin.Status = *pluginStatus
	plugin.StatusID = pluginStatus.ID

	DB.Save(&plugin)

	return c.Redirect(Plugin.Index)
}

//Uninstall ...
func (c Plugin) Uninstall(ID uint) revel.Result {
	plugin := &models.Plugin{}

	DB.Where("ID = ?", ID).Find(&plugin)

	deletePluginRoutes(plugin.Uuid)
	deletePluginControllers(plugin.Uuid)
	deleteJSFile(plugin.Uuid)
	deleteCSSFile(plugin.Uuid)
	deleteHTMLFile(plugin.Uuid)
	deleteDownloadedFile(c, plugin.CleanName)
	deletePluginConf(plugin.Uuid)

	DB.Where("ID = ?", plugin.ID).Delete(&plugin)

	return c.Redirect(Plugin.Index)
}

func deletePluginRoutes(uuid string) {
	path, err := os.Getwd()

	check(err)
	data, err := ioutil.ReadFile(path + "/conf/routes")

	check(err)

	lines := strings.Split(string(data), "\n")

	startLine := 0
	endLine := 0

	for i, line := range lines {
		if strings.Contains(line, "###"+uuid+"###") {
			if startLine == 0 {
				startLine = i
			} else {
				endLine = i
				break
			}
		}
	}

	// Remove all the lines of the plugin from the route file
	for startLine < endLine {
		lines = append(lines[:startLine], lines[startLine+1:]...)
		startLine++
	}

	lines = lines[:len(lines)-1]

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path+"/conf/routes", []byte(output), 0644)
	check(err)
}

func deletePluginConf(uuid string) {
	path, err := os.Getwd()

	check(err)
	data, err := ioutil.ReadFile(path + "/conf/app.conf")

	check(err)

	lines := strings.Split(string(data), "\n")

	startLine := 0
	endLine := 0

	for i, line := range lines {
		if strings.Contains(line, "###"+uuid+"###") {
			if startLine == 0 {
				startLine = i
			} else {
				endLine = i
				break
			}
		}
	}

	// Remove all the lines of the plugin from the route file
	for startLine < endLine {
		lines = append(lines[:startLine], lines[startLine+1:]...)
		startLine++
	}

	lines = lines[:len(lines)-1]

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path+"/conf/app.conf", []byte(output), 0644)
	check(err)
}

func deletePluginControllers(uuid string) {
	path, err := os.Getwd()

	check(err)

	matches := findFile(path+"/app/controllers/plugins/", uuid+"*")

	for _, match := range matches {
		err := os.Remove(match)

		check(err)
	}
}

func deleteJSFile(uuid string) {
	path, err := os.Getwd()

	check(err)

	matches := findFile("/app/views/plugins/js/", uuid+"*")

	for _, match := range matches {
		err = os.Remove(path + match)

		check(err)
	}
}

func deleteCSSFile(uuid string) {
	path, err := os.Getwd()

	check(err)

	matches := findFile("/app/views/plugins/css/", uuid+"*")

	for _, match := range matches {
		err = os.Remove(path + match)

		check(err)
	}
}

func deleteHTMLFile(uuid string) {
	path, err := os.Getwd()

	check(err)

	matches := findFile(path+"/app/views/Plugins/html/side_menu_admin/", uuid+"*")

	for _, match := range matches {
		err = os.Remove(match)

		check(err)
	}

	matches = findFile(path+"/app/views/Plugins/html/side_menu_user/", uuid+"*")

	for _, match := range matches {
		err = os.Remove(match)

		check(err)
	}

	err = os.RemoveAll(path + "/app/views/plugins/html/" + uuid)

	check(err)
}

func deleteDownloadedFile(c Plugin, pluginName string) {
	pluginBasePath := getRootFolder() + strings.Replace(pluginName, ".zip", "", -1)
	err := os.RemoveAll(pluginBasePath)

	if err != nil {
		log.Fatal(err)
		c.Validation.Error("Plugin folder deletion error.").Message(err.Error())
		c.Validation.Keep()
		c.FlashParams()
		c.Redirect(Plugin.Index)
	}

	err = os.Remove(pluginBasePath + ".zip")

	if err != nil {
		log.Fatal(err)
		c.Validation.Error("Plugin file deletion error.").Message(err.Error())
		c.Validation.Keep()
		c.FlashParams()
		c.Redirect(Plugin.Index)
	}
}

func findFile(targetDir string, pattern string) []string {
	matches, err := filepath.Glob(targetDir + pattern)

	check(err)

	return matches
}
