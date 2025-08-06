package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/RaikyD/UserSegmentationService/internal/storage"
	"github.com/google/uuid"
)

type SegmentService interface {
	CreateSegment(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	GetSegmentByID(ctx context.Context, id uuid.UUID) (*models.Segment, error)
	ListSegments(ctx context.Context) ([]*models.Segment, error)
	UpdateSegment(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	DeleteSegment(ctx context.Context, id uuid.UUID) error
}

type segmentService struct {
	repo storage.SegmentRepository
}

func NewSegmentService(repo storage.SegmentRepository) SegmentService {
	return &segmentService{repo: repo}
}

func (s *segmentService) CreateSegment(ctx context.Context, seg *models.Segment) (*models.Segment, error) {
	seg.CreatedOn = time.Now()
	if err := s.repo.Create(ctx, seg); err != nil {
		return nil, err
	}
	return seg, nil
}

func (s *segmentService) GetSegmentByID(ctx context.Context, id uuid.UUID) (*models.Segment, error) {
	seg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if seg == nil {
		return nil, fmt.Errorf("segment %s not found", id)
	}
	return seg, nil
}

func (s *segmentService) ListSegments(ctx context.Context) ([]*models.Segment, error) {
	return s.repo.List(ctx)
}

func (s *segmentService) UpdateSegment(ctx context.Context, seg *models.Segment) (*models.Segment, error) {
	existing, err := s.repo.GetByID(ctx, seg.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("segment %s not found", seg.ID)
	}
	if err := s.repo.Update(ctx, seg); err != nil {
		return nil, err
	}
	return seg, nil
}

func (s *segmentService) DeleteSegment(ctx context.Context, id uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("segment %s not found", id)
	}
	return s.repo.Delete(ctx, id)
}
