package models

import (
	"time"
)

type Gross struct {
	Id             int       `json:"id,omitempty"`
	WeekEndDate    time.Time `json:"weekEndDate"`
	Name           string    `json:"name"`
	Theater        string    `json:"theater"`
	Gross          float32   `json:"gross"`
	Diff           float32   `json:"diff"`
	AvgTicket      float32   `json:"avgTicket"`
	TopTicket      float32   `json:"topTicket"`
	SeatsSold      int       `json:"seatsSold"`
	SeatsInTheater int       `json:"seatsInTheater"`
	Performances   int       `json:"performances"`
	Previews       int       `json:"previews"`
	PercentCap     float32   `json:"percentCap"`
	DiffPercentCap float32   `json:"diffPercentCap"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
