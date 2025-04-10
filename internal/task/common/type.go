package common

import (
	"log"
	"time"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"gorm.io/gorm"
)

type Parameters map[string]string

type Context struct {
	DB     *gorm.DB
	Config *config.Config
}

type Step struct {
	Name     string
	NiceName string
	Func     func(*Context, Parameters) (string, error)
}

type Task struct {
	Name  string
	Steps []Step
}

func (t *Task) Run(ctx *Context) {
	// select all task matching object + action
	log.Println("start runner for task: " + t.Name)
	var tasks []model.Task
	ctx.DB.Where("name = ? and status = ?", t.Name, "todo").Find(&tasks)

	for _, task := range tasks {
		t.runTask(ctx, &task)
	}
}

func (t *Task) runTask(ctx *Context, task *model.Task) {
	// set first step if first run
	if task.Step == "" {
		task.Step = t.Steps[0].Name
		ctx.DB.Save(task)
	}

	log.SetPrefix(t.Name + "/" + task.ID + " ")

	log.Println("running on step: ", task.Step)

	stepNB := findStepNB(t.Steps, task.Step)

	errMsg, err := t.Steps[stepNB].Func(ctx, task.Parameters)

	if err != nil {
		log.Println("task failed:", errMsg, err.Error())
		task.Status = model.TaskStatusError
		ctx.DB.Save(&task)
		return
	}

	stepNB++
	if stepNB == len(t.Steps) {
		task.Status = model.TaskStatusDone
		now := time.Now()
		task.DoneAt = &now
		ctx.DB.Save(&task)
		log.Println("task done")
		return
	}

	task.Step = t.Steps[stepNB].Name
	ctx.DB.Save(&task)
	t.runTask(ctx, task)
}

func findStepNB(steps []Step, step string) int {
	for stepNB := range steps {
		if steps[stepNB].Name == step {
			return stepNB
		}
	}

	return -1
}
