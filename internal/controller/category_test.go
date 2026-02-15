package controller

import (
	"encoding/json"
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
)

func setupCategoryController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Category{}, &model.CategorySub{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{})

	return ctrl
}

func TestController_CategoryList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/category", nil)

	ctrl.CategoryList(c)

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

func TestController_CategoryAdd_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/category/add", strings.NewReader(`Name=Comedy`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:errcheck
	c.Request.PostForm.Set("Name", "Comedy")

	ctrl.CategoryAdd(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var cat model.Category
	if ctrl.datastore.First(&cat, "name = ?", "Comedy").RowsAffected < 1 {
		t.Error("category not created in DB")
	}
}

func TestController_CategoryAdd_EmptyName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/category/add", strings.NewReader(`name=`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctrl.CategoryAdd(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "category name cannot be empty" {
		t.Errorf("expected empty name error, got %v", errMsg)
	}
}

func TestController_CategorySubGet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)
	cat := &model.Category{Name: "Parent"}
	ctrl.datastore.Create(cat)
	sub := &model.CategorySub{Name: "Sub1", Category: cat.ID}
	ctrl.datastore.Create(sub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: sub.ID}}
	c.Request = httptest.NewRequest("GET", "/api/category/sub/"+sub.ID, nil)

	ctrl.CategorySubGet(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_CategorySubGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/category/sub/nonexistent", nil)

	ctrl.CategorySubGet(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_CategoryDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)
	cat := &model.Category{Name: "ToDelete"}
	ctrl.datastore.Create(cat)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: cat.ID}}
	c.Request = httptest.NewRequest("DELETE", "/api/category/"+cat.ID, nil)

	ctrl.CategoryDelete(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ctrl.datastore.First(&model.Category{}, "id = ?", cat.ID).RowsAffected > 0 {
		t.Error("category should be deleted")
	}
}

func TestController_CategoryDelete_WithSubCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)
	cat := &model.Category{Name: "Parent"}
	ctrl.datastore.Create(cat)
	sub := &model.CategorySub{Name: "Sub", Category: cat.ID}
	ctrl.datastore.Create(sub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: cat.ID}}
	c.Request = httptest.NewRequest("DELETE", "/api/category/"+cat.ID, nil)

	ctrl.CategoryDelete(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "category cannot be deleted with values presents" {
		t.Errorf("expected sub-category error, got %v", errMsg)
	}
}
