package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateSegmentRequest — payload для POST /segments
type CreateSegmentRequest struct {
	Name        string          `json:"name"        validate:"required"`
	Type        string          `json:"type"        validate:"required,oneof=static dynamic dynamic_rule"`
	Config      json.RawMessage `json:"config"      validate:"required,json"`
	Description *string         `json:"description"` // опционально
	IsActive    *bool           `json:"is_active"`   // опционально, default=true
	ValidFrom   *time.Time      `json:"valid_from"`  // опционально
	ValidTo     *time.Time      `json:"valid_to"`    // опционально
}

// UpdateSegmentRequest — payload для PUT /segments/{id}
type UpdateSegmentRequest struct {
	Name        *string          `json:"name"` // опционально
	Type        *string          `json:"type"        validate:"omitempty,oneof=static dynamic dynamic_rule"`
	Config      *json.RawMessage `json:"config"      validate:"omitempty,json"`
	Description *string          `json:"description"`
	IsActive    *bool            `json:"is_active"`
	ValidFrom   *time.Time       `json:"valid_from"`
	ValidTo     *time.Time       `json:"valid_to"`
}

// SegmentResponse — то, что возвращаем клиенту по GET /segments[/{id}]
type SegmentResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Config      json.RawMessage `json:"config"`
	Description *string         `json:"description,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedOn   time.Time       `json:"created_on"`
	ValidFrom   *time.Time      `json:"valid_from,omitempty"`
	ValidTo     *time.Time      `json:"valid_to,omitempty"`
}
