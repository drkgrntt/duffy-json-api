package models

import (
	"time"
)

type Production struct {
	Id            int       `json:"id,omitempty"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	LastShownAt   time.Time `json:"lastShownAt"`
	LastScannedAt time.Time `json:"lastScannedAt"`
	HasTktsData   bool      `json:"hasTktsData"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	CompetitionGroupId int               `json:"competitionGroupId"`
	CompetitionGroup   *CompetitionGroup `gorm:"foreignKey:CompetitionGroupId" json:"competitionGroup,omitempty"`

	Shows []*Show `gorm:"foreignKey:ProductionId" json:"shows,omitempty"`
}
