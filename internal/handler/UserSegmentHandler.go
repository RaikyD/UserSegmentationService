package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RaikyD/UserSegmentationService/internal/handler/dto"
	"github.com/RaikyD/UserSegmentationService/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserSegmentHandler struct {
	svc service.UserSegmentService
}

func NewUserSegmentHandler(svc service.UserSegmentService) *UserSegmentHandler {
	return &UserSegmentHandler{svc: svc}
}

func (h *UserSegmentHandler) Register(r chi.Router) {
	r.Post("/segments/{segmentID}/users", h.AssignUser)
	r.Delete("/segments/{segmentID}/users/{userID}", h.UnassignUser)
	r.Get("/users/{userID}/segments", h.ListUserSegments)
	r.Get("/segments/{segmentID}/users", h.ListSegmentUsers)
	r.Post("/segments/mass-assign", h.MassAssignSegment)
}

func (h *UserSegmentHandler) AssignUser(w http.ResponseWriter, r *http.Request) {
	segmentID, err := uuid.Parse(chi.URLParam(r, "segmentID"))
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	var req dto.AssignUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("AssignUser: req.UserID = %s", req.UserID)

	if err := h.svc.AssignUser(r.Context(), segmentID, req.UserID); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "segment not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserSegmentHandler) UnassignUser(w http.ResponseWriter, r *http.Request) {
	segmentID, err := uuid.Parse(chi.URLParam(r, "segmentID"))
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if err := h.svc.UnassignUser(r.Context(), segmentID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserSegmentHandler) ListUserSegments(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	segs, err := h.svc.ListUserSegments(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp dto.UserSegmentsResponse
	resp.UserID = userID
	for _, s := range segs {
		resp.Segments = append(resp.Segments, dto.SegmentResponse{
			ID:          s.ID,
			Name:        s.SegmentName,
			Type:        string(s.Type),
			Config:      s.Config,
			Description: &s.Description,
			IsActive:    s.IsActive,
			CreatedOn:   s.CreatedOn,
		})
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserSegmentHandler) ListSegmentUsers(w http.ResponseWriter, r *http.Request) {
	segmentID, err := uuid.Parse(chi.URLParam(r, "segmentID"))
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	users, err := h.svc.ListSegmentUsers(r.Context(), segmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.SegmentUsersResponse{
		SegmentID: segmentID,
		UserIDs:   users,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserSegmentHandler) MassAssignSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.MassAssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Percent < 1 || req.Percent > 100 {
		http.Error(w, "percent must be between 1 and 100", http.StatusBadRequest)
		return
	}

	result, err := h.svc.MassAssignSegment(r.Context(), req.SegmentID, req.Percent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
