package controller

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
)

// --- mock provider implementation ---

type mockProvider struct {
	slug                string
	name                string
	searchActor         bool
	scrapePicture       bool
	actorSearchURL      string
	actorSearchResult   error
	actorGetThumbBytes  []byte
	actorGetThumbResult error
}

func (m *mockProvider) SlugGet() string               { return m.slug }
func (m *mockProvider) NiceName() string              { return m.name }
func (m *mockProvider) CapabilitySearchActor() bool   { return m.searchActor }
func (m *mockProvider) CapabilityScrapePicture() bool { return m.scrapePicture }
func (m *mockProvider) ActorSearch(bool, string) (string, error) {
	return m.actorSearchURL, m.actorSearchResult
}

func (m *mockProvider) ActorGetThumb(bool, string, string) ([]byte, error) {
	return m.actorGetThumbBytes, m.actorGetThumbResult
}

// --- setup helper ---

func setupTestController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&model.Provider{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	_logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&_logger)
	ctrl.DatabaseRegister(db)
	ctrl.providers = make(map[string]provider.Provider)

	return ctrl
}

// --- tests ---

func TestController_ProviderRegister_NewProvider(t *testing.T) {
	ctrl := setupTestController(t)

	mockProv := &mockProvider{
		slug:          "mockprov",
		name:          "Mock Provider",
		searchActor:   true,
		scrapePicture: false,
	}

	err := ctrl.ProviderRegister(mockProv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check map registration
	if _, ok := ctrl.providers["mockprov"]; !ok {
		t.Error("provider not registered in controller map")
	}

	// Check DB record
	var rec model.Provider
	if err := ctrl.datastore.First(&rec, "id = ?", "mockprov").Error; err != nil {
		t.Fatalf("expected provider in DB, got error: %v", err)
	}
	if rec.NiceName != "Mock Provider" {
		t.Errorf("expected NiceName 'Mock Provider', got %q", rec.NiceName)
	}
	if !rec.AbleToSearchActor {
		t.Error("expected AbleToSearchActor true")
	}
	if rec.AbleToScrapePicture {
		t.Error("expected AbleToScrapePicture false")
	}
}

func TestController_ProviderRegister_UpdateExisting(t *testing.T) {
	ctrl := setupTestController(t)

	// Create an existing record
	existing := model.Provider{
		ID:                  "mockprov",
		NiceName:            "Old Name",
		AbleToSearchActor:   false,
		AbleToScrapePicture: false,
	}
	if err := ctrl.datastore.Create(&existing).Error; err != nil {
		t.Fatalf("failed to create existing record: %v", err)
	}

	mockProv := &mockProvider{
		slug:          "mockprov",
		name:          "Updated Provider",
		searchActor:   true,
		scrapePicture: true,
	}

	err := ctrl.ProviderRegister(mockProv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var rec model.Provider
	if err := ctrl.datastore.First(&rec, "id = ?", "mockprov").Error; err != nil {
		t.Fatalf("expected provider in DB, got error: %v", err)
	}
	if rec.NiceName != "Updated Provider" {
		t.Errorf("expected NiceName 'Updated Provider', got %q", rec.NiceName)
	}
	if !rec.AbleToSearchActor || !rec.AbleToScrapePicture {
		t.Error("expected updated capability flags to be true")
	}
}

func TestController_ProviderGet_Success(t *testing.T) {
	ctrl := setupTestController(t)

	mockProv := &mockProvider{slug: "p1"}
	ctrl.providers["p1"] = mockProv

	got, err := ctrl.ProviderGet("p1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.SlugGet() != "p1" {
		t.Errorf("expected slug 'p1', got %s", got.SlugGet())
	}
}

func TestController_ProviderGet_NotFound(t *testing.T) {
	ctrl := setupTestController(t)

	_, err := ctrl.ProviderGet("doesnotexist")
	if err == nil {
		t.Fatal("expected error for missing provider, got nil")
	}
	if !errors.Is(err, err) && err.Error() != "provider not found" {
		t.Errorf("unexpected error: %v", err)
	}
}
