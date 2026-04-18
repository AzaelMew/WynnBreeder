package handlers

import (
	"net/http"

	"wynnbreeder/models"
)

func (h *Handler) AnalyticsPage(w http.ResponseWriter, r *http.Request) {
	statRows, err := h.DB.GetStatInheritance()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	colorRows, err := h.DB.GetColorInheritance()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	potStats, err := h.DB.GetPotentialStats()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	h.render(w, r, "analytics.html", PageData{
		Title: "Analytics",
		Data: map[string]any{
			"StatRows":      statRows,
			"ColorRows":     colorRows,
			"PotentialStats": potStats,
		},
	})
}

func (h *Handler) APIAnalyticsStats(w http.ResponseWriter, r *http.Request) {
	statRows, err := h.DB.GetStatInheritance()
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	colorRows, err := h.DB.GetColorInheritance()
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	potStats, err := h.DB.GetPotentialStats()
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]any{
		"stat_inheritance": statRows,
		"color_inheritance": colorRows,
		"potential":        potStats,
	})
}

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	subs, total, err := h.DB.ListSubmissions(5, 0)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	user := UserFromContext(r.Context())
	var pendingSubs []models.Submission
	if user != nil {
		pendingSubs, err = h.DB.ListPendingByUser(user.ID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
	}

	h.render(w, r, "dashboard.html", PageData{
		Title: "Dashboard",
		Data: map[string]any{
			"RecentSubmissions": subs,
			"TotalSubmissions":  total,
			"PendingBreedings":  pendingSubs,
		},
	})
}
