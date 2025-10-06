package provider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- Basic metadata tests ---

func TestPornhub_BasicMethods(t *testing.T) {
	p := &Pornhub{}

	if slug := p.SlugGet(); slug != "pornhub" {
		t.Errorf("expected slug 'pornhub', got %q", slug)
	}
	if name := p.NiceName(); name != "PornHub" {
		t.Errorf("expected NiceName 'PornHub', got %q", name)
	}
	if !p.CapabilitySearchActor() {
		t.Error("expected CapabilitySearchActor() to return true")
	}
	if !p.CapabilityScrapePicture() {
		t.Error("expected CapabilityScrapePicture() to return true")
	}
}

// --- ActorSearch tests ---

func TestPornhub_ActorSearch_OfflineMode(t *testing.T) {
	p := &Pornhub{}
	_, err := p.ActorSearch(true, "Jane Doe")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestPornhub_ActorSearch_PornstarFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "/pornstar/") {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Pornhub{}
	testHTTPTransports(server, nil, []string{"www.pornhub.com"}, nil, func() {
		url, err := p.ActorSearch(false, "Jane Doe")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "pornstar/jane-doe") {
			t.Errorf("unexpected url: %s", url)
		}
	})
}

func TestPornhub_ActorSearch_ModelFallback(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if strings.Contains(r.URL.Path, "/pornstar/") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(r.URL.Path, "/model/") {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	p := &Pornhub{}
	testHTTPTransports(server, nil, []string{"www.pornhub.com"}, nil, func() {
		url, err := p.ActorSearch(false, "Jane Doe")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "model/jane-doe") {
			t.Errorf("unexpected fallback url: %s", url)
		}
		if requestCount < 2 {
			t.Errorf("expected two HTTP requests, got %d", requestCount)
		}
	})
}

func TestPornhub_ActorSearch_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Pornhub{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorSearch(false, "Unknown Person")
		if err == nil {
			t.Fatal("expected error for not found actor")
		}
	})
}

// --- ActorGetThumb tests ---

func TestPornhub_ActorGetThumb_OfflineMode(t *testing.T) {
	p := &Pornhub{}
	_, err := p.ActorGetThumb(true, "Jane Doe", "")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestPornhub_ActorGetThumb_LegacyImage(t *testing.T) {
	// Step 1: mock the thumb image server
	thumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("thumb-bytes")) //nolint:all
	}))
	defer thumbServer.Close()

	// Step 2: mock the profile page with <img id="getAvatar">
	html := `<img id="getAvatar" src="` + thumbServer.URL + `">`
	pageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer pageServer.Close()

	p := &Pornhub{}
	testHTTPTransports(pageServer, thumbServer, nil, nil, func() {
		data, err := p.ActorGetThumb(false, "Jane Doe", pageServer.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(data) != "thumb-bytes" {
			t.Errorf("expected 'thumb-bytes', got %q", string(data))
		}
	})
}

func TestPornhub_ActorGetThumb_NewImageFormat(t *testing.T) {
	// Step 1: mock the thumb image server
	thumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("new-thumb")) //nolint:all
	}))
	defer thumbServer.Close()

	// Step 2: mock the profile page with <div class="thumbImage">
	html := `<div class="thumbImage"><img src="` + thumbServer.URL + `">`
	pageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer pageServer.Close()

	p := &Pornhub{}
	testHTTPTransports(pageServer, thumbServer, nil, nil, func() {
		data, err := p.ActorGetThumb(false, "Jane Doe", pageServer.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(data) != "new-thumb" {
			t.Errorf("expected 'new-thumb', got %q", string(data))
		}
	})
}

func TestPornhub_ActorGetThumb_NoImageFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("<html><body>no image here</body></html>")) //nolint:all
	}))
	defer server.Close()

	p := &Pornhub{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Jane Doe", server.URL)
		if err == nil {
			t.Fatal("expected error for missing thumbnail")
		}
	})
}

func TestPornhub_ActorGetThumb_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Pornhub{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Jane Doe", server.URL)
		if err == nil {
			t.Fatal("expected error for HTTP 404")
		}
	})
}
