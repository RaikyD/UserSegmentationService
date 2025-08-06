package dto

import "github.com/google/uuid"

// MassAssignRequest запрос на массовое назначение сегмента
type MassAssignRequest struct {
	SegmentID uuid.UUID `json:"segmentID"`
	Percent   int       `json:"percent"` // от 1 до 100
}

// MassAssignResult результат массового назначения сегмента
type MassAssignResult struct {
	TotalUsers int `json:"totalUsers"`
	Assigned   int `json:"assigned"`
	Skipped    int `json:"skipped"`
} 