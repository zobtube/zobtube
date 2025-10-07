package controller

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/task/common"
)

// --- setup helpers ---

func setupCleanupTestController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&model.UserSession{}, &model.Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)

	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.runner = &runner.Runner{} // default runner (can mock below)

	return ctrl
}

// --- tests ---

func TestController_SessionCleanup_RemovesExpiredSessions(t *testing.T) {
	ctrl := setupCleanupTestController(t)

	// Insert expired and valid sessions
	expired := model.UserSession{ValidUntil: time.Now().Add(-1 * time.Hour)}
	valid := model.UserSession{ValidUntil: time.Now().Add(1 * time.Hour)}
	if err := ctrl.datastore.Create(&expired).Error; err != nil {
		t.Fatalf("failed to insert expired session: %v", err)
	}
	if err := ctrl.datastore.Create(&valid).Error; err != nil {
		t.Fatalf("failed to insert valid session: %v", err)
	}

	ctrl.sessionCleanup()

	var sessions []model.UserSession
	if err := ctrl.datastore.Find(&sessions).Error; err != nil {
		t.Fatalf("db query failed: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 remaining session, got %d", len(sessions))
	}
	if sessions[0].ID != valid.ID {
		t.Errorf("expected valid session to remain, got ID=%s", sessions[0].ID)
	}
}

func TestController_SessionCleanup_DBError(t *testing.T) {
	ctrl := setupCleanupTestController(t)

	// Force datastore error by closing DB (simulate query failure)
	sqlDB, _ := ctrl.datastore.DB()
	sqlDB.Close()

	// Should not panic
	ctrl.sessionCleanup()
}

// --- Task restart tests ---

func TestController_TaskRestart_RetriesTodoTasks(t *testing.T) {
	ctrl := setupCleanupTestController(t)

	// Add one task in TODO state
	task1 := model.Task{Name: "T1", Status: model.TaskStatusTodo}
	if err := ctrl.datastore.Create(&task1).Error; err != nil {
		t.Fatalf("failed to insert T1: %v", err)
	}

	// Channel to detect that our step ran
	done := make(chan struct{}, 1)

	// Define a test step that signals completion
	step := common.Step{
		Name:     "step1",
		NiceName: "Step 1",
		Func: func(ctx *common.Context, p common.Parameters) (string, error) {
			done <- struct{}{}
			return "ok", nil
		},
	}

	// Register a Runner with that task
	r := &runner.Runner{}
	r.RegisterTask(&common.Task{
		Name:  "T1",
		Steps: []common.Step{step},
	})
	r.Start(&config.Config{}, ctrl.datastore)
	ctrl.RunnerRegister(r)

	// Call taskRestart (should trigger TaskRetry)
	ctrl.taskRestart()

	select {
	case <-done:
		// Success: our step ran
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for step execution after taskRestart")
	}
}

func TestController_TaskRestart_DBError(t *testing.T) {
	ctrl := setupCleanupTestController(t)
	sqlDB, _ := ctrl.datastore.DB()
	_ = sqlDB.Close()
	ctrl.taskRestart() // should not panic
}
