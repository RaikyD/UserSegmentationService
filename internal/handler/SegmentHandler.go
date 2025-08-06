package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/RaikyD/UserSegmentationService/internal/handler/dto"
	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/RaikyD/UserSegmentationService/internal/service"
)

// SegmentHandler обрабатывает HTTP-запросы, связанные с сегментами.
type SegmentHandler struct {
	svc service.SegmentService
}

// NewSegmentHandler создаёт новый HTTP-хэндлер для сегментов.
func NewSegmentHandler(svc service.SegmentService) *SegmentHandler {
	return &SegmentHandler{svc: svc}
}

// CreateSegment обрабатывает POST /segments
func (h *SegmentHandler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	seg := &models.Segment{
		SegmentName: req.Name,
		Type:        models.SegmentType(req.Type),
		Config:      req.Config,
		Description: "",
		IsActive:    true,
		CreatedOn:   time.Now(),
	}
	if req.Description != nil {
		seg.Description = *req.Description
	}
	if req.IsActive != nil {
		seg.IsActive = *req.IsActive
	}
	created, err := h.svc.CreateSegment(r.Context(), seg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.SegmentResponse{
		ID:          created.ID,
		Name:        created.SegmentName,
		Type:        string(created.Type),
		Config:      created.Config,
		Description: &created.Description,
		IsActive:    created.IsActive,
		CreatedOn:   created.CreatedOn,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// ListSegments обрабатывает GET /segments
func (h *SegmentHandler) ListSegments(w http.ResponseWriter, r *http.Request) {
	segments, err := h.svc.ListSegments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp []dto.SegmentResponse
	for _, s := range segments {
		resp = append(resp, dto.SegmentResponse{
			ID:          s.ID,
			Name:        s.SegmentName,
			Type:        string(s.Type),
			Config:      s.Config,
			Description: &s.Description,
			IsActive:    s.IsActive,
			CreatedOn:   s.CreatedOn,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetSegment обрабатывает GET /segments/{id}
func (h *SegmentHandler) GetSegment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	seg, err := h.svc.GetSegmentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if seg == nil {
		http.Error(w, "segment not found", http.StatusNotFound)
		return
	}
	resp := dto.SegmentResponse{
		ID:          seg.ID,
		Name:        seg.SegmentName,
		Type:        string(seg.Type),
		Config:      seg.Config,
		Description: &seg.Description,
		IsActive:    seg.IsActive,
		CreatedOn:   seg.CreatedOn,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateSegment обрабатывает PUT /segments/{id}
func (h *SegmentHandler) UpdateSegment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	// Fetch existing
	existing, err := h.svc.GetSegmentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "segment not found", http.StatusNotFound)
		return
	}
	if req.Name != nil {
		existing.SegmentName = *req.Name
	}
	if req.Type != nil {
		existing.Type = models.SegmentType(*req.Type)
	}
	if req.Config != nil {
		existing.Config = *req.Config
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	updated, err := h.svc.UpdateSegment(r.Context(), existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.SegmentResponse{
		ID:          updated.ID,
		Name:        updated.SegmentName,
		Type:        string(updated.Type),
		Config:      updated.Config,
		Description: &updated.Description,
		IsActive:    updated.IsActive,
		CreatedOn:   updated.CreatedOn,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteSegment обрабатывает DELETE /segments/{id}
func (h *SegmentHandler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid segment id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteSegment(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Register регистрирует маршруты сегментов в роутере
func (h *SegmentHandler) Register(r chi.Router) {
	r.Post("/", h.CreateSegment)
	r.Get("/", h.ListSegments)
	r.Get("/{id}", h.GetSegment)
	r.Put("/{id}", h.UpdateSegment)
	r.Delete("/{id}", h.DeleteSegment)
}
