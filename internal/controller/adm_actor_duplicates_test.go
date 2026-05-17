package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func setupAdmActorDuplicatesController(t *testing.T) *Controller {
	t.Helper()
	ctrl := setupAdmController(t)
	if err := ctrl.datastore.AutoMigrate(&model.ActorDismissedDuplicate{}); err != nil {
		t.Fatalf("migrate ActorDismissedDuplicate: %v", err)
	}
	return ctrl
}

func TestAdmActorDuplicates_FindsDuplicates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmActorDuplicatesController(t)

	a1 := &model.Actor{Name: "Alice", Sex: "f"}
	a2 := &model.Actor{Name: "alice", Sex: "f"}
	a3 := &model.Actor{Name: "Bob", Sex: "m"}
	ctrl.datastore.Create(a1)
	ctrl.datastore.Create(a2)
	ctrl.datastore.Create(a3)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/adm/actor/duplicates", nil)
	ctrl.AdmActorDuplicates(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["total"].(float64) != 1 {
		t.Fatalf("expected 1 group, got %v", body["total"])
	}
	groups, ok := body["groups"].([]any)
	if !ok || len(groups) != 1 {
		t.Fatalf("expected 1 group in response, got %#v", body["groups"])
	}
}

func TestAdmActorDuplicates_ExcludesDismissed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmActorDuplicatesController(t)

	a1 := &model.Actor{Name: "Alice", Sex: "f"}
	a2 := &model.Actor{Name: "alice", Sex: "f"}
	ctrl.datastore.Create(a1)
	ctrl.datastore.Create(a2)

	id1, id2 := model.NormalizeActorPair(a1.ID, a2.ID)
	ctrl.datastore.Create(&model.ActorDismissedDuplicate{ActorID1: id1, ActorID2: id2})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/adm/actor/duplicates", nil)
	ctrl.AdmActorDuplicates(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["total"].(float64) != 0 {
		t.Fatalf("expected 0 groups after dismiss, got %v", body["total"])
	}
}

func TestAdmActorDuplicateDismiss(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmActorDuplicatesController(t)

	a1 := &model.Actor{Name: "Alice", Sex: "f"}
	a2 := &model.Actor{Name: "alice", Sex: "f"}
	ctrl.datastore.Create(a1)
	ctrl.datastore.Create(a2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"actor_id_1":"` + a1.ID + `","actor_id_2":"` + a2.ID + `"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/api/adm/actor/duplicates/dismiss", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmActorDuplicateDismiss(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var count int64
	ctrl.datastore.Model(&model.ActorDismissedDuplicate{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 dismissed record, got %d", count)
	}
}

func TestAdmActorDuplicateDismissRemove(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmActorDuplicatesController(t)

	a1 := &model.Actor{Name: "Alice", Sex: "f"}
	a2 := &model.Actor{Name: "alice", Sex: "f"}
	ctrl.datastore.Create(a1)
	ctrl.datastore.Create(a2)

	id1, id2 := model.NormalizeActorPair(a1.ID, a2.ID)
	record := &model.ActorDismissedDuplicate{ActorID1: id1, ActorID2: id2}
	ctrl.datastore.Create(record)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: record.ID}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/adm/actor/duplicates/dismiss/"+record.ID, nil)
	ctrl.AdmActorDuplicateDismissRemove(c)

	if w.Code != http.StatusNoContent {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 204, got %d: %s", w.Code, string(body))
	}

	var count int64
	ctrl.datastore.Model(&model.ActorDismissedDuplicate{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 dismissed records, got %d", count)
	}
}
