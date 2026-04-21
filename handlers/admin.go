package handlers

import (
	"net/http"
	"strconv"

	"wynnbreeder/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.ListUsers()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	subs, total, err := h.DB.ListSubmissions(10, 0)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "admin.html", PageData{
		Title: "Admin Panel",
		Data: map[string]any{
			"Users":       users,
			"Submissions": subs,
			"TotalSubs":   total,
		},
	})
}

func (h *Handler) APICreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Username) < 3 {
		jsonError(w, "username must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		jsonError(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	user, err := h.DB.CreateUser(req.Username, string(hash), req.IsAdmin)
	if err != nil {
		if err == database.ErrUsernameTaken {
			jsonError(w, "username already taken", http.StatusConflict)
			return
		}
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, user)
}

func (h *Handler) APIListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.ListUsers()
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, users)
}

func (h *Handler) APIDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	me := UserFromContext(r.Context())

	target, err := h.DB.GetUserByID(id)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	// Prevent self-deletion
	if me != nil && me.ID == id {
		jsonError(w, "cannot delete your own account", http.StatusBadRequest)
		return
	}

	// Only superadmin can delete other admins
	if target.IsAdmin && (me == nil || !me.IsSuperAdmin) {
		jsonError(w, "only the superadmin can delete admin accounts", http.StatusForbidden)
		return
	}

	// Prevent deleting the last superadmin
	if target.IsSuperAdmin {
		jsonError(w, "cannot delete the superadmin account", http.StatusBadRequest)
		return
	}

	if err := h.DB.DeleteUser(id); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

func (h *Handler) APISetUserRole(w http.ResponseWriter, r *http.Request) {
	me := UserFromContext(r.Context())
	if me == nil || !me.IsSuperAdmin {
		jsonError(w, "only superadmin can change roles", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if me.ID == id {
		jsonError(w, "cannot change your own role", http.StatusBadRequest)
		return
	}

	var req struct {
		Role string `json:"role"` // "user", "admin", "superadmin"
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var isAdmin, isSuperAdmin bool
	switch req.Role {
	case "user":
		isAdmin, isSuperAdmin = false, false
	case "admin":
		isAdmin, isSuperAdmin = true, false
	case "superadmin":
		isAdmin, isSuperAdmin = true, true
	default:
		jsonError(w, "role must be user, admin, or superadmin", http.StatusBadRequest)
		return
	}

	updated, err := h.DB.SetUserRole(id, isAdmin, isSuperAdmin)
	if err != nil {
		if err == database.ErrUserNotFound {
			jsonError(w, "user not found", http.StatusNotFound)
			return
		}
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, updated)
}
