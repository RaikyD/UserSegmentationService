package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/RaikyD/UserSegmentationService/internal/storage"
	"github.com/google/uuid"
)

// UserSegmentService описывает логику работы с привязками пользователей к сегментам.
type UserSegmentService interface {
	AssignUser(ctx context.Context, segmentID, userID uuid.UUID) error
	UnassignUser(ctx context.Context, segmentID, userID uuid.UUID) error
	ListUserSegments(ctx context.Context, userID uuid.UUID) ([]*models.Segment, error)
	ListSegmentUsers(ctx context.Context, segmentID uuid.UUID) ([]uuid.UUID, error)
	MassAssignSegment(ctx context.Context, segmentID uuid.UUID, percent int) (*models.MassAssignResult, error)
}

type userSegmentService struct {
	segRepo storage.SegmentRepository
	usRepo  storage.UserRepository
}

func NewUserSegmentService(segRepo storage.SegmentRepository, usRepo storage.UserRepository) UserSegmentService {
	return &userSegmentService{segRepo: segRepo, usRepo: usRepo}
}

func (u *userSegmentService) AssignUser(ctx context.Context, segmentID, userID uuid.UUID) error {
	log.Printf("Service.AssignUser: userID = %s", userID)

	seg, err := u.segRepo.GetByID(ctx, segmentID)
	if err != nil {
		return err
	}
	if seg == nil {
		return fmt.Errorf("segment %s not found", segmentID)
	}
	asg := &models.UserSegmentAssignment{
		SegmentID:      segmentID,
		UserID:         userID,
		AssignmentType: models.AssignmentManual,
		AssignedAt:     time.Now(),
	}
	return u.usRepo.Add(ctx, asg)
}

func (u *userSegmentService) UnassignUser(ctx context.Context, segmentID, userID uuid.UUID) error {
	return u.usRepo.Delete(ctx, segmentID, userID)
}

func (u *userSegmentService) ListUserSegments(ctx context.Context, userID uuid.UUID) ([]*models.Segment, error) {
	assignments, err := u.usRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var segments []*models.Segment
	for _, asg := range assignments {
		seg, err := u.segRepo.GetByID(ctx, asg.SegmentID)
		if err != nil {
			return nil, err
		}
		if seg != nil {
			segments = append(segments, seg)
		}
	}
	return segments, nil
}

func (u *userSegmentService) ListSegmentUsers(ctx context.Context, segmentID uuid.UUID) ([]uuid.UUID, error) {
	assignments, err := u.usRepo.ListBySegment(ctx, segmentID)
	if err != nil {
		return nil, err
	}
	var userIDs []uuid.UUID
	for _, asg := range assignments {
		userIDs = append(userIDs, asg.UserID)
	}
	return userIDs, nil
}

// Стоит отметить, что ,наверное, для наиболее правдоподобной работы стоило делать 3 таблицу,
// в которой хранились бы пользователи с их данными, но так как это мне показалось нерационально, то решил
// брать просто пользователей, которые сами добавляем в ходе работы.
func (u *userSegmentService) MassAssignSegment(ctx context.Context, segmentID uuid.UUID, percent int) (*models.MassAssignResult, error) {
	seg, err := u.segRepo.GetByID(ctx, segmentID)
	if err != nil {
		return nil, err
	}
	if seg == nil {
		return nil, fmt.Errorf("segment %s not found", segmentID)
	}

	// Собираем все уникаильные ключи юзеров
	userIDs, err := u.usRepo.GetAllUserIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return &models.MassAssignResult{
			TotalUsers: 0,
			Assigned:   0,
			Skipped:    0,
		}, nil
	}

	// Берём наш процент
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIDs), func(i, j int) {
		userIDs[i], userIDs[j] = userIDs[j], userIDs[i]
	})

	count := len(userIDs) * percent / 100
	if count == 0 && percent > 0 {
		count = 1
	}

	selected := userIDs[:count]

	assigned := 0
	skipped := 0

	for _, userID := range selected {
		err := u.AssignUser(ctx, segmentID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
				skipped++
			} else {
				return nil, err
			}
		} else {
			assigned++
		}
	}

	return &models.MassAssignResult{
		TotalUsers: len(userIDs),
		Assigned:   assigned,
		Skipped:    skipped,
	}, nil
}
