package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) AdmHome(g *gin.Context) {
	// get counts
	var (
		videoCount    int64
		actorCount    int64
		channelCount  int64
		userCount     int64
		categoryCount int64
	)

	c.datastore.Table("videos").Where("deleted_at is null").Count(&videoCount)
	c.datastore.Table("actors").Where("deleted_at is null").Count(&actorCount)
	c.datastore.Table("channels").Where("deleted_at is null").Count(&channelCount)
	c.datastore.Table("users").Where("deleted_at is null").Count(&userCount)
	c.datastore.Table("categories").Where("deleted_at is null").Count(&categoryCount)

	binaryPath, err := os.Executable()
	if err != nil {
		c.logger.Warn().Err(err).Msg("unable to get binary path")
		binaryPath = ""
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		c.logger.Warn().Err(err).Msg("unable to get binary working directory")
		workingDirectory = ""
	}

	c.HTML(g, http.StatusOK, "adm/home.html", gin.H{
		"Build":            c.build,
		"VideoCount":       videoCount,
		"ActorCount":       actorCount,
		"ChannelCount":     channelCount,
		"UserCount":        userCount,
		"CategoryCount":    categoryCount,
		"GolangVersion":    runtime.Version(),
		"DBDriver":         c.config.DB.Driver,
		"BinaryPath":       binaryPath,
		"StartupDirectory": workingDirectory,
		"HealthErrors":     c.healthError,
	})
}

func (c *Controller) AdmVideoList(g *gin.Context) {
	var videos []model.Video

	c.datastore.Find(&videos)

	c.HTML(g, http.StatusOK, "adm/object-list.html", gin.H{
		"ObjectName": "Video",
		"Objects":    videos,
	})
}

func (c *Controller) AdmActorList(g *gin.Context) {
	var actors []model.Actor

	c.datastore.Find(&actors)

	c.HTML(g, http.StatusOK, "adm/object-list.html", gin.H{
		"ObjectName": "Actor",
		"Objects":    actors,
	})
}

func (c *Controller) AdmChannelList(g *gin.Context) {
	var channels []model.Channel

	c.datastore.Find(&channels)

	c.HTML(g, http.StatusOK, "adm/object-list.html", gin.H{
		"ObjectName": "Channel",
		"Objects":    channels,
	})
}

func (c *Controller) AdmTaskView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	task := &model.Task{
		ID: id,
	}
	result := c.datastore.First(task)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	c.HTML(g, http.StatusOK, "adm/task-view.html", gin.H{
		"Task": task,
	})
}

func (c *Controller) AdmTaskList(g *gin.Context) {
	tasks := []model.Task{}
	c.datastore.Find(&tasks)

	c.HTML(g, http.StatusOK, "adm/task-list.html", gin.H{
		"Tasks": tasks,
	})
}

func (c *Controller) AdmTaskRetry(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	task := &model.Task{
		ID: id,
	}
	result := c.datastore.First(task)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	task.Status = model.TaskStatusTodo

	err := c.datastore.Save(task).Error
	if err != nil {
		c.ErrFatal(g, err.Error())
		return
	}

	c.runner.TaskRetry(task.Name)

	taskURL := fmt.Sprintf("/adm/task/%s", task.ID)
	g.Redirect(http.StatusFound, taskURL)
}

func (c *Controller) AdmUserList(g *gin.Context) {
	var users []model.User

	c.datastore.Find(&users)

	c.HTML(g, http.StatusOK, "adm/user-list.html", gin.H{
		"Objects": users,
	})
}

type AdmUserNewForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Admin    string `form:"admin"`
}

func (c *Controller) AdmUserNew(g *gin.Context) {
	var err error

	if g.Request.Method == "POST" {
		var form AdmUserNewForm
		err = g.ShouldBind(&form)
		if err == nil {
			if form.Password == "" {
				err = errors.New("password cannot be empty")
			} else {
				// get user
				userExists := &model.User{}
				result := c.datastore.First(userExists, "username = ?", form.Username)
				if result.RowsAffected > 0 {
					err = errors.New("username already taken")
				} else {
					now := time.Now()
					passwordHex := sha256.Sum256([]byte(form.Password))
					password := hex.EncodeToString(passwordHex[:])
					newUser := &model.User{
						Username:  form.Username,
						Admin:     form.Admin != "",
						CreatedAt: now,
						UpdatedAt: now,
						Password:  password,
					}

					err = c.datastore.Create(&newUser).Error
					if err == nil {
						g.Redirect(http.StatusFound, "/adm/users")
						return
					}
				}
			}
		}
	}

	c.HTML(g, http.StatusOK, "adm/user-new.html", gin.H{
		"Error": err,
	})
}

