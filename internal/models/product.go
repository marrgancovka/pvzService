package models

import (
	"github.com/google/uuid"
	"time"
)

type ProductType string

const (
	TypeElectronics ProductType = "электроника"
	TypeClothes     ProductType = "одежда"
	TypeShoes       ProductType = "обувь"
)

type Product struct {
	ID          uuid.UUID   `json:"id"`
	DateTime    time.Time   `json:"dateTime"`
	Type        ProductType `json:"type"`
	ReceptionID uuid.UUID   `json:"receptionId"`
}
