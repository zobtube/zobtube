package provider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBabepedia_BasicMethods(t *testing.T) {
	p := &Babepedia{}

	if slug := p.SlugGet(); slug != "babepedia" {
		t.Errorf("expected slug 'babepedia', got %q", slug)
	}

	if name := p.NiceName(); name != "Babepedia" {
		t.Errorf("expected NiceName 'Babepedia', got %q", name)
	}

	if !p.CapabilitySearchActor() {
		t.Error("expected CapabilitySearchActor to return true")
	}
	if !p.CapabilityScrapePicture() {
		t.Error("expected CapabilityScrapePicture to return true")
	}
}

func TestBabepedia_ActorSearch_OfflineMode(t *testing.T) {
	p := &Babepedia{}

	_, err := p.ActorSearch(true, "Jane Doe")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBabepedia_ActorSearch_HTTPStatus200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &Babepedia{}
	testHTTPTransports(server, nil, []string{"www.babepedia.com"}, nil, func() {
		url, err := p.ActorSearch(false, "Jane Doe")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "babepedia.com/babe/") {
			t.Errorf("unexpected url: %s", url)
		}
	})
}

func TestBabepedia_ActorSearch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Babepedia{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorSearch(false, "Missing Person")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestBabepedia_ActorGetThumb_OfflineMode(t *testing.T) {
	p := &Babepedia{}
	_, err := p.ActorGetThumb(true, "Jane Doe", "")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBabepedia_ActorGetThumb_EmptyURL(t *testing.T) {
	p := &Babepedia{}
	_, err := p.ActorGetThumb(false, "Jane Doe", "")
	if err == nil {
		t.Fatal("expected error for empty url")
	}
	if err != nil && !strings.Contains(err.Error(), "actor name not found in url") {
		t.Errorf("expected 'actor name not found in url', got %v", err)
	}
}

func TestBabepedia_ActorGetThumb_Success(t *testing.T) {
	wantData := []byte("fake_image_data")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write(wantData) //nolint:all
	}))
	defer server.Close()

	p := &Babepedia{}
	testHTTPTransports(server, nil, []string{"www.babepedia.com"}, nil, func() {
		got, err := p.ActorGetThumb(false, "Jane Doe", "https://www.babepedia.com/babe/Jane_Doe")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(got) != string(wantData) {
			t.Errorf("expected thumb data %q, got %q", string(wantData), string(got))
		}
	})
}

func TestBabepedia_ActorGetThumb_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Babepedia{}
	testHTTPTransports(server, nil, []string{"www.babepedia.com"}, nil, func() {
		_, err := p.ActorGetThumb(false, "Missing Person", "https://www.babepedia.com/babe/Missing_Person")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
