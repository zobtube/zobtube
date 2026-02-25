package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCategorySubDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&CategorySub{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestCategorySub_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupCategorySubDB(t)
	c := &CategorySub{Name: "sub", Category: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(c.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", c.ID, err)
	}
}

func TestCategorySub_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupCategorySubDB(t)
	c := &CategorySub{ID: "00000000-0000-0000-0000-000000000000", Name: "sub", Category: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.ID == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected zero UUID to be replaced")
	}
	if _, err := uuid.Parse(c.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", c.ID, err)
	}
}

func TestCategorySub_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupCategorySubDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	c := &CategorySub{ID: wantID, Name: "sub", Category: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, c.ID)
	}
}

func TestCategorySub_URLThumb(t *testing.T) {
	id := "deadbeef-1234-5678-abcd-000000000000"
	c := &CategorySub{ID: id}
	if got := c.URLThumb(); got != "/api/category-sub/"+id+"/thumb" {
		t.Errorf("URLThumb() = %q", got)
	}
}
