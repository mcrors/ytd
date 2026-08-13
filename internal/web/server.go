package web

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/mcrors/ytd/internal/format"
)

//go:embed templates static
var embeddedFiles embed.FS

type Queue interface {
	GetTitle(ctx context.Context, url string) (string, error)
	Enqueue(ctx context.Context, url, targetDir, newName string, format format.Format, title string) (int64, error)
	Cancel(id int64)
}

type server struct {
	queue     Queue
	baseDir   string
	db        *sql.DB
	tmpls     map[string]*template.Template
	fragTmpls map[string]*template.Template
	dev       bool
}

func RegisterRoutes(mux *http.ServeMux, queue Queue, baseDir string, db *sql.DB, dev bool) error {
	s := &server{queue: queue, baseDir: baseDir, db: db, dev: dev}

	if !dev {
		tmpls, err := loadTemplates(embeddedFiles)
		if err != nil {
			return fmt.Errorf("loading templates: %w", err)
		}
		s.tmpls = tmpls

		frags, err := loadFragments(embeddedFiles)
		if err != nil {
			return fmt.Errorf("loading fragments: %w", err)
		}
		s.fragTmpls = frags
	}

	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		return fmt.Errorf("static sub-fs: %w", err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /healthz", s.healthzHandler)
	mux.HandleFunc("GET /readyz", s.readyzHandler)
	mux.HandleFunc("POST /downloads", s.submitHandler)
	mux.HandleFunc("GET /downloads/{id}/status", s.statusHandler)
	mux.HandleFunc("DELETE /downloads/{id}/cancel", s.cancelHandler)
	mux.HandleFunc("GET /downloads/history", s.historyHandler)
	mux.HandleFunc("GET /api/directories", s.getDirectoriesHandler)
	mux.HandleFunc("POST /api/directory", s.createDirectoryHandler)

	return nil
}

func loadHTMLFiles(fsys fs.FS, glob string, extras ...string) (map[string]*template.Template, error) {
	files, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}
	prefix := glob[:strings.LastIndex(glob, "/")+1]
	tmpls := make(map[string]*template.Template, len(files))
	for _, f := range files {
		parseFiles := append(extras, f)
		t, err := template.New("").ParseFS(fsys, parseFiles...)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		tmpls[f[len(prefix):]] = t
	}
	return tmpls, nil
}

func loadTemplates(fsys fs.FS) (map[string]*template.Template, error) {
	return loadHTMLFiles(fsys, "templates/pages/*.html", "templates/layout.html")
}

func loadFragments(fsys fs.FS) (map[string]*template.Template, error) {
	return loadHTMLFiles(fsys, "templates/fragments/*.html")
}

