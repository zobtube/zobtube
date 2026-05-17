package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupActorDismissedDuplicateDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&ActorDismissedDuplicate{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestNormalizeActorPair(t *testing.T) {
	a, b := NormalizeActorPair("b-id", "a-id")
	if a != "a-id" || b != "b-id" {
		t.Fatalf("expected (a-id, b-id), got (%s, %s)", a, b)
	}
}

func TestActorDismissedDuplicate_BeforeCreateAssignsID(t *testing.T) {
	db := setupActorDismissedDuplicateDB(t)
	record := &ActorDismissedDuplicate{
		ActorID1: "00000000-0000-0000-0000-000000000001",
		ActorID2: "00000000-0000-0000-0000-000000000002",
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if record.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
}
