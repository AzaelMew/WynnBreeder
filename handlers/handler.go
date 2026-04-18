package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"wynnmounts/database"
	"wynnmounts/models"
)

type Handler struct {
	DB         *database.DB
	// Per-page templates: map[filename]*template.Template, each pre-parsed with layout.
	Templates  map[string]*template.Template
	SessionTTL time.Duration
}

type PageData struct {
	User      *models.User
	Title     string
	Data      any
	Flash     string
	FlashType string
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, tmpl string, data PageData) {
	if data.User == nil {
		data.User = UserFromContext(r.Context())
	}
	t, ok := h.Templates[tmpl]
	if !ok {
		http.Error(w, "template not found: "+tmpl, http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
