package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) AdmHome(g *gin.Context) {
	var videoCount, actorCount, channelCount, userCount, categoryCount int64
	c.datastore.Table("videos").Where(NOT_DELETED).Count(&videoCount)
	c.datastore.Table("actors").Where(NOT_DELETED).Count(&actorCount)
	c.datastore.Table("channels").Where(NOT_DELETED).Count(&channelCount)
	c.datastore.Table("users").Where(NOT_DELETED).Count(&userCount)
	c.datastore.Table("categories").Where(NOT_DELETED).Count(&categoryCount)
	binaryPath, _ := os.Executable()
	workingDirectory, _ := os.Getwd()
	g.JSON(http.StatusOK, gin.H{
		"build":             c.build,
		"video_count":       videoCount,
		"actor_count":       actorCount,
		"channel_count":     channelCount,
		"user_count":        userCount,
		"category_count":    categoryCount,
		"golang_version":    runtime.Version(),
		"db_driver":         c.config.DB.Driver,
		"binary_path":       binaryPath,
		"startup_directory": workingDirectory,
		"health_errors":     c.healthError,
	})
}

func (c *Controller) AdmVideoList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

func (c *Controller) AdmActorList(g *gin.Context) {
	var actors []model.Actor
	c.datastore.Find(&actors)
	g.JSON(http.StatusOK, gin.H{"items": actors, "total": len(actors)})
}

func (c *Controller) AdmChannelList(g *gin.Context) {
	var channels []model.Channel
	c.datastore.Find(&channels)
	g.JSON(http.StatusOK, gin.H{"items": channels, "total": len(channels)})
}

func (c *Controller) AdmCategory(g *gin.Context) {
	var categories []model.Category
	result := c.datastore.Preload("Sub").Find(&categories)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusOK, gin.H{"items": []model.Category{}, "total": 0})
		return
	}
	g.JSON(http.StatusOK, gin.H{"items": categories, "total": len(categories)})
}

func (c *Controller) AdmTaskList(g *gin.Context) {
	var tasks []model.Task
	c.datastore.Find(&tasks)
	g.JSON(http.StatusOK, gin.H{"items": tasks, "total": len(tasks)})
}

func (c *Controller) AdmTaskHome(g *gin.Context) {
	var tasks []model.Task
	c.datastore.Limit(5).Order("created_at DESC").Find(&tasks)
	g.JSON(http.StatusOK, gin.H{"items": tasks, "total": len(tasks)})
}

func (c *Controller) AdmTaskView(g *gin.Context) {
	id := g.Param("id")
	task := &model.Task{ID: id}
	if result := c.datastore.First(task); result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	g.JSON(http.StatusOK, task)
}

func (c *Controller) AdmTaskRetry(g *gin.Context) {
	id := g.Param("id")
	task := &model.Task{ID: id}
	if result := c.datastore.First(task); result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	task.Status = model.TaskStatusTodo
	if err := c.datastore.Save(task).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.runner.TaskRetry(task.Name)
	g.JSON(http.StatusOK, gin.H{"redirect": "/adm/task/" + task.ID})
}

func (c *Controller) AdmUserList(g *gin.Context) {
	var users []model.User
	c.datastore.Find(&users)
	g.JSON(http.StatusOK, gin.H{"items": users, "total": len(users)})
}

func (c *Controller) AdmUserNew(g *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Password == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
		return
	}
	var existing model.User
	if result := c.datastore.First(&existing, "username = ?", body.Username); result.RowsAffected > 0 {
		g.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}
	passwordHex := sha256.Sum256([]byte(body.Password))
	newUser := &model.User{
		Username:  body.Username,
		Admin:     body.Admin,
		Password:  hex.EncodeToString(passwordHex[:]),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.datastore.Create(newUser).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"id": newUser.ID, "redirect": "/adm/users"})
}

func (c *Controller) AdmUserDelete(g *gin.Context) {
	id := g.Param("id")
	user := model.User{ID: id}
	if result := c.datastore.First(&user); result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.datastore.Delete(&user).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusNoContent, gin.H{})
}

func (c *Controller) AdmConfigAuth(g *gin.Context) {
	dbconfig := &model.Configuration{}
	if result := c.datastore.First(dbconfig); result.RowsAffected < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": ErrConfigAbsent})
		return
	}
	g.JSON(http.StatusOK, gin.H{"authentication_enabled": dbconfig.UserAuthentication})
}

func (c *Controller) AdmConfigAuthUpdate(g *gin.Context) {
	action := g.Param("action")
	if action != "enable" && action != "disable" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}
	dbconfig := &model.Configuration{}
	if result := c.datastore.First(dbconfig); result.RowsAffected < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": ErrConfigAbsent})
		return
	}
	dbconfig.UserAuthentication = action == "enable"
	if err := c.datastore.Save(dbconfig).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.ConfigurationFromDBApply(dbconfig)
	g.JSON(http.StatusOK, gin.H{"redirect": "/adm/config/auth"})
}

func (c *Controller) AdmConfigProvider(g *gin.Context) {
	var providers []model.Provider
	c.datastore.Find(&providers)
	dbconfig := &model.Configuration{}
	if result := c.datastore.First(dbconfig); result.RowsAffected < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": ErrConfigAbsent})
		return
	}
	providerLoaded := make(map[string]string)
	for k, p := range c.providers {
		providerLoaded[k] = p.NiceName()
	}
	g.JSON(http.StatusOK, gin.H{
		"providers":       providers,
		"provider_loaded": providerLoaded,
		"offline_mode":    dbconfig.OfflineMode,
	})
}

func (c *Controller) AdmConfigProviderSwitch(g *gin.Context) {
	providerID := g.Param("id")
	provider := model.Provider{ID: providerID}
	if result := c.datastore.First(&provider); result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}
	provider.Enabled = !provider.Enabled
	if err := c.datastore.Save(&provider).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"redirect": "/adm/config/provider"})
}

func (c *Controller) AdmConfigOfflineMode(g *gin.Context) {
	dbconfig := &model.Configuration{}
	if result := c.datastore.First(dbconfig); result.RowsAffected < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": ErrConfigAbsent})
		return
	}
	g.JSON(http.StatusOK, gin.H{"offline_mode": dbconfig.OfflineMode})
}

func (c *Controller) AdmConfigOfflineModeUpdate(g *gin.Context) {
	action := g.Param("action")
	if action != "enable" && action != "disable" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}
	dbconfig := &model.Configuration{}
	if result := c.datastore.First(dbconfig); result.RowsAffected < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": ErrConfigAbsent})
		return
	}
	dbconfig.OfflineMode = action == "enable"
	if err := c.datastore.Save(dbconfig).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.ConfigurationFromDBApply(dbconfig)
	g.JSON(http.StatusOK, gin.H{"redirect": "/adm/config/offline"})
}
