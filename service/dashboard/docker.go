package dashboard

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/docker"
)

type DockerMessage struct {
	CPU docker.CgroupCPUStat
}

type DockerModel struct {
	listenOnly bool
	Ports      []Port
}

func getDockerInfo() DockerMessage {
	cpuStat, err := docker.GetDockerStat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cpuStat)
	return DockerMessage{}
}

func (d DockerModel) getDockerMessage() DockerMessage {
	return getDockerInfo()
}

func (d DockerModel) Init() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return DockerModel{}.getDockerMessage()
	})
}

func (d DockerModel) Update(msg DockerMessage) DockerModel {
	return d
}

func NewDockermon() DockerModel {
	return DockerModel{}
}

func (s DockerModel) View() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Docker Monitor\n")

	return builder.String()
}
