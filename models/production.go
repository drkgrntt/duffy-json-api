package models

import (
	"time"
)

type Production struct {
	// ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ID                int       `json:"id,omitempty"`
	Name              string    `gorm:"type:varchar(255);not null" json:"name"`
	LastShownAt       time.Time `json:"lastShownAt"`
	LastScannedAt     time.Time `json:"lastScannedAt"`
	HasTktsData       bool      `json:"hasTktsData"`
	CollectionGroupId int       `json:"collectionGroupId"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
