package storage

import (
	"context"
	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Add(ctx context.Context, asg *models.UserSegmentAssignment) error
	Delete(ctx context.Context, segmentID, userID uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserSegmentAssignment, error)
	ListBySegment(ctx context.Context, segmentID uuid.UUID) ([]*models.UserSegmentAssignment, error)
	GetAllUserIDs(ctx context.Context) ([]uuid.UUID, error)
}
