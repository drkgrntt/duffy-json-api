package models

import (
	"time"
)

type CompetitionGroup struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Productions []Production `gorm:"foreignKey:CompetitionGroupId" json:"productions"`
}
