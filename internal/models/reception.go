package models

import (
	"github.com/google/uuid"
	"time"
)

type ReceptionType string

const (
	StatusInProgress ReceptionType = "in_progress"
	StatusClose      ReceptionType = "close"
)

type Reception struct {
	ID       uuid.UUID     `json:"id"`
	DateTime time.Time     `json:"dateTime"`
	PvzID    uuid.UUID     `json:"pvzId"`
	Status   ReceptionType `json:"status"`
}
