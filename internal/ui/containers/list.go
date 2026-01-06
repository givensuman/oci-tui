package containers

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/context"
)

var style lipgloss.Style = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	list               list.Model
	selectedContainers map[string]int
	keybindings        *keybindings
}

var _ tea.Model = (*Model)(nil)

func NewContainersList() Model {
	containers := context.GetClient().GetContainers()
	var containerItems []list.Item
	for _, container := range containers {
		containerItems = append(
			containerItems,
			ContainerItem{
				Container:  container,
				isSelected: false,
			},
		)
	}

	width, height := context.GetWindowSize()
	list := list.New(containerItems, newDefaultDelegate(), width, height)

	keybindings := newKeybindings()
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keybindings.pauseContainer,
			keybindings.unpauseContainer,
			keybindings.startContainer,
			keybindings.stopContainer,
			keybindings.toggleSelect,
		}
	}

	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)

	selectedContainers := make(map[string]int)

	return Model{list, selectedContainers, keybindings}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		widthOffset, heightOffset := style.GetFrameSize()
		m.list.SetWidth(msg.Width - widthOffset)
		m.list.SetHeight(msg.Height - heightOffset)
	}

	list, cmd := m.list.Update(msg)
	m.list = list

	return m, cmd
}

func (m Model) View() string {
	return style.Render(m.list.View())
}
