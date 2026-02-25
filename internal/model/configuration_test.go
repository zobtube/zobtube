package model

import (
	"testing"
)

func TestConfiguration_ZeroValue(t *testing.T) {
	var c Configuration
	if c.OfflineMode != false {
		t.Errorf("zero value OfflineMode = %v, want false", c.OfflineMode)
	}
}

func TestConfiguration_FieldAssignment(t *testing.T) {
	c := Configuration{
		ID:                 1,
		UserAuthentication: true,
		OfflineMode:        true,
	}
	if c.ID != 1 || !c.UserAuthentication || !c.OfflineMode {
		t.Errorf("Configuration fields not set correctly: %+v", c)
	}
}
