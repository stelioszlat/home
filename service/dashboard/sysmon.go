package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type SysmonMessage struct {
	Cpu         []cpu.InfoStat
	Times       []cpu.TimesStat
	Percent     []float64
	Cores       int
	Vmem        mem.VirtualMemoryStat
	Swap        *mem.SwapMemoryStat
	SwapDevices []*mem.SwapDevice
	Partitions  []disk.PartitionStat
	DiskUsage   disk.UsageStat
}

type SysmonModel struct {
	CpuModel
	MemModel
	DiskModel
	spinner spinner.Model
}

type CpuModel struct {
	cpu     []cpu.InfoStat
	times   []cpu.TimesStat
	percent []float64
	cores   int
}

type MemModel struct {
	vmem        mem.VirtualMemoryStat
	swap        *mem.SwapMemoryStat
	swapDevices []*mem.SwapDevice
}

type DiskModel struct {
	partitions []disk.PartitionStat
	diskUsage  disk.UsageStat
}

func getCPUInfo() CpuModel {
	info, err := cpu.Info()
	percent, err := cpu.Percent(time.Second, true)
	counts, err := cpu.Counts(true)
	times, err := cpu.Times(true)
	if err != nil || (len(info) == 0 && len(times) == 0) {
		return CpuModel{}
	}
	return CpuModel{cpu: info, percent: percent, cores: counts, times: times}
}

func getMemInfo() MemModel {
	vmStat, err := mem.VirtualMemory()
	swap, err := mem.SwapMemory()
	swapDevices, err := mem.SwapDevices()
	if err != nil {
		return MemModel{}
	}
	return MemModel{vmem: *vmStat, swap: swap, swapDevices: swapDevices}
}

func getDiskInfo() DiskModel {
	partitions, err := disk.Partitions(false)
	if err != nil || len(partitions) == 0 {
		return DiskModel{}
	}

	// usage := disk.Usage(parti).Usage(partitions[0].Mountpoint)
	// if err != nil {
	return DiskModel{partitions: partitions}
	// }
	// return DiskModel{partitions: partitions, diskUsage: *usage}
}

func (s SysmonMessage) getSysmonMessage() SysmonMessage {
	cpuInfo := getCPUInfo()
	memInfo := getMemInfo()
	return SysmonMessage{
		Cpu:         cpuInfo.cpu,
		Percent:     cpuInfo.percent,
		Times:       cpuInfo.times,
		Cores:       cpuInfo.cores,
		Vmem:        memInfo.vmem,
		Swap:        memInfo.swap,
		SwapDevices: memInfo.swapDevices,
	}
}

func (s SysmonModel) Init() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return SysmonMessage{}.getSysmonMessage()
	})
}

func (s SysmonModel) Update(msg SysmonMessage) SysmonModel {
	s.cpu = msg.Cpu
	s.times = msg.Times
	s.percent = msg.Percent
	s.cores = msg.Cores
	s.vmem = msg.Vmem
	s.swap = msg.Swap
	s.swapDevices = msg.SwapDevices
	s.partitions = msg.Partitions
	s.diskUsage = msg.DiskUsage

	return s
}

func NewSysmon() SysmonModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return SysmonModel{spinner: s}
}

func sysmonBar(pct float64, width int) string {
	if width <= 0 {
		return ""
	}
	filled := int(pct / 100 * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

func getHeartbeat() string {
	// Blinks every second
	if time.Now().Second()%2 == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render("●")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("235")).Render("○")
}

func (s SysmonModel) View(width int) string {
	barWidth := width - 40
	var builder strings.Builder
	usage := s.CpuModel.percent
	fmt.Fprintf(&builder, "CPU %s\n", getHeartbeat())
	for i, cpu := range s.CpuModel.cpu {
		times := s.CpuModel.times
		fmt.Fprintf(&builder, "CPU %d - Core %s 	%s %.1f%%%s\n", i, cpu.CoreID, sysmonBar(usage[i], barWidth), usage[i], s.spinner.View())
		fmt.Fprintf(&builder, "User: %.1fm - System: %.1fm - Idle: %.1fm - Steal: %.1fm\n", times[i].User/60.0, times[i].System/60.0, times[i].Idle/60.0, times[i].Steal/60.0)
	}

	fmt.Fprintf(&builder, "\nMemory %s %.2f%% Available %.2f Mb of %.2f Mb\n", sysmonBar(float64(s.MemModel.vmem.UsedPercent), barWidth), s.MemModel.vmem.UsedPercent, float64(s.MemModel.vmem.Available)*0.000001, float64(s.MemModel.vmem.Total)*0.000001)

	// fmt.Fprintf(&builder, "\nDisk %s Free %dB Used %dB Total %dB\n", s.DiskModel.diskUsage.Path, s.DiskModel.diskUsage.Free, s.DiskModel.diskUsage.Used, s.DiskModel.diskUsage.Total)
	// for _, part := range s.DiskModel.partitions {
	// 	fmt.Fprintf(&builder, "Partition %s\n", part.Device)
	// }
	return builder.String()
}
