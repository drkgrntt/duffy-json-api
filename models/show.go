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
