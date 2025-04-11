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

type ProductRequest struct {
	Type  ProductType `json:"type"`
	PvzID uuid.UUID   `json:"pvzId"`
}

func (productType *ProductType) IsValid() bool {
	return *productType == TypeElectronics || *productType == TypeClothes || *productType == TypeShoes
}
