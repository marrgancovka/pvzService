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

type ReceptionRequest struct {
	PvzID uuid.UUID `json:"pvzId"`
}

type ReceptionWithProducts struct {
	Reception *Reception `json:"reception"`
	Products  []*Product `json:"products"`
}

func (receptionType ReceptionType) IsValid() bool {
	switch receptionType {
	case StatusInProgress, StatusClose:
		return true
	default:
		return false
	}
}
