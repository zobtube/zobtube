package provider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- Basic metadata tests ---

func TestIAFD_BasicMethods(t *testing.T) {
	p := &IAFD{}

	if slug := p.SlugGet(); slug != "iafd" {
		t.Errorf("expected slug 'iafd', got %q", slug)
	}
	if name := p.NiceName(); name != "IAFD" {
		t.Errorf("expected NiceName 'IAFD', got %q", name)
	}
	if !p.CapabilitySearchActor() {
		t.Error("expected CapabilitySearchActor() to return true")
	}
	if !p.CapabilityScrapePicture() {
		t.Error("expected CapabilityScrapePicture() to return true")
	}
}

// --- ActorSearch tests ---

func TestIAFD_ActorSearch_OfflineMode(t *testing.T) {
	p := &IAFD{}
	_, err := p.ActorSearch(true, "Alice")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Fatalf("expected ErrOfflineMode, got %v", err)
	}
}

func TestIAFD_ActorSearch_SingleMatch(t *testing.T) {
	html := `<tr><td><a href="/person.rme/perfid=alice_1/gender=f/alice"></a></td></tr>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		url, err := p.ActorSearch(false, "Alice")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(url, "/person.rme/") {
			t.Errorf("unexpected url returned: %s", url)
		}
	})
}

func TestIAFD_ActorSearch_NoMatch(t *testing.T) {
	html := `<html><body>No results</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		_, err := p.ActorSearch(false, "Unknown")
		if err == nil {
			t.Fatal("expected error for no match")
		}
	})
}

func TestIAFD_ActorSearch_MultipleMatches(t *testing.T) {
	html := `<tr><td><a href="/person.rme/a1"></a></td></tr><tr><td><a href="/person.rme/a2"></a></td></tr>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		_, err := p.ActorSearch(false, "Alice")
		if err == nil || !strings.Contains(err.Error(), "more than one") {
			t.Fatalf("expected multiple match error, got %v", err)
		}
	})
}

func TestIAFD_ActorSearch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		_, err := p.ActorSearch(false, "Alice")
		if err == nil {
			t.Fatal("expected error on HTTP 404")
		}
	})
}

// --- ActorGetThumb tests ---

func TestIAFD_ActorGetThumb_OfflineMode(t *testing.T) {
	p := &IAFD{
		client: &http.Client{},
	}
	_, err := p.ActorGetThumb(true, "Alice", "")
	if err == nil || !errors.Is(err, ErrOfflineMode) {
		t.Fatalf("expected ErrOfflineMode, got %v", err)
	}
}

func TestIAFD_ActorGetThumb_Success(t *testing.T) {
	// Step 1: mock thumb image server
	thumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte("thumb-data")) //nolint:all
	}))
	defer thumbServer.Close()

	// Step 2: mock profile page
	html := `<div id="headshot"><img src="` + thumbServer.URL + `"></div>`
	pageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer pageServer.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(pageServer, thumbServer, []string{"www.iafd.com"}, nil, func() {
		data, err := p.ActorGetThumb(false, "Alice", pageServer.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(data) != "thumb-data" {
			t.Errorf("expected 'thumb-data', got %q", string(data))
		}
	})
}

func TestIAFD_ActorGetThumb_NoImageFound(t *testing.T) {
	html := `<html><body>No headshot here</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//revive:disable-next-line
		w.Write([]byte(html)) //nolint:all
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		_, err := p.ActorGetThumb(false, "Alice", server.URL)
		if err == nil {
			t.Fatal("expected error for missing thumbnail")
		}
	})
}

func TestIAFD_ActorGetThumb_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := &IAFD{
		client: &http.Client{},
	}
	testHTTPTransports(server, nil, []string{"www.iafd.com"}, nil, func() {
		_, err := p.ActorGetThumb(false, "Alice", server.URL)
		if err == nil {
			t.Fatal("expected error for HTTP 404")
		}
	})
}
