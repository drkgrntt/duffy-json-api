package models

import (
	"time"
)

type Gross struct {
	Id             int       `json:"id,omitempty"`
	WeekEndDate    time.Time `json:"weekEndDate"`
	Name           string    `json:"name"`
	Theater        string    `json:"theater"`
	Gross          float64   `json:"gross"`
	Diff           float64   `json:"diff"`
	AvgTicket      float64   `json:"avgTicket"`
	TopTicket      float64   `json:"topTicket"`
	SeatsSold      int       `json:"seatsSold"`
	SeatsInTheater int       `json:"seatsInTheater"`
	Performances   int       `json:"performances"`
	Previews       int       `json:"previews"`
	PercentCap     float64   `json:"percentCap"`
	DiffPercentCap float64   `json:"diffPercentCap"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
