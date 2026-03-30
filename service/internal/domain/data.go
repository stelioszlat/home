package domain

import "time"

type Data struct {
	Location  string    `json:"location"`
	Activity  string    `json:"activity"`
	Device    string    `json:"device"`
	CreatedAt time.Time `json:"createdAt" gorm:"createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"updatedAt"`
}
