package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type SegmentType string

const (
	SegmentTypeStatic      SegmentType = "static"
	SegmentTypeDynamic     SegmentType = "dynamic"
	SegmentTypeDynamicRule SegmentType = "dynamic_rule"
)

type Segment struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	SegmentName string          `db:"segmentName" json:"segmentName"`
	Type        SegmentType     `db:"segmentType" json:"segmentType"`
	Config      json.RawMessage `db:"config" json:"config"`
	Description string          `db:"description" json:"description"`
	IsActive    bool            `db:"isActive" json:"isActive"`
	CreatedOn   time.Time       `db:"createdOn" json:"createdOn"`
}
