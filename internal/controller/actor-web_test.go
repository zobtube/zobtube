package controller

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func TestController_ActorWeb_ActorList(t *testing.T) {
	ctrl := setupActorController(t)

	// seed test data
	ctrl.datastore.Create(&model.Actor{Name: "Alice"})
	ctrl.datastore.Create(&model.Actor{Name: "Bob"})

	// set gin test mode
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	// setup user
	user := &model.User{}
	ctrl.GetFirstUser(user)
	ctx.Set("user", user)

	tmpl := template.New("actor/list.html")
	tmpl, err := tmpl.Parse(`
		{{- range .Actors }}
		<div class="actor">{{ .Name }}</div>
		{{- end }}
	`)
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}
	engine.SetHTMLTemplate(tmpl)

	ctrl.ActorList(ctx)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	out := w.Body.String()
	if !strings.Contains(out, "Alice") {
		t.Fatalf("expected output to contain actor 'Alice', got: %s", out)
	}
	if !strings.Contains(out, "Bob") {
		t.Fatalf("expected output to contain actor 'Bob', got: %s", out)
	}
}
