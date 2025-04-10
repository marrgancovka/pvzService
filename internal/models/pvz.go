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

type PVZ struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             City      `json:"city"`
}

func (citi *City) IsValid() bool {
	return *citi == CityMoscow || *citi == CitySpb || *citi == CityKazan
}
