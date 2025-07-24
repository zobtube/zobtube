package runner

import (
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

type RunnerEvent int

const (
	NewTaskEvent RunnerEvent = 1
)

type Runner struct {
	tasks        []*common.Task
	tasksChannel map[string]chan RunnerEvent
	ctx          *common.Context
}

func (r *Runner) RegisterTask(t *common.Task) {
	r.tasks = append(r.tasks, t)
	if r.tasksChannel == nil {
		r.tasksChannel = make(map[string]chan RunnerEvent)
	}
	r.tasksChannel[t.Name] = make(chan RunnerEvent)
}

func (r *Runner) Start(cfg *config.Config, db *gorm.DB) {
	r.ctx = &common.Context{
		DB:     db,
		Config: cfg,
	}
	for _, task := range r.tasks {
		go func() {
			for {
				event := <-r.tasksChannel[task.Name]
				if event == NewTaskEvent {
					task.Run(r.ctx)
				}
			}
		}()
	}
}

func (r *Runner) NewTask(action string, params map[string]string) error {
	task := &model.Task{
		Name:       action,
		Parameters: params,
	}

	err := r.ctx.DB.Create(&task).Error
	if err != nil {
		return err
	}

	r.tasksChannel[action] <- NewTaskEvent

	return nil
}

func (r *Runner) TaskRetry(action string) {
	r.tasksChannel[action] <- NewTaskEvent
}
