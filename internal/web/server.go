package web

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/mcrors/ytd/internal/download"
)

//go:embed templates static
var embeddedFiles embed.FS

type Downloader interface {
	Download(context.Context, download.DownloadCommand) (*download.DownloadResult, error)
}

type server struct {
	dl      Downloader
	baseDir string
	db      *sql.DB
	tmpls   map[string]*template.Template
	dev     bool
}

func RegisterRoutes(mux *http.ServeMux, dl Downloader, baseDir string, db *sql.DB, dev bool) error {
	s := &server{dl: dl, baseDir: baseDir, db: db, dev: dev}

	if !dev {
		tmpls, err := loadTemplates(embeddedFiles)
		if err != nil {
			return fmt.Errorf("loading templates: %w", err)
		}
		s.tmpls = tmpls
	}

	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		return fmt.Errorf("static sub-fs: %w", err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /healthz", s.healthzHandler)
	mux.HandleFunc("GET /readyz", s.readyzHandler)
	mux.HandleFunc("POST /api/download", s.downloadHandler)
	mux.HandleFunc("GET /api/directories", s.getDirectoriesHandler)
	mux.HandleFunc("POST /api/directory", s.createDirectoryHandler)

	return nil
}

func loadTemplates(fsys fs.FS) (map[string]*template.Template, error) {
	pages, err := fs.Glob(fsys, "templates/pages/*.html")
	if err != nil {
		return nil, err
	}
	tmpls := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		name := page[len("templates/pages/"):]
		t, err := template.New("").ParseFS(fsys, "templates/layout.html", page)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", page, err)
		}
		tmpls[name] = t
	}
	return tmpls, nil
}

func (s *server) render(w http.ResponseWriter, page string, data any) {
	var tmpl *template.Template

	if s.dev {
		var err error
		diskFS := os.DirFS("internal/web")
		tmpl, err = template.New("").ParseFS(diskFS, "templates/layout.html", "templates/pages/"+page)
		if err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var ok bool
		tmpl, ok = s.tmpls[page]
		if !ok {
			http.Error(w, "template not found: "+page, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
	}
}
