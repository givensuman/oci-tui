// Package tabs implements the tab navigation component for the TUI.
package tabs

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type Tab int

const (
	Containers Tab = iota
	Images
	Volumes
	Networks
)

func (t Tab) String() string {
	return [...]string{
		"Containers",
		"Images",
		"Volumes",
		"Networks",
	}[t]
}

type KeyMap struct {
	SwitchToContainers key.Binding
	SwitchToImages     key.Binding
	SwitchToVolumes    key.Binding
	SwitchToNetworks   key.Binding
	NextTab            key.Binding
	PrevTab            key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		SwitchToContainers: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "containers"),
		),
		SwitchToImages: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "images"),
		),
		SwitchToVolumes: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "volumes"),
		),
		SwitchToNetworks: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "networks"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab", "right"),
			key.WithHelp("tab/right", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab", "left"),
			key.WithHelp("shift+tab/left", "prev tab"),
		),
	}
}

type Model struct {
	shared.Component
	ActiveTab Tab
	Tabs      []Tab
	KeyMap    KeyMap
}

func New() Model {
	return Model{
		ActiveTab: Containers,
		Tabs:      []Tab{Containers, Images, Volumes, Networks},
		KeyMap:    NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.NextTab):
			m.ActiveTab = (m.ActiveTab + 1) % Tab(len(m.Tabs))
		case key.Matches(msg, m.KeyMap.PrevTab):
			m.ActiveTab = (m.ActiveTab - 1 + Tab(len(m.Tabs))) % Tab(len(m.Tabs))
		case key.Matches(msg, m.KeyMap.SwitchToContainers):
			m.ActiveTab = Containers
		case key.Matches(msg, m.KeyMap.SwitchToImages):
			m.ActiveTab = Images
		case key.Matches(msg, m.KeyMap.SwitchToVolumes):
			m.ActiveTab = Volumes
		case key.Matches(msg, m.KeyMap.SwitchToNetworks):
			m.ActiveTab = Networks
		}
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	var tabs []string
	for _, t := range m.Tabs {
		if m.ActiveTab == t {
			tabs = append(tabs, activeTabStyle.Render(" "+t.String()+" "))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(" "+t.String()+" "))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Fill the rest of the line with the gap style
	gapWidth := max(0, m.WindowWidth-lipgloss.Width(row))
	gap := tabGapStyle.Render(strings.Repeat(" ", gapWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, row, gap)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	// Active tab: Primary color background, distinct text
	activeTabStyle = lipgloss.NewStyle().
			Foreground(colors.Black()).
			Background(colors.White()).
			Bold(true)

	// Inactive tab: Border/Gray background, Muted/Gray text (or White if background is dark enough)
	// Using Border color (Gray) for background to create a 'bar' effect.
	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(colors.Gray())

	// Gap: Same background as inactive tabs to create a continuous bar
	tabGapStyle = lipgloss.NewStyle()
)
