package dashboard

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	clock     ClockModel
	sysmon    SysmonModel
	portmon   PortmonModel
	dockermon DockerModel
	lights    LightsModel
	width     int
	height    int
}

func initialModel() model {
	return model{
		clock:     NewClock(),
		sysmon:    NewSysmon(),
		portmon:   NewPortmon(),
		dockermon: NewDockermon(),
		lights:    NewLights(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), m.sysmon.Init(), m.portmon.Init(), m.dockermon.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "l":
			m.lights = m.lights.Toggle()
		case "p":
			m.portmon = m.portmon.ToggleListen()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.clock = m.clock.Tick(time.Time(msg))
		cmds = append(cmds, tick())

	case SysmonMessage:
		var spCmd tea.Cmd
		m.sysmon = m.sysmon.Update(msg)
		m.sysmon.spinner, spCmd = m.sysmon.spinner.Update(msg)
		cmds = append(cmds, m.sysmon.Init())
		cmds = append(cmds, spCmd)
	case PortmonMessage:
		m.portmon = m.portmon.Update(msg)
		cmds = append(cmds, m.portmon.Init())
	case DockerMessage:
		m.dockermon = m.dockermon.Update(msg)
		cmds = append(cmds, m.dockermon.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	colWidth := m.width/2 - 2

	bgStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")). // Tokyo Night Dark
		Foreground(lipgloss.Color("#a9b1d6")).
		Width(m.width).
		Height(m.height)

	left := lipgloss.JoinVertical(lipgloss.Left,
		m.renderBlock("Clock", m.clock.View(), colWidth, "c"),
		m.renderBlock("Lights", m.lights.View(), colWidth, "l"),
		m.renderBlock("Ports", m.portmon.View(), colWidth, "p"),
	)
	right := lipgloss.JoinVertical(lipgloss.Right,
		m.renderBlock("System", m.sysmon.View(colWidth), colWidth, "s"),
		m.renderBlock("Docker", m.dockermon.View(), colWidth, "d"),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	help := m.renderBlock("Help", "  [q] quit  [l] toggle lights  [p] toggle listen ports", colWidth, "")

	dashboardLayout := lipgloss.JoinVertical(lipgloss.Left, body, help)

	return bgStyle.Render(dashboardLayout)
}

func (m model) renderBlock(title string, content string, width int, hotkey string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00ee00")).
		// Padding(0, 1).
		Width(width).
		Render(content)
}

func Dashboard() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
