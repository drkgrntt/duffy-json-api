package models

import (
	"time"
)

type TktsSalesDay struct {
	Id          int       `json:"id,omitempty"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	TicketsSold int       `json:"ticketsSold"`
	SalesTotal  float64   `json:"salesTotal"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
