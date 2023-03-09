package models

import (
	"time"
)

type Show struct {
	Id        int       `json:"id,omitempty"`
	Showtime  time.Time `json:"showtime"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	ProductionId int         `json:"productionId"`
	Production   *Production `gorm:"foreignKey:ProductionId" json:"production,omitempty"`

	Listings []*Listing `gorm:"foreignKey:ShowId" json:"listings,omitempty"`
}

func (s *Show) AveragePrice() float64 {
	var total float64
	for _, listing := range s.Listings {
		priceRange := listing.ParsedPriceRange()
		total += priceRange.High
		total += priceRange.Low
	}
	average := total / float64(len(s.Listings)*2)
	return average
}
