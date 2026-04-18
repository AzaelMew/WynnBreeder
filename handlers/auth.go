package handlers

import (
	"net/http"

	"wynnmounts/database"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Already logged in → redirect
	if cookie, err := r.Cookie(cookieName); err == nil {
		if s, err := h.DB.GetSession(cookie.Value); err == nil && s != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	h.render(w, r, "login.html", PageData{Title: "Login"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.render(w, r, "login.html", PageData{Title: "Login", Flash: "Invalid request", FlashType: "error"})
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.DB.GetUserByUsername(username)
	if err != nil {
		h.render(w, r, "login.html", PageData{Title: "Login", Flash: "Invalid credentials", FlashType: "error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.render(w, r, "login.html", PageData{Title: "Login", Flash: "Invalid credentials", FlashType: "error"})
		return
	}

	session, err := h.DB.CreateSession(user.ID, h.SessionTTL)
	if err != nil {
		h.render(w, r, "login.html", PageData{Title: "Login", Flash: "Server error", FlashType: "error"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		_ = h.DB.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// API login (JSON)
func (h *Handler) APILogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.DB.GetUserByUsername(req.Username)
	if err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := h.DB.CreateSession(user.ID, h.SessionTTL)
	if err != nil {
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	jsonOK(w, map[string]any{"ok": true, "is_admin": user.IsAdmin})
}

func (h *Handler) APILogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		_ = h.DB.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, MaxAge: -1, Path: "/"})
	jsonOK(w, map[string]string{"ok": "logged out"})
}

var _ = database.ErrUserNotFound // keep import used
