package http

import (
	"embed"
	"html/template"
	"io/fs"
	"regexp"
	"strings"
)

func loadAndAddToRoot(
	funcMap template.FuncMap,
	rootTemplate *template.Template,
	embedFS embed.FS,
	pattern string,
) error {
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".*")

	err := fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if matched, _ := regexp.MatchString(pattern, path); !d.IsDir() && matched {
			data, readErr := embedFS.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			t := rootTemplate.New(path).Funcs(funcMap)
			if _, parseErr := t.Parse(string(data)); parseErr != nil {
				return parseErr
			}
		}
		return nil
	})
	return err
}

func (s *Server) LoadHTMLFromEmbedFS(globPath string) {
	root := template.New("")
	tmpl := template.Must(
		root,
		loadAndAddToRoot(
			s.Router.FuncMap,
			root,
			*s.FS,
			globPath),
	)
	s.Router.SetHTMLTemplate(tmpl)
}
