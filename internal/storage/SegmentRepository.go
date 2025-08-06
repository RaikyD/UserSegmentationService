package storage

import (
	"context"
	"time"

	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentDB struct {
	pool *pgxpool.Pool
}

// NewSegmentDB конструирует репозиторий сегментов на основе пула соединений.
func NewSegmentDB(pool *pgxpool.Pool) *SegmentDB {
	return &SegmentDB{pool: pool}
}

func (db *SegmentDB) Create(ctx context.Context, seg *models.Segment) error {
	seg.ID = uuid.New()
	seg.CreatedOn = time.Now()

	const sql = `
INSERT INTO segments
  (id, segment_name, type, config, description, is_active, created_on)
VALUES
  ($1, $2, $3, $4, $5, $6, $7);
`
	_, err := db.pool.Exec(ctx, sql,
		seg.ID,
		seg.SegmentName,
		seg.Type,
		seg.Config,
		seg.Description,
		seg.IsActive,
		seg.CreatedOn,
	)
	return err
}

func (db *SegmentDB) GetByID(ctx context.Context, id uuid.UUID) (*models.Segment, error) {
	const sql = `
SELECT id, segment_name, type, config, description, is_active, created_on
  FROM segments
 WHERE id = $1;
`
	row := db.pool.QueryRow(ctx, sql, id)
	var seg models.Segment
	if err := row.Scan(
		&seg.ID,
		&seg.SegmentName,
		&seg.Type,
		&seg.Config,
		&seg.Description,
		&seg.IsActive,
		&seg.CreatedOn,
	); err != nil {
		return nil, err
	}
	return &seg, nil
}

func (db *SegmentDB) List(ctx context.Context) ([]*models.Segment, error) {
	const sql = `
SELECT id, segment_name, type, config, description, is_active, created_on
  FROM segments
 ORDER BY created_on DESC;
`
	rows, err := db.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Segment
	for rows.Next() {
		var seg models.Segment
		if err := rows.Scan(
			&seg.ID,
			&seg.SegmentName,
			&seg.Type,
			&seg.Config,
			&seg.Description,
			&seg.IsActive,
			&seg.CreatedOn,
		); err != nil {
			return nil, err
		}
		out = append(out, &seg)
	}
	return out, rows.Err()
}

func (db *SegmentDB) Update(ctx context.Context, seg *models.Segment) error {
	const sql = `
UPDATE segments
   SET segment_name = $2,
       type         = $3,
       config       = $4,
       description  = $5,
       is_active    = $6
 WHERE id = $1;
`
	_, err := db.pool.Exec(ctx, sql,
		seg.ID,
		seg.SegmentName,
		seg.Type,
		seg.Config,
		seg.Description,
		seg.IsActive,
	)
	return err
}

func (db *SegmentDB) Delete(ctx context.Context, id uuid.UUID) error {
	const sql = `
DELETE FROM segments
 WHERE id = $1;
`
	_, err := db.pool.Exec(ctx, sql, id)
	return err
}
