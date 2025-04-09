package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}

type TokenPayload struct {
	ID   uuid.UUID
	Role Role
	Exp  time.Time
}
