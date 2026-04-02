package handler

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/albertomateo10/url-shortener/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	urlSvc      *service.URLService
	clickLogger *service.ClickLogger
}

func NewRedirectHandler(urlSvc *service.URLService, clickLogger *service.ClickLogger) *RedirectHandler {
	return &RedirectHandler{urlSvc: urlSvc, clickLogger: clickLogger}
}

func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	u, err := h.urlSvc.ResolveShortCode(r.Context(), shortCode)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if u == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "short code not found"})
		return
	}

	// Log click asynchronously
	h.clickLogger.Log(&model.ClickEvent{
		URLID:     u.ID,
		ClickedAt: time.Now(),
		IPAddress: extractIP(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
	})

	http.Redirect(w, r, u.OriginalURL, http.StatusFound)
}

func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
