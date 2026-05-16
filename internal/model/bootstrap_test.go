package model

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupBootstrapDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Library{}, &Organization{}, &Video{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestEnsureDefaultOrganization_CreatesWhenEmpty(t *testing.T) {
	db := setupBootstrapDB(t)
	id, err := EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("EnsureDefaultOrganization: %v", err)
	}
	if id != DefaultOrganizationUUID {
		t.Errorf("expected default UUID, got %s", id)
	}
	var org Organization
	if err := db.First(&org, "id = ?", id).Error; err != nil {
		t.Fatalf("organization not persisted: %v", err)
	}
	if !org.Active {
		t.Error("default organization must be active")
	}
	if org.Template != DefaultOrganizationTemplate {
		t.Errorf("template = %q, want %q", org.Template, DefaultOrganizationTemplate)
	}
}

func TestEnsureDefaultOrganization_ReturnsActiveWhenExists(t *testing.T) {
	db := setupBootstrapDB(t)
	custom := Organization{ID: "11111111-1111-1111-1111-111111111111", Name: "Custom", Template: "$TYPE/$ID/v.mp4", Active: true}
	if err := db.Create(&custom).Error; err != nil {
		t.Fatalf("create custom: %v", err)
	}
	id, err := EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("EnsureDefaultOrganization: %v", err)
	}
	if id != custom.ID {
		t.Errorf("expected existing active id %s, got %s", custom.ID, id)
	}
	var count int64
	db.Model(&Organization{}).Count(&count)
	if count != 1 {
		t.Errorf("must not create a duplicate, got %d rows", count)
	}
}

func TestBackfillVideoOrganization_AssignsPathAndOrgToImported(t *testing.T) {
	db := setupBootstrapDB(t)
	orgID, err := EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("ensure org: %v", err)
	}
	imported := Video{ID: "11111111-1111-1111-1111-111111111111", Type: "v", Name: "n", Filename: "f.mp4", Imported: true}
	if err := db.Create(&imported).Error; err != nil {
		t.Fatalf("create imported: %v", err)
	}
	staging := Video{ID: "22222222-2222-2222-2222-222222222222", Type: "v", Name: "n2", Filename: "f2.mp4", Imported: false}
	if err := db.Create(&staging).Error; err != nil {
		t.Fatalf("create staging: %v", err)
	}

	if err := BackfillVideoOrganization(db, orgID); err != nil {
		t.Fatalf("BackfillVideoOrganization: %v", err)
	}

	var got Video
	if err := db.First(&got, "id = ?", imported.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.OrganizationID == nil || *got.OrganizationID != orgID {
		t.Errorf("expected organization_id = %q, got %v", orgID, got.OrganizationID)
	}
	wantPath := filepath.Join("/videos", imported.ID, "video.mp4")
	if got.Path == nil || *got.Path != wantPath {
		t.Errorf("expected path %q, got %v", wantPath, got.Path)
	}

	var stagingAfter Video
	if err := db.First(&stagingAfter, "id = ?", staging.ID).Error; err != nil {
		t.Fatal(err)
	}
	if stagingAfter.OrganizationID != nil {
		t.Errorf("non-imported videos must not be backfilled, got %v", stagingAfter.OrganizationID)
	}
}

func TestBackfillVideoOrganization_SkipsTriageInPlacePath(t *testing.T) {
	db := setupBootstrapDB(t)
	orgID, err := EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("ensure org: %v", err)
	}
	triagePath := "triage/keep.mp4"
	v := Video{ID: "55555555-5555-5555-5555-555555555555", Type: "v", Name: "n", Filename: "keep.mp4", Imported: true, Path: &triagePath}
	if err := db.Create(&v).Error; err != nil {
		t.Fatalf("create video: %v", err)
	}
	if err := BackfillVideoOrganization(db, orgID); err != nil {
		t.Fatalf("backfill: %v", err)
	}
	var got Video
	if err := db.First(&got, "id = ?", v.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.OrganizationID != nil {
		t.Errorf("in-place triage video must keep organization_id nil, got %v", got.OrganizationID)
	}
	if got.Path == nil || *got.Path != triagePath {
		t.Errorf("path = %v, want %q", got.Path, triagePath)
	}
}

func TestBackfillVideoOrganization_DoesNotOverwriteExisting(t *testing.T) {
	db := setupBootstrapDB(t)
	orgID, err := EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("ensure org: %v", err)
	}
	existingOrgID := "33333333-3333-3333-3333-333333333333"
	existingPath := "custom/x.mp4"
	v := Video{ID: "44444444-4444-4444-4444-444444444444", Type: "v", Name: "n", Filename: "f.mp4", Imported: true, OrganizationID: &existingOrgID, Path: &existingPath}
	if err := db.Create(&v).Error; err != nil {
		t.Fatalf("create video: %v", err)
	}
	if err := BackfillVideoOrganization(db, orgID); err != nil {
		t.Fatalf("backfill: %v", err)
	}
	var got Video
	if err := db.First(&got, "id = ?", v.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.OrganizationID == nil || *got.OrganizationID != existingOrgID {
		t.Errorf("must not overwrite existing organization_id, got %v", got.OrganizationID)
	}
	if got.Path == nil || *got.Path != existingPath {
		t.Errorf("must not overwrite existing path, got %v", got.Path)
	}
}
