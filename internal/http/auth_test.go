package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_authRedirectURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("path_only", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "http://example.com/profile", nil)
		got := authRedirectURL(c)
		want := "/auth?next=%2Fprofile"
		if got != want {
			t.Errorf("authRedirectURL() = %q, want %q", got, want)
		}
	})

	t.Run("path_with_query", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "http://example.com/profile?foo=bar", nil)
		got := authRedirectURL(c)
		want := "/auth?next=%2Fprofile%3Ffoo%3Dbar"
		if got != want {
			t.Errorf("authRedirectURL() = %q, want %q", got, want)
		}
	})
}

func Test_apiUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/foo", nil)

	apiUnauthorized(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("apiUnauthorized() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if body["error"] != "unauthorized" {
		t.Errorf("body[error] = %q, want %q", body["error"], "unauthorized")
	}
}
