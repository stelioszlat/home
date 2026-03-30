package models

import "time"

type NotificationLog struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	Title  string    `json:"title"`
	Body   string    `json:"body"`
	Data   string    `gorm:"type:jsonb" json:"data"` // JSON field
	SentAt time.Time `json:"sentAt"`
	Status string    `json:"status"` // "success", "failed"
}
