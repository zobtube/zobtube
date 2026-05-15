package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/metamigrate"
)

// AdmMetadataStorage godoc
//
//	@Summary	Get metadata storage configuration (admin, read-only)
//	@Tags		admin
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/adm/metadata-storage [get]
func (c *Controller) AdmMetadataStorage(g *gin.Context) {
	if c.config == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "configuration not available"})
		return
	}
	resp := gin.H{
		"type":    c.config.Metadata.Type,
		"source":  "cli",
		"message": "Metadata storage is configured via CLI flags, environment variables, or config.yml. Restart the server after changing settings.",
	}
	switch c.config.Metadata.Type {
	case "filesystem":
		resp["path"] = c.config.Metadata.Path
	case "s3":
		resp["bucket"] = c.config.Metadata.S3Bucket
		resp["region"] = c.config.Metadata.S3Region
		resp["prefix"] = c.config.Metadata.S3Prefix
		resp["endpoint"] = c.config.Metadata.S3Endpoint
		if c.config.Metadata.S3AccessKeyID != "" {
			resp["access_key_id"] = c.config.Metadata.S3AccessKeyID
			resp["access_key_configured"] = true
		}
		if c.config.Metadata.S3SecretAccessKey != "" {
			resp["secret_access_key_configured"] = true
		}
	}
	g.JSON(http.StatusOK, resp)
}

// AdmMetadataStorageMigrate godoc
//
//	@Summary	Enqueue metadata thumbnail migration (admin)
//	@Tags		admin
//	@Success	202	{object}	map[string]interface{}
//	@Failure	409	{object}	map[string]interface{}
//	@Router		/adm/metadata-storage/migrate [post]
func (c *Controller) AdmMetadataStorageMigrate(g *gin.Context) {
	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	var pending int64
	c.datastore.Model(&model.Task{}).
		Where("name = ? AND status IN ?", metamigrate.TaskName, []model.TaskStatus{model.TaskStatusTodo, model.TaskStatusInProgress}).
		Count(&pending)
	if pending > 0 {
		g.JSON(http.StatusConflict, gin.H{"error": "metadata migration is already queued or running"})
		return
	}
	if err := c.runner.NewTask(metamigrate.TaskName, map[string]string{}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusAccepted, gin.H{
		"message":  "metadata migration task queued",
		"task":     metamigrate.TaskName,
		"redirect": "/adm/tasks",
	})
}
