package models

import (
	"time"
)

type Listing struct {
	Id              int       `json:"id,omitempty"`
	ScannedAt       time.Time `json:"showtime"`
	PriceRange      string    `json:"priceRange"`
	PercentDiscount string    `json:"percentDiscount"`
	IsPlayOnly      bool      `json:"isPlayOnly"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	ShowId int   `json:"productionId"`
	Show   *Show `gorm:"foreignKey:ShowId" json:"show,omitempty"`
}

// type Tabler interface {
// 	TableName() string
// }

func (Listing) TableName() string {
	return "show_listings"
}
