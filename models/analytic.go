package models

import (
	"database/sql"
	"time"
)

type Analytic struct {
	Id        int            `json:"id,omitempty"`
	Page      string         `json:"page"`
	Query     string         `json:"query"`
	UserAgent string         `json:"useragent"`
	Country   string         `json:"country"`
	City      string         `json:"city"`
	County    sql.NullString `json:"county"`
	State     string         `json:"state"`
	Ip        string         `json:"ip"`
	Domain    string         `json:"domain"`
	UserId    sql.NullInt64  `json:"userId"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}
