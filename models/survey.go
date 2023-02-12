package models

import (
	"time"
)

type Survey struct {
	Date         time.Time `gorm:"column:demo_date" json:"date"`
	Age          string    `gorm:"column:demo_age" json:"age"`
	Residence    string    `gorm:"column:demo_reside" json:"residence"`
	ShowsPerYear string    `gorm:"column:demo_shows_per_year" json:"showsPerYear"`
	ShopTkts     string    `gorm:"column:demo_shop_tkts" json:"shopTkts"`
	Noob         string    `gorm:"column:noob" json:"noob"`
	Timestamp    time.Time `gorm:"column:primaryId" json:"timestamp"`
}

type Tabler interface {
	TableName() string
}

func (Survey) TableName() string {
	return "demo_survey"
}
