package controller

import (
	"time"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CleanupRoutine() {
	go c.taskRestart()
	for {
		// wait loop
		time.Sleep(time.Minute)

		c.sessionCleanup()
	}
}

func (c *Controller) sessionCleanup() {
	var sessions []model.UserSession
	result := c.datastore.Find(&sessions)
	err := result.Error
	if err != nil {
		c.logger.Warn().Err(err).Msg("error while querying sessions")
		return
	}

	for _, session := range sessions {
		if session.ValidUntil.Before(time.Now()) {
			c.datastore.Delete(&session)
		}
	}
}

func (c *Controller) taskRestart() {
	var tasks []model.Task
	result := c.datastore.Where("status = ?", model.TaskStatusTodo).Find(&tasks)
	err := result.Error
	if err != nil {
		c.logger.Warn().Err(err).Msg("error while querying tasks")
		return
	}

	for _, task := range tasks {
		c.logger.Warn().Str("kind", "tasks").Str("task-id", task.ID).Msg("stuck in todo, restarting")
		c.runner.TaskRetry(task.Name)
	}
}
