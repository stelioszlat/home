package dashboard

import (
	"service/lights"

	"github.com/charmbracelet/lipgloss"
)

type LightsModel struct {
	on bool
}

func NewLights() LightsModel {
	return LightsModel{}
}

func (l LightsModel) Toggle() LightsModel {
	l.on = !l.on
	return l
}

func (l LightsModel) View() string {
	status := "OFF 🌑"
	color := "240"
	if l.on {
		status = "ON  💡"
		color = "226"
	}

	content := lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Render("Living Room Lights\n\n" + status + "\n\n[l] to toggle")

	lights.ToggleLights()

	return content
}
