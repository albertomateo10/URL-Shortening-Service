package handler

import (
	"net/http"
	"strconv"

	"github.com/albertomateo10/url-shortener/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) GetClicksOverTime(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "7d"
	}

	resp, err := h.svc.GetClicksOverTime(r.Context(), id, period)
	if err != nil {
		if err.Error() == "url not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "url not found"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AnalyticsHandler) GetSources(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "7d"
	}

	resp, err := h.svc.GetSources(r.Context(), id, period)
	if err != nil {
		if err.Error() == "url not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "url not found"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
