package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/task/common"
)

func setupVideoController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Video{}, &model.Actor{}, &model.Channel{}, &model.Category{},
		&model.CategorySub{}, &model.VideoView{}, &model.Task{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{})

	// Runner with no-op video/create task for VideoCreate handler
	r := &runner.Runner{}
	r.RegisterTask(&common.Task{
		Name: "video/create",
		Steps: []common.Step{{
			Name:     "noop",
			NiceName: "No-op",
			Func:     func(*common.Context, common.Parameters) (string, error) { return "", nil },
		}},
	})
	r.Start(&config.Config{}, db)
	ctrl.RunnerRegister(r)

	return ctrl
}

func TestController_VideoList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/video", nil)

	ctrl.VideoList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 0 {
		t.Errorf("expected total 0, got %v", total)
	}
}

func TestController_VideoList_WithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	ctrl.datastore.Create(&model.Video{Name: "V1", Filename: "v1.mp4", Type: "v"})
	ctrl.datastore.Create(&model.Video{Name: "V2", Filename: "v2.mp4", Type: "v"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/video", nil)

	ctrl.VideoList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 2 {
		t.Errorf("expected total 2, got %v", total)
	}
}

func TestController_ClipList_WithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	ctrl.datastore.Create(&model.Video{Name: "C1", Filename: "c1.mp4", Type: "c"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/clip", nil)

	ctrl.ClipList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 1 {
		t.Errorf("expected total 1, got %v", total)
	}
}

func TestController_MovieList_WithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	ctrl.datastore.Create(&model.Video{Name: "M1", Filename: "m1.mp4", Type: "m"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/movie", nil)

	ctrl.MovieList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 1 {
		t.Errorf("expected total 1, got %v", total)
	}
}

func TestController_VideoView_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "Test", Filename: "test.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("GET", "/api/video/"+vid.ID, nil)

	ctrl.VideoView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["video"] == nil {
		t.Error("expected video in response")
	}
}

func TestController_VideoView_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/video/nonexistent", nil)

	ctrl.VideoView(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_VideoEdit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "Test", Filename: "test.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("GET", "/api/video/"+vid.ID+"/edit", nil)

	ctrl.VideoEdit(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["video"] == nil {
		t.Error("expected video in response")
	}
}

func TestController_VideoEdit_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/video/nonexistent/edit", nil)

	ctrl.VideoEdit(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_VideoActors_Put_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	actor := &model.Actor{Name: "A"}
	ctrl.datastore.Create(vid)
	ctrl.datastore.Create(actor)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}, {Key: "actor_id", Value: actor.ID}}
	c.Request = httptest.NewRequest("PUT", "/api/video/"+vid.ID+"/actor/"+actor.ID, nil)

	ctrl.VideoActors(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}
}

func TestController_VideoActors_VideoNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	actor := &model.Actor{Name: "A"}
	ctrl.datastore.Create(actor)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}, {Key: "actor_id", Value: actor.ID}}
	c.Request = httptest.NewRequest("PUT", "/api/video/nonexistent/actor/"+actor.ID, nil)

	ctrl.VideoActors(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_VideoCategories_Put_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	cat := &model.Category{Name: "Cat"}
	ctrl.datastore.Create(cat)
	sub := &model.CategorySub{Name: "Sub", Category: cat.ID}
	ctrl.datastore.Create(sub)
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}, {Key: "category_id", Value: sub.ID}}
	c.Request = httptest.NewRequest("PUT", "/api/video/"+vid.ID+"/category/"+sub.ID, nil)

	ctrl.VideoCategories(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}
}

func TestController_VideoCategories_SubCategoryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}, {Key: "category_id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("PUT", "/api/video/"+vid.ID+"/category/nonexistent", nil)

	ctrl.VideoCategories(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_VideoRename_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "Old", Filename: "old.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("POST", "/api/video/"+vid.ID+"/rename", strings.NewReader(`name=NewName`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctrl.VideoRename(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var updated model.Video
	ctrl.datastore.First(&updated, "id = ?", vid.ID)
	if updated.Name != "NewName" {
		t.Errorf("expected name 'NewName', got %q", updated.Name)
	}
}

func TestController_VideoRename_WrongMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("GET", "/api/video/"+vid.ID+"/rename", nil)

	ctrl.VideoRename(c)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestController_VideoCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/video", strings.NewReader(`name=MyVideo&filename=my.mp4&type=v`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctrl.VideoCreate(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["video_id"] == nil || body["video_id"] == "" {
		t.Error("expected video_id in response")
	}
}

func TestController_VideoCreate_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/video", strings.NewReader(`name=X&filename=x.mp4&type=x`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctrl.VideoCreate(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "invalid input" {
		t.Errorf("expected invalid input error, got %v", errMsg)
	}
}
