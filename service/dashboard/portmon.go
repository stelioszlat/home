package dashboard

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type PortmonMessage struct {
	listenOnly bool
	Ports      []Port
}

type PortmonModel struct {
	listenOnly bool
	Ports      []Port
}

type Port struct {
	PortNr     uint32
	PortApp    string
	PortStatus string
}

func getPortInfo(listenOnly bool) PortmonMessage {
	connections, err := net.Connections("all")
	if err != nil {
		panic(err)
	}
	var ports []Port
	for _, connection := range connections {
		if listenOnly || connection.Status == "LISTEN" {
			p, err := process.NewProcess(connection.Pid)
			if err != nil {
				continue
			}
			cmdline, err := p.CmdlineSlice()
			if err != nil {
				continue
			}
			ports = append(ports, Port{PortNr: connection.Laddr.Port, PortApp: cmdline[len(cmdline)-1], PortStatus: connection.Status})
		}
	}

	return PortmonMessage{Ports: ports}
}

func (s PortmonMessage) getPortmonMessage() PortmonMessage {
	return getPortInfo(s.listenOnly)
}

func (s PortmonModel) ToggleListen() PortmonModel {
	s.listenOnly = !s.listenOnly
	return s
}

func (s PortmonModel) Init() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return PortmonMessage{listenOnly: s.listenOnly}.getPortmonMessage()
	})
}

func (s PortmonModel) Update(msg PortmonMessage) PortmonModel {
	s.Ports = msg.Ports
	s.listenOnly = msg.listenOnly
	return s
}

func NewPortmon() PortmonModel {
	return PortmonModel{}
}

func (s PortmonModel) View() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Ports Monitor\n")
	for i := range s.Ports {
		p := s.Ports[i]
		fmt.Fprintf(&builder, "%s %d    %s    \n", p.PortStatus, p.PortNr, p.PortApp)
	}
	return builder.String()
}
