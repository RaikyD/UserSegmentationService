package storage

import (
	"context"

	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/google/uuid"
)

// SegmentRepository описывает CRUD-операции над сегментами.
type SegmentRepository interface {
	// Create вставляет новый сегмент, заполняя у него ID и CreatedOn
	Create(ctx context.Context, seg *models.Segment) error
	// GetByID возвращает сегмент по его UUID или ошибку
	GetByID(ctx context.Context, id uuid.UUID) (*models.Segment, error)
	// List возвращает все сегменты, отсортированные по CreatedOn
	List(ctx context.Context) ([]*models.Segment, error)
	// Update перезаписывает все изменяемые поля у существующего сегмента
	Update(ctx context.Context, seg *models.Segment) error
	// Delete удаляет сегмент по ID
	Delete(ctx context.Context, id uuid.UUID) error
}
