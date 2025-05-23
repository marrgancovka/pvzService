package models

import (
	"github.com/google/uuid"
	"time"
)

type City string

const (
	CityMoscow City = "Москва"
	CitySpb    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

type Pvz struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             City      `json:"city"`
}

type PvzWithReceptions struct {
	ID               uuid.UUID                `json:"id"`
	RegistrationDate time.Time                `json:"registrationDate"`
	City             City                     `json:"city"`
	Receptions       []*ReceptionWithProducts `json:"receptions"`
}

func (city City) IsValid() bool {
	switch city {
	case CitySpb, CityMoscow, CityKazan:
		return true
	default:
		return false
	}
}
