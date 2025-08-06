package storage

import (
	"context"
	"log"

	"github.com/RaikyD/UserSegmentationService/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDB struct {
	pool *pgxpool.Pool
}

func (db *UserDB) Add(ctx context.Context, asg *models.UserSegmentAssignment) error {
	log.Printf("Repo.Add: asg.UserID = %s", asg.UserID)

	const sql = `
		INSERT INTO user_segment_assignment
		    (segment_id, user_id, assignment_type, assigned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (segment_id, user_id) DO UPDATE 
		SET assignment_type = EXCLUDED.assignment_type,
      	assigned_at     = EXCLUDED.assigned_at;
	`
	_, err := db.pool.Exec(ctx, sql,
		asg.SegmentID,
		asg.UserID,
		asg.AssignmentType,
		asg.AssignedAt)
	return err
}

func (db *UserDB) Delete(ctx context.Context, segmentID, userID uuid.UUID) error {
	const sql = `
DELETE FROM user_segment_assignment
 WHERE segment_id = $1
   AND user_id    = $2;
`
	_, err := db.pool.Exec(ctx, sql, segmentID, userID)
	return err
}

func (db *UserDB) ListByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserSegmentAssignment, error) {
	const sql = `
SELECT segment_id, user_id, assignment_type, assigned_at
  FROM user_segment_assignment
 WHERE user_id = $1;
`
	rows, err := db.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.UserSegmentAssignment
	for rows.Next() {
		var asg models.UserSegmentAssignment
		if err := rows.Scan(
			&asg.SegmentID,
			&asg.UserID,
			&asg.AssignmentType,
			&asg.AssignedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &asg)
	}
	return list, rows.Err()
}

func (db *UserDB) ListBySegment(ctx context.Context, segmentID uuid.UUID) ([]*models.UserSegmentAssignment, error) {
	const sql = `
SELECT segment_id, user_id, assignment_type, assigned_at
  FROM user_segment_assignment
 WHERE segment_id = $1;
`
	rows, err := db.pool.Query(ctx, sql, segmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.UserSegmentAssignment
	for rows.Next() {
		var asg models.UserSegmentAssignment
		if err := rows.Scan(
			&asg.SegmentID,
			&asg.UserID,
			&asg.AssignmentType,
			&asg.AssignedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &asg)
	}
	return list, rows.Err()
}

// GetAllUserIDs получает все уникальные user_id из таблицы user_segment
func (db *UserDB) GetAllUserIDs(ctx context.Context) ([]uuid.UUID, error) {
	const sql = `
		SELECT DISTINCT user_id 
		FROM user_segment_assignment
		ORDER BY user_id;
	`
	rows, err := db.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

func NewUserDB(pool *pgxpool.Pool) *UserDB {
	return &UserDB{pool: pool}
}
