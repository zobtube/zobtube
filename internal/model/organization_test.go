package model

import (
	"errors"
	"testing"
)

func TestValidateOrganizationTemplate(t *testing.T) {
	cases := []struct {
		name    string
		tmpl    string
		wantErr error
	}{
		{"empty", "", errors.New("any")},
		{"whitespace only", "   ", errors.New("any")},
		{"no $ID", "videos/video.mp4", ErrOrganizationTemplateMissingID},
		{"valid default", DefaultOrganizationTemplate, nil},
		{"valid custom", "$TYPE/$BASENAME-$ID$EXT", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOrganizationTemplate(tc.tmpl)
			if tc.wantErr == nil && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if tc.wantErr != nil && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tc.wantErr != nil && errors.Is(tc.wantErr, ErrOrganizationTemplateMissingID) && !errors.Is(err, ErrOrganizationTemplateMissingID) {
				t.Fatalf("expected ErrOrganizationTemplateMissingID, got %v", err)
			}
		})
	}
}

func TestOrganization_Render_DefaultTemplateMatchesLegacy(t *testing.T) {
	org := &Organization{Template: DefaultOrganizationTemplate}
	cases := []struct {
		typeLetter string
		want       string
	}{
		{"v", "videos/abc/video.mp4"},
		{"c", "clips/abc/video.mp4"},
		{"m", "movies/abc/video.mp4"},
	}
	for _, tc := range cases {
		t.Run(tc.typeLetter, func(t *testing.T) {
			v := &Video{ID: "abc", Type: tc.typeLetter, Filename: "raw.mp4"}
			if got := org.Render(v); got != tc.want {
				t.Fatalf("Render() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOrganization_Render_AllVariables(t *testing.T) {
	org := &Organization{Template: "$TYPE_LETTER/$TYPE/$TYPE_NAME/$ID/$BASENAME$EXT"}
	v := &Video{ID: "deadbeef", Type: "v", Filename: "My File.mp4"}
	got := org.Render(v)
	want := "v/videos/video/deadbeef/My File.mp4"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestOrganization_Render_StripsLeadingSlash(t *testing.T) {
	org := &Organization{Template: "/library/$TYPE/$ID/video.mp4"}
	v := &Video{ID: "id1", Type: "v", Filename: "raw.mp4"}
	if got := org.Render(v); got != "library/videos/id1/video.mp4" {
		t.Fatalf("Render() = %q", got)
	}
}

func TestOrganization_Render_HonorsFilename(t *testing.T) {
	org := &Organization{Template: "$TYPE/$FILENAME"}
	v := &Video{ID: "id1", Type: "v", Filename: "subdir/raw.mp4"}
	got := org.Render(v)
	if got != "videos/raw.mp4" {
		t.Fatalf("Render() = %q, want filename basename", got)
	}
}
