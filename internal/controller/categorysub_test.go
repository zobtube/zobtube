package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func TestController_CategorySubAdd_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupCategoryController(t)

	parent := &model.Category{Name: "Parent"}
	if err := ctrl.datastore.Create(parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/category-sub", strings.NewReader("Name=Sub+Item&Parent="+parent.ID))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctrl.CategorySubAdd(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	id, _ := body["id"].(string)
	if id == "" {
		t.Fatalf("expected id in response, got %v", body)
	}

	var sub model.CategorySub
	if ctrl.datastore.First(&sub, "id = ?", id).RowsAffected < 1 {
		t.Fatal("sub-category not created in DB")
	}
	if sub.Name != "Sub Item" || sub.Category != parent.ID {
		t.Errorf("unexpected sub: name=%q category=%q", sub.Name, sub.Category)
	}
}
