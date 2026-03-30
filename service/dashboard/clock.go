package dashboard

import (
	"fmt"
	"time"
)

type ClockModel struct {
	now time.Time
}

func NewClock() ClockModel {
	return ClockModel{now: time.Now()}
}

func (c ClockModel) Tick(t time.Time) ClockModel {
	c.now = t
	return c
}

func (c ClockModel) View() string {
	content := fmt.Sprintf("🕐 %s\n📅 %s",
		c.now.Format("15:04:05"),
		c.now.Format("Mon, Jan 2 2006"),
	)

	return content
}
