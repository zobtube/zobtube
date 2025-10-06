package provider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- Simple metadata tests ---

func TestBabesDirectory_BasicMethods(t *testing.T) {
	p := &BabesDirectory{}

	if slug := p.SlugGet(); slug != "babesdirectory" {
		t.Errorf("expected slug 'babesdirectory', got %q", slug)
	}
	if name := p.NiceName(); name != "Babes Directory" {
		t.Errorf("expected NiceName 'Babes Directory', got %q", name)
	}
	if !p.CapabilitySearchActor() {
		t.Error("expected CapabilitySearchActor to return true")
	}
	if !p.CapabilityScrapePicture() {
		t.Error("expected CapabilityScrapePicture to return true")
	}
}

// --- ActorSearch tests ---

func TestBabesDirectory_ActorSearch_OfflineMode(t *testing.T) {
	p := &BabesDirectory{}
	_, err := p.ActorSearch(true, "Alice Example")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBabesDirectory_ActorSearch_FirstTry200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &BabesDirectory{}
	testHTTPTransports(server, nil, []string{"babesdirectory.online"}, nil, func() {
		url, err := p.ActorSearch(false, "Alice Example")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "babesdirectory.online/profile/alice-example") {
			t.Errorf("unexpected URL: %s", url)
		}
	})
}

func TestBabesDirectory_ActorSearch_SecondTry200(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &BabesDirectory{}
	testHTTPTransports(server, nil, []string{"babesdirectory.online"}, nil, func() {
		url, err := p.ActorSearch(false, "Alice Example")
		if err != nil {
			t.Fatalf("expected no error on second try, got %v", err)
		}
		if !strings.Contains(url, "-pornstar") {
			t.Errorf("expected '-pornstar' suffix in URL, got %q", url)
		}
	})
}

func TestBabesDirectory_ActorSearch_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &BabesDirectory{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorSearch(false, "No Match")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// --- ActorGetThumb tests ---

func TestBabesDirectory_ActorGetThumb_OfflineMode(t *testing.T) {
	p := &BabesDirectory{}
	_, err := p.ActorGetThumb(true, "Alice Example", "")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBabesDirectory_ActorGetThumb_Success(t *testing.T) {
	// Step 1: serve HTML with <img src="thumb.jpg">
	serverThumb := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("thumb-bytes")) //nolint:all
	}))
	defer serverThumb.Close()

	html := `<div class="pill-image"><img src="` + serverThumb.URL + `">`

	// Step 2: first call returns the HTML
	serverPage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer serverPage.Close()

	p := &BabesDirectory{}
	testHTTPTransports(serverPage, serverThumb, nil, nil, func() {
		data, err := p.ActorGetThumb(false, "Alice Example", serverPage.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(data) != "thumb-bytes" {
			t.Errorf("expected 'thumb-bytes', got %q", string(data))
		}
	})
}

func TestBabesDirectory_ActorGetThumb_NoThumbFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("<html><body>no image here</body></html>")) //nolint:all
	}))
	defer server.Close()

	p := &BabesDirectory{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Alice Example", server.URL)
		if err == nil {
			t.Fatal("expected error for missing thumbnail")
		}
	})
}

func TestBabesDirectory_ActorGetThumb_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &BabesDirectory{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Alice Example", server.URL)
		if err == nil {
			t.Fatal("expected error for HTTP 404")
		}
	})
}
