package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"wynnbreeder/database"
	"wynnbreeder/models"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) SubmitPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "submit.html", PageData{Title: "Submit Breeding"})
}

func (h *Handler) SubmissionsPage(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	if page < 1 {
		page = 1
	}
	const perPage = 20
	offset := (page - 1) * perPage

	subs, total, err := h.DB.ListSubmissions(perPage, offset)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	h.render(w, r, "submissions.html", PageData{
		Title: "Submissions",
		Data: map[string]any{
			"Submissions": subs,
			"Total":       total,
			"Page":        page,
			"TotalPages":  totalPages,
		},
	})
}

type StatDelta struct {
	Name        string
	ParentA     int
	ParentALim  int
	ParentAMax  int
	ParentB     int
	ParentBLim  int
	ParentBMax  int
	Avg         float64
	Offspring   int
	OffspringLim int
	OffspringMax int
	Delta       float64
}

func computeDeltas(pa, pb, off *models.Mount) []StatDelta {
	if pa == nil || pb == nil || off == nil {
		return nil
	}
	d := func(name string, aVal, aLim, aMax, bVal, bLim, bMax, oVal, oLim, oMax int) StatDelta {
		avg := float64(aVal+bVal) / 2.0
		return StatDelta{
			Name: name,
			ParentA: aVal, ParentALim: aLim, ParentAMax: aMax,
			ParentB: bVal, ParentBLim: bLim, ParentBMax: bMax,
			Avg:       avg,
			Offspring: oVal, OffspringLim: oLim, OffspringMax: oMax,
			Delta: float64(oVal) - avg,
		}
	}
	return []StatDelta{
		d("Speed", pa.SpeedVal, pa.SpeedLim, pa.SpeedMax, pb.SpeedVal, pb.SpeedLim, pb.SpeedMax, off.SpeedVal, off.SpeedLim, off.SpeedMax),
		d("Acceleration", pa.AccelVal, pa.AccelLim, pa.AccelMax, pb.AccelVal, pb.AccelLim, pb.AccelMax, off.AccelVal, off.AccelLim, off.AccelMax),
		d("Altitude", pa.AltitudeVal, pa.AltitudeLim, pa.AltitudeMax, pb.AltitudeVal, pb.AltitudeLim, pb.AltitudeMax, off.AltitudeVal, off.AltitudeLim, off.AltitudeMax),
		d("Energy", pa.EnergyStatVal, pa.EnergyStatLim, pa.EnergyStatMax, pb.EnergyStatVal, pb.EnergyStatLim, pb.EnergyStatMax, off.EnergyStatVal, off.EnergyStatLim, off.EnergyStatMax),
		d("Handling", pa.HandlingVal, pa.HandlingLim, pa.HandlingMax, pb.HandlingVal, pb.HandlingLim, pb.HandlingMax, off.HandlingVal, off.HandlingLim, off.HandlingMax),
		d("Toughness", pa.ToughnessVal, pa.ToughnessLim, pa.ToughnessMax, pb.ToughnessVal, pb.ToughnessLim, pb.ToughnessMax, off.ToughnessVal, off.ToughnessLim, off.ToughnessMax),
		d("Boost", pa.BoostVal, pa.BoostLim, pa.BoostMax, pb.BoostVal, pb.BoostLim, pb.BoostMax, off.BoostVal, off.BoostLim, off.BoostMax),
		d("Training", pa.TrainingVal, pa.TrainingLim, pa.TrainingMax, pb.TrainingVal, pb.TrainingLim, pb.TrainingMax, off.TrainingVal, off.TrainingLim, off.TrainingMax),
		d("Potential", pa.Potential, 0, 0, pb.Potential, 0, 0, off.Potential, 0, 0),
	}
}

