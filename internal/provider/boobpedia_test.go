package provider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- Basic metadata tests ---

func TestBoobpedia_BasicMethods(t *testing.T) {
	p := &Boobpedia{}

	if slug := p.SlugGet(); slug != "boobpedia" {
		t.Errorf("expected slug 'boobpedia', got %q", slug)
	}
	if name := p.NiceName(); name != "Boobpedia" {
		t.Errorf("expected NiceName 'Boobpedia', got %q", name)
	}
	if !p.CapabilitySearchActor() {
		t.Error("expected CapabilitySearchActor() to return true")
	}
	if !p.CapabilityScrapePicture() {
		t.Error("expected CapabilityScrapePicture() to return true")
	}
}

// --- ActorSearch tests ---

func TestBoobpedia_ActorSearch_OfflineMode(t *testing.T) {
	p := &Boobpedia{}
	_, err := p.ActorSearch(true, "Jane Doe")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBoobpedia_ActorSearch_HTTP200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &Boobpedia{}
	testHTTPTransports(server, nil, []string{"www.boobpedia.com"}, nil, func() {
		url, err := p.ActorSearch(false, "Jane Doe")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "boobpedia.com/boobs/Jane_Doe") {
			t.Errorf("unexpected url: %s", url)
		}
	})
}

func TestBoobpedia_ActorSearch_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Boobpedia{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorSearch(false, "Missing Person")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// --- ActorGetThumb tests ---

func TestBoobpedia_ActorGetThumb_OfflineMode(t *testing.T) {
	p := &Boobpedia{}
	_, err := p.ActorGetThumb(true, "Jane Doe", "")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Errorf("expected ErrOfflineMode, got %v", err)
	}
}

func TestBoobpedia_ActorGetThumb_Success(t *testing.T) {
	// Step 1: serve actual thumb bytes
	thumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("thumb-data")) //nolint:all
	}))
	defer thumbServer.Close()

	// Step 2: serve HTML page with the correct <img src="">
	html := `<a class="mw-file-description"><img src="` + thumbServer.URL + `">`
	pageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer pageServer.Close()

	p := &Boobpedia{}
	testHTTPTransports(pageServer, thumbServer, nil, []string{"www.boobpedia.com"}, func() {
		data, err := p.ActorGetThumb(false, "Jane Doe", pageServer.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(data) != "thumb-data" {
			t.Errorf("expected thumb data %q, got %q", "thumb-data", string(data))
		}
	})
}

func TestBoobpedia_ActorGetThumb_NoThumbFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("<html><body>no image</body></html>")) //nolint:all
	}))
	defer server.Close()

	p := &Boobpedia{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Jane Doe", server.URL)
		if err == nil {
			t.Fatal("expected error for missing thumb")
		}
	})
}

func TestBoobpedia_ActorGetThumb_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &Boobpedia{}
	testHTTPTransports(server, nil, nil, nil, func() {
		_, err := p.ActorGetThumb(false, "Jane Doe", server.URL)
		if err == nil {
			t.Fatal("expected error for HTTP 404")
		}
	})
}
