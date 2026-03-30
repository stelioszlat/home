package domain

import "time"

type Module struct {
	ID          string    `json:"_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	IsEnabled   bool      `json:"-" gorm:"default:false"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
