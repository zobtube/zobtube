package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func TestController_AdmOrganizationList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/organizations", nil)
	ctrl.AdmOrganizationList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	items, _ := body["items"].([]any)
	if len(items) != 0 {
		t.Errorf("expected empty items, got %d", len(items))
	}
}

func TestController_AdmOrganizationCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/organizations",
		strings.NewReader(`{"name":"v2","template":"$TYPE/$ID/v.mp4","active":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmOrganizationCreate(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var count int64
	ctrl.datastore.Model(&model.Organization{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 organization, got %d", count)
	}
}

func TestController_AdmOrganizationCreate_RejectsTemplateWithoutID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/organizations",
		strings.NewReader(`{"name":"bad","template":"videos/video.mp4"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmOrganizationCreate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestController_AdmOrganizationCreate_ActiveDeactivatesOthers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Organization{ID: "11111111-1111-1111-1111-111111111111", Name: "old", Template: "$TYPE/$ID/v.mp4", Active: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/organizations",
		strings.NewReader(`{"name":"v2","template":"alt/$ID/v.mp4","active":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmOrganizationCreate(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var activeCount int64
	ctrl.datastore.Model(&model.Organization{}).Where("active = ?", true).Count(&activeCount)
	if activeCount != 1 {
		t.Errorf("expected exactly 1 active organization, got %d", activeCount)
	}
}

func TestController_AdmOrganizationActivate_Switches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Organization{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "A", Template: "$TYPE/$ID/v.mp4", Active: true})
	ctrl.datastore.Create(&model.Organization{ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", Name: "B", Template: "alt/$ID/v.mp4", Active: false})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"}}
	c.Request = httptest.NewRequest("POST", "/api/adm/organizations/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb/activate", nil)
	ctrl.AdmOrganizationActivate(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var orgA, orgB model.Organization
	ctrl.datastore.First(&orgA, "id = ?", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	ctrl.datastore.First(&orgB, "id = ?", "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	if orgA.Active {
		t.Error("organization A should have been deactivated")
	}
	if !orgB.Active {
		t.Error("organization B should now be active")
	}
}

func TestController_AdmOrganizationDelete_RejectsActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Organization{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "A", Template: "$TYPE/$ID/v.mp4", Active: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}}
	c.Request = httptest.NewRequest("DELETE", "/api/adm/organizations/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", nil)
	ctrl.AdmOrganizationDelete(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for active organization, got %d", w.Code)
	}
}

func TestController_AdmOrganizationDelete_RejectsWhenVideosLinked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	orgID := "11111111-1111-1111-1111-111111111111"
	ctrl.datastore.Create(&model.Organization{ID: orgID, Name: "A", Template: "$TYPE/$ID/v.mp4", Active: false})
	ctrl.datastore.Create(&model.Video{Name: "V", Filename: "v.mp4", Type: "v", OrganizationID: &orgID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: orgID}}
	c.Request = httptest.NewRequest("DELETE", "/api/adm/organizations/"+orgID, nil)
	ctrl.AdmOrganizationDelete(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when videos are linked, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestController_AdmOrganizationUpdate_LocksTemplateWhenUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	orgID := "11111111-1111-1111-1111-111111111111"
	ctrl.datastore.Create(&model.Organization{ID: orgID, Name: "A", Template: "$TYPE/$ID/v.mp4", Active: true})
	ctrl.datastore.Create(&model.Video{Name: "V", Filename: "v.mp4", Type: "v", OrganizationID: &orgID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: orgID}}
	c.Request = httptest.NewRequest("PUT", "/api/adm/organizations/"+orgID,
		strings.NewReader(`{"template":"other/$ID/v.mp4"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmOrganizationUpdate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestController_AdmOrganizationUpdate_AllowsRenameAlways(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	orgID := "11111111-1111-1111-1111-111111111111"
	ctrl.datastore.Create(&model.Organization{ID: orgID, Name: "A", Template: "$TYPE/$ID/v.mp4", Active: true})
	ctrl.datastore.Create(&model.Video{Name: "V", Filename: "v.mp4", Type: "v", OrganizationID: &orgID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: orgID}}
	c.Request = httptest.NewRequest("PUT", "/api/adm/organizations/"+orgID,
		strings.NewReader(`{"name":"Renamed"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.AdmOrganizationUpdate(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var got model.Organization
	ctrl.datastore.First(&got, "id = ?", orgID)
	if got.Name != "Renamed" {
		t.Errorf("expected name Renamed, got %q", got.Name)
	}
}

func TestController_AdmOrganizationReorganize_QueuesTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.runner.RegisterTask(&common.Task{
		Name:  "video/reorganize",
		Steps: []common.Step{{Name: "noop", Func: func(*common.Context, common.Parameters) (string, error) { return "", nil }}},
	})

	srcOrg := "11111111-1111-1111-1111-111111111111"
	dstOrg := "22222222-2222-2222-2222-222222222222"
	ctrl.datastore.Create(&model.Organization{ID: srcOrg, Name: "src", Template: "$TYPE/$ID/v.mp4", Active: true})
	ctrl.datastore.Create(&model.Organization{ID: dstOrg, Name: "dst", Template: "alt/$ID/v.mp4", Active: false})

	v := &model.Video{Name: "V", Filename: "v.mp4", Type: "v", Imported: true, OrganizationID: &srcOrg}
	ctrl.datastore.Create(v)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: dstOrg}}
	c.Request = httptest.NewRequest("POST", "/api/adm/organizations/"+dstOrg+"/reorganize", nil)
	ctrl.AdmOrganizationReorganize(c)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d body=%s", w.Code, w.Body.String())
	}
	var count int64
	ctrl.datastore.Model(&model.Task{}).Where("name = ?", "video/reorganize").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 reorganize task queued, got %d", count)
	}
}

func TestController_AdmConfigReorganizeOnImport_Toggle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Configuration{ID: 1, ReorganizeOnImport: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "action", Value: "disable"}}
	c.Request = httptest.NewRequest("GET", "/api/adm/config/reorganize-on-import/disable", nil)
	ctrl.AdmConfigReorganizeOnImportUpdate(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var cfg model.Configuration
	ctrl.datastore.First(&cfg)
	if cfg.ReorganizeOnImport {
		t.Error("expected ReorganizeOnImport to be disabled")
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Params = gin.Params{{Key: "action", Value: "enable"}}
	c2.Request = httptest.NewRequest("GET", "/api/adm/config/reorganize-on-import/enable", nil)
	ctrl.AdmConfigReorganizeOnImportUpdate(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 after enable, got %d", w2.Code)
	}
	ctrl.datastore.First(&cfg)
	if !cfg.ReorganizeOnImport {
		t.Error("expected ReorganizeOnImport to be enabled again")
	}
}

func TestController_AdmConfigReorganizeOnImport_InvalidAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "action", Value: "bogus"}}
	c.Request = httptest.NewRequest("GET", "/api/adm/config/reorganize-on-import/bogus", nil)
	ctrl.AdmConfigReorganizeOnImportUpdate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
