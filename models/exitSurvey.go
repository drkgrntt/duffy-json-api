package models

import (
	"time"
)

type ExitSurvey struct {
	Price     string    `json:"price"`
	Expect    string    `json:"expect"`
	Kids      string    `json:"kids"`
	Choice    string    `json:"choice"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Timestamp time.Time `gorm:"column:primaryId" json:"timestamp"`
}

// type Tabler interface {
// 	TableName() string
// }

func (ExitSurvey) TableName() string {
	return "demo_exit"
}
