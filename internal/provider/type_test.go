package provider

import (
	"testing"
)

// --- Tests for constants and interface ---

func TestErrOfflineMode(t *testing.T) {
	if ErrOfflineMode == nil {
		t.Fatal("expected ErrOfflineMode to be defined")
	}
	if ErrOfflineMode.Error() != "not possible in offline mode" {
		t.Errorf("unexpected error text: %s", ErrOfflineMode.Error())
	}
}
