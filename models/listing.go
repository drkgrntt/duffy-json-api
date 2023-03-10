package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Listing struct {
	Id              int       `json:"id,omitempty"`
	ScannedAt       time.Time `json:"showtime"`
	PriceRange      string    `json:"priceRange"`
	PercentDiscount string    `json:"percentDiscount"`
	IsPlayOnly      bool      `json:"isPlayOnly"`
	Broadway        bool      `json:"broadway"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	ShowId int   `json:"productionId"`
	Show   *Show `gorm:"foreignKey:ShowId" json:"show,omitempty"`
}

func (Listing) TableName() string {
	return "show_listings"
}

type PriceRange struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

func (l *Listing) ParsedPriceRange() PriceRange {
	pieces := strings.Split(l.PriceRange, "$")
	prices := strings.Split(pieces[len(pieces)-1], "-")
	var low float64
	var high float64
	var err error
	if low, err = strconv.ParseFloat(prices[0], 64); err != nil {
		fmt.Println("Error parsing low price", prices[0], "to an int")
		fmt.Println(prices)
	}
	if high, err = strconv.ParseFloat(prices[len(prices)-1], 64); err != nil {
		fmt.Println("Error parsing high price", prices[len(prices)-1], "to an int")
		fmt.Println(prices)
	}
	priceRange := PriceRange{low, high}
	return priceRange
}

func (l *Listing) IsWithinParams(
	minPrice float64,
	maxPrice float64,
	minDiscount float64,
	maxDiscount float64,
) bool {
	priceRange := l.ParsedPriceRange()
	if float64(priceRange.High) > maxPrice {
		return false
	}

	if float64(priceRange.Low) < minPrice {
		return false
	}

	discount, _ := strconv.ParseFloat(l.PercentDiscount, 64)
	if discount > maxDiscount {
		return false
	}

	if discount < minDiscount {
		return false
	}

	return true
}
