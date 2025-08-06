package dto

import (
	"github.com/google/uuid"
)

// AssignUserRequest — payload для POST /segments/{id}/users
type AssignUserRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

// UserSegmentsResponse — ответ на GET /users/{user_id}/segments
type UserSegmentsResponse struct {
	UserID   uuid.UUID         `json:"user_id"`
	Segments []SegmentResponse `json:"segments"`
}

// SegmentUsersResponse — ответ на GET /segments/{id}/users
type SegmentUsersResponse struct {
	SegmentID uuid.UUID   `json:"segment_id"`
	UserIDs   []uuid.UUID `json:"user_ids"`
}
