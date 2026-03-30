package notifications

import "time"

type Notification struct {
	Message   string
	App       string
	CreatedAt time.Time
}
