package models

import "time"

type DeviceToken struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Token      string    `gorm:"unique;not null" json:"token"`
	DeviceInfo string    `gorm:"type:jsonb" json:"deviceInfo"` // JSON field
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
