package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, struct {
		Error string `json:"error"`
	}{Error: message})
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

func (s *server) renderFragment(w http.ResponseWriter, file string, data any) {
	name := strings.TrimSuffix(file, ".html")
	var tmpl *template.Template

	if s.dev {
		diskFS := os.DirFS("internal/web")
		t, err := template.New("").ParseFS(diskFS, "templates/fragments/"+file)
		if err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl = t
	} else {
		t, ok := s.fragTmpls[file]
		if !ok {
			http.Error(w, "fragment not found: "+file, http.StatusInternalServerError)
			return
		}
		tmpl = t
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
	}
}