func (h *Handler) SubmissionDetailPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	sub, err := h.DB.GetSubmission(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	mountMap := map[models.MountRole]*models.Mount{}
	for i := range sub.Mounts {
		m := &sub.Mounts[i]
		mountMap[m.Role] = m
	}

	pa := mountMap[models.RoleParentA]
	pb := mountMap[models.RoleParentB]
	off := mountMap[models.RoleOffspring]

	user := UserFromContext(r.Context())
	isOwner := user != nil && (user.ID == sub.UserID || user.IsAdmin)

	h.render(w, r, "submission_detail.html", PageData{
		Title: fmt.Sprintf("Submission #%d", sub.ID),
		Data: map[string]any{
			"Submission": sub,
			"ParentA":    pa,
			"ParentB":    pb,
			"Offspring":  off,
			"Deltas":     computeDeltas(pa, pb, off),
			"IsOwner":    isOwner,
		},
	})
}

// APICreateSubmission handles POST /api/submissions
// Offspring is optional — omit it to save a pending (in-progress) breeding.
func (h *Handler) APICreateSubmission(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	var req models.SubmitRequest
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateMountJSON(req.ParentA, "parent_a"); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateMountJSON(req.ParentB, "parent_b"); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	mounts := []models.Mount{
		models.MountFromJSON(req.ParentA, 0, models.RoleParentA),
		models.MountFromJSON(req.ParentB, 0, models.RoleParentB),
	}

	status := "pending"
	if req.Offspring != nil {
		if err := validateMountJSON(*req.Offspring, "offspring"); err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		mounts = append(mounts, models.MountFromJSON(*req.Offspring, 0, models.RoleOffspring))
		status = "complete"
	}

	sub, err := h.DB.CreateSubmission(user.ID, req.Notes, mounts, status)
	if err != nil {
		jsonError(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, sub)
}

// APIAddOffspring handles PATCH /api/submissions/:id/offspring
// Only the owner or an admin may complete a pending submission.
func (h *Handler) APIAddOffspring(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.DB.GetSubmission(id)
	if err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	if !user.IsAdmin && sub.UserID != user.ID {
		jsonError(w, "forbidden — only the submitter can add offspring", http.StatusForbidden)
		return
	}

	var offJSON models.MountJSON
	if err := decodeJSON(r, &offJSON); err != nil {
		jsonError(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateMountJSON(offJSON, "offspring"); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	offspring := models.MountFromJSON(offJSON, id, models.RoleOffspring)
	updated, err := h.DB.AddOffspring(id, offspring)
	if err != nil {
		if err == database.ErrSubmissionAlreadyComplete {
			jsonError(w, "submission already has offspring", http.StatusConflict)
			return
		}
		jsonError(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, updated)
}

// APIListSubmissions handles GET /api/submissions
func (h *Handler) APIListSubmissions(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	if page < 1 {
		page = 1
	}
	const perPage = 20
	subs, total, err := h.DB.ListSubmissions(perPage, (page-1)*perPage)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]any{"submissions": subs, "total": total, "page": page})
}

// APIGetSubmission handles GET /api/submissions/:id
func (h *Handler) APIGetSubmission(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}
	sub, err := h.DB.GetSubmission(id)
	if err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	jsonOK(w, sub)
}

// APIDeleteSubmission handles DELETE /api/submissions/:id
func (h *Handler) APIDeleteSubmission(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.DB.GetSubmission(id)
	if err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	if !user.IsAdmin && sub.UserID != user.ID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.DB.DeleteSubmission(id); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

func validateMountJSON(m models.MountJSON, role string) error {
	if strings.TrimSpace(m.Type) == "" {
		return fmt.Errorf("%s: missing type", role)
	}
	if m.Potential < 0 {
		return fmt.Errorf("%s: potential must be non-negative", role)
	}
	if strings.TrimSpace(m.Color) == "" {
		return fmt.Errorf("%s: missing color", role)
	}
	return nil
}

func decodeJSON(r *http.Request, v any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB max
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func queryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

var _ = database.ErrSubmissionNotFound
