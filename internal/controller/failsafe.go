package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/zobtube/zobtube/internal/config"
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
					//TODO: improve with http server kill and self execve
					fmt.Println("configuration set, exiting to allow a restart")
					os.Exit(0)
				}
			}
		}
	}

	g.HTML(http.StatusOK, "failsafe/config.html", gin.H{
		"Error": err,
	})
}
