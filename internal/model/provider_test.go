package model

import (
	"testing"
)

func TestProvider_ZeroValue(t *testing.T) {
	var p Provider
	if p.Enabled != false {
		t.Errorf("zero value Enabled = %v, want false", p.Enabled)
	}
}

func TestProvider_FieldAssignment(t *testing.T) {
	p := Provider{
		ID:                "test-provider",
		NiceName:          "Test",
		Enabled:           true,
		AbleToSearchActor: true,
	}
	if p.ID != "test-provider" || p.NiceName != "Test" || !p.Enabled || !p.AbleToSearchActor {
		t.Errorf("Provider fields not set correctly: %+v", p)
	}
}
