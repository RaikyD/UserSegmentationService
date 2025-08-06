package models

import (
	"time"

	"github.com/google/uuid"
)

type AssignmentType string

const (
	AssignmentManual AssignmentType = "manual"
	AssignmentAuto   AssignmentType = "auto"
)

// UserSegmentAssignment хранит одну запись о том,
// что пользователь user_id в момент assigned_at
// получил сегмент segment_id типом assignmentType.
type UserSegmentAssignment struct {
	SegmentID      uuid.UUID      `db:"segment_id" json:"segment_id"`
	UserID         uuid.UUID      `db:"user_id"    json:"user_id"`
	AssignmentType AssignmentType `db:"assignment_type" json:"assignment_type"` // Необходимо для понимания,
	// был ли пользователь добавлен "руками" или при случайной выборке (Полезно может быть для условной категории Стримеров/VIP и тд)
	AssignedAt time.Time `db:"assigned_at" json:"assigned_at"`
}
