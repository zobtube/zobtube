package model

import (
	"testing"
)

func TestVideoView_ZeroValue(t *testing.T) {
	var v VideoView
	if v.Count != 0 {
		t.Errorf("zero value Count = %d, want 0", v.Count)
	}
}

func TestVideoView_FieldAssignment(t *testing.T) {
	v := VideoView{
		VideoID: "vid-123",
		UserID:  "user-456",
		Count:   3,
	}
	if v.VideoID != "vid-123" || v.UserID != "user-456" || v.Count != 3 {
		t.Errorf("VideoView fields not set correctly: %+v", v)
	}
}
