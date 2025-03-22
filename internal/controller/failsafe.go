package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
)

type ConfigNewForm struct {
	Bind         string `form:"bind"`
	Media        string `form:"media-path"`
	DBDriver     string `form:"db-driver"`
	DBConnString string `form:"db-connstring"`
}

func (c *Controller) FailsafeConfiguration(g *gin.Context) {
	var err error

	if g.Request.Method == "POST" {
		var form ConfigNewForm
		err = g.ShouldBind(&form)
		if err == nil {
			newConfig := &config.Config{}
			newConfig.Server.Bind = form.Bind
			newConfig.Media.Path = form.Media
			newConfig.DB.Driver = form.DBDriver
			newConfig.DB.Connstring = form.DBConnString
			file, err := os.OpenFile("config.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err == nil {
				defer file.Close()

				encoder := yaml.NewEncoder(file)

				err = encoder.Encode(newConfig)

				if err == nil {
					c.Restart()
					g.Redirect(http.StatusFound, "/")
					return
				}
			}
		}
	}

	g.HTML(http.StatusOK, "failsafe/config.html", gin.H{
		"Error": err,
	})
}

type UserNewForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (c *Controller) FailsafeUser(g *gin.Context) {
	var err error

	if g.Request.Method == "POST" {
		var form UserNewForm
		err = g.ShouldBind(&form)
		if err == nil {
			if form.Password == "" {
				err = errors.New("password cannot be empty")
			} else {
				now := time.Now()
				passwordHex := sha256.Sum256([]byte(form.Password))
				password := hex.EncodeToString(passwordHex[:])
				admin := &model.User{
					Username:  form.Username,
					Admin:     true,
					CreatedAt: now,
					UpdatedAt: now,
					Password:  password,
				}

				err = c.datastore.Create(&admin).Error
				if err == nil {
					c.Restart()
					g.Redirect(http.StatusFound, "/")
					return
				}
			}
		}
	}

	g.HTML(http.StatusOK, "failsafe/user.html", gin.H{
		"Error": err,
	})
}