func (c *Controller) AdmUserDelete(g *gin.Context) {
	var err error

	// get alias id from path
	userID := g.Param("id")

	user := model.User{
		ID: userID,
	}
	result := c.datastore.First(&user)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	err = c.datastore.Delete(&user).Error
	if err != nil {
		c.ErrFatal(g, err.Error())
		return
	}

	g.Redirect(http.StatusFound, "/adm/users")
}

func (c *Controller) AdmCategory(g *gin.Context) {
	categories := []model.Category{}
	result := c.datastore.Preload("Sub").Find(&categories)
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	c.HTML(g, http.StatusOK, "adm/category.html", gin.H{
		"Categories": categories,
	})
}

func (c *Controller) AdmTaskHome(g *gin.Context) {
	var tasks []model.Task
	c.datastore.Limit(5).Order("created_at DESC").Find(&tasks)

	c.HTML(g, http.StatusOK, "adm/task-home.html", gin.H{
		"Tasks": tasks,
	})
}

func (c *Controller) AdmConfigAuth(g *gin.Context) {
	c.HTML(g, http.StatusOK, "adm/config-auth.html", gin.H{})
}

func (c *Controller) AdmConfigAuthUpdate(g *gin.Context) {
	action := g.Param("action")
	if action != "enable" && action != "disable" {
		g.Redirect(http.StatusFound, "/adm/config/auth")
		return
	}

	// loading configuration from database
	dbconfig := &model.Configuration{}
	result := c.datastore.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		c.HTML(g, http.StatusOK, "adm/config-auth.html", gin.H{
			"error": errors.New("configuration not found"),
		})
		return
	}

	dbconfig.UserAuthentication = action == "enable"
	err := c.datastore.Save(&dbconfig).Error
	if result.RowsAffected < 1 {
		c.HTML(g, http.StatusOK, "adm/config-auth.html", gin.H{
			"error": err,
		})
		return
	}

	c.ConfigurationFromDBApply(dbconfig)
	g.Redirect(http.StatusFound, "/adm/config/auth")
}

func (c *Controller) AdmConfigProvider(g *gin.Context) {
	// find providers
	var providers []model.Provider
	c.datastore.Find(&providers)

	// loading configuration from database
	dbconfig := &model.Configuration{}
	result := c.datastore.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		c.HTML(g, 500, "err/fatal.html", gin.H{
			"error": "configuration not found, restarting the appliaction should fix the issue",
		})
		return
	}

	c.HTML(g, http.StatusOK, "adm/config-provider.html", gin.H{
		"Providers":      providers,
		"ProviderLoaded": c.providers,
		"OfflineMode":    dbconfig.OfflineMode,
	})
}

func (c *Controller) AdmConfigProviderSwitch(g *gin.Context) {
	providerID := g.Param("id")

	// find providers
	provider := model.Provider{
		ID: providerID,
	}

	result := c.datastore.First(&provider)

	// check result
	if result.RowsAffected < 1 {
		c.HTML(g, 404, "err/fatal.html", gin.H{
			"error": "provider not found",
		})
		return
	}

	provider.Enabled = !provider.Enabled
	err := c.datastore.Save(&provider).Error
	// check result
	if err != nil {
		c.HTML(g, 500, "err/fatal.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	g.Redirect(http.StatusFound, "/adm/config/provider")
}

func (c *Controller) AdmConfigOfflineMode(g *gin.Context) {
	// loading configuration from database
	dbconfig := &model.Configuration{}
	result := c.datastore.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		c.HTML(g, 500, "err/fatal.html", gin.H{
			"error": "configuration not found, restarting the appliaction should fix the issue",
		})
		return
	}

	c.HTML(g, http.StatusOK, "adm/config-offline.html", gin.H{
		"OfflineMode": dbconfig.OfflineMode,
	})
}

func (c *Controller) AdmConfigOfflineModeUpdate(g *gin.Context) {
	action := g.Param("action")
	if action != "enable" && action != "disable" {
		g.Redirect(http.StatusFound, "/adm/config/offline")
		return
	}

	// loading configuration from database
	dbconfig := &model.Configuration{}
	result := c.datastore.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		c.HTML(g, 500, "err/fatal.html", gin.H{
			"error": "configuration not found, restarting the appliaction should fix the issue",
		})
		return
	}

	dbconfig.OfflineMode = action == "enable"
	err := c.datastore.Save(&dbconfig).Error
	if result.RowsAffected < 1 {
		c.HTML(g, http.StatusOK, "adm/config-offline.html", gin.H{
			"error": err,
		})
		return
	}

	c.ConfigurationFromDBApply(dbconfig)
	g.Redirect(http.StatusFound, "/adm/config/offline")
}
