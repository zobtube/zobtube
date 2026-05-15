package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
	"github.com/zobtube/zobtube/internal/task/metamigrate"
)

func TestController_AdmMetadataStorage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/metadata-storage", nil)

	ctrl.AdmMetadataStorage(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["type"] != "filesystem" {
		t.Fatalf("expected filesystem type, got %v", body["type"])
	}
	if body["path"] != "/tmp" {
		t.Fatalf("expected path /tmp, got %v", body["path"])
	}
	if _, ok := body["secret_access_key"]; ok {
		t.Fatal("secret access key must not be exposed")
	}
}

func TestController_AdmMetadataStorageMigrate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.runner.RegisterTask(&common.Task{
		Name: metamigrate.TaskName,
		Steps: []common.Step{{
			Name: "noop",
			Func: func(*common.Context, common.Parameters) (string, error) { return "", nil },
		}},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/metadata-storage/migrate", nil)
	ctrl.AdmMetadataStorageMigrate(c)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["task"] != metamigrate.TaskName {
		t.Fatalf("expected task %q, got %v", metamigrate.TaskName, body["task"])
	}
	var count int64
	ctrl.datastore.Model(&model.Task{}).Where("name = ?", metamigrate.TaskName).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 queued task, got %d", count)
	}
}

func TestController_AdmMetadataStorageMigrate_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	if err := ctrl.datastore.Create(&model.Task{
		Name:   metamigrate.TaskName,
		Status: model.TaskStatusInProgress,
	}).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/metadata-storage/migrate", nil)
	ctrl.AdmMetadataStorageMigrate(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
}
