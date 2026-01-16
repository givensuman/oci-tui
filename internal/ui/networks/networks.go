// Package networks defines the networks component.
package networks

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type keybindings struct {
	toggleSelection      key.Binding
	toggleSelectionOfAll key.Binding
}

func newKeybindings() *keybindings {
	return &keybindings{
		toggleSelection: key.NewBinding(
			key.WithKeys(tea.KeySpace.String()),
			key.WithHelp("space", "toggle selection"),
		),
		toggleSelectionOfAll: key.NewBinding(
			key.WithKeys(tea.KeyCtrlA.String()),
			key.WithHelp("ctrl+a", "toggle selection of all"),
		),
	}
}

// selectedNetworks is map of a network's ID to
// its index in the list
type selectedNetworks struct {
	selections map[string]int
}

func newSelectedNetworks() *selectedNetworks {
	return &selectedNetworks{
		selections: make(map[string]int),
	}
}

func (sn *selectedNetworks) selectNetworkInList(id string, index int) {
	sn.selections[id] = index
}

func (sn selectedNetworks) unselectNetworkInList(id string) {
	delete(sn.selections, id)
}

type Model struct {
	shared.Component
	style            lipgloss.Style
	list             list.Model
	selectedNetworks *selectedNetworks
	keybindings      *keybindings
}

var (
	_ tea.Model             = (*Model)(nil)
	_ shared.ComponentModel = (*Model)(nil)
)

func New() Model {
	networks, err := context.GetClient().GetNetworks()
	if err != nil {
		networks = []client.Network{}
	}
	items := make([]list.Item, 0, len(networks))
	for _, net := range networks {
		items = append(items, NetworkItem{Network: net})
	}

	width, height := context.GetWindowSize()
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		PaddingTop(1)

	delegate := newDefaultDelegate()
	l := list.New(items, delegate, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	keybindings := newKeybindings()
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keybindings.toggleSelection,
			keybindings.toggleSelectionOfAll,
		}
	}

	return Model{
		style:            style,
		list:             l,
		selectedNetworks: newSelectedNetworks(),
		keybindings:      keybindings,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateWindowDimensions(msg)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keybindings.toggleSelection):
			m.handleToggleSelection()
		case key.Matches(msg, m.keybindings.toggleSelectionOfAll):
			m.handleToggleSelectionOfAll()
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) handleToggleSelection() {
	index := m.list.Index()
	selectedItem, ok := m.list.SelectedItem().(NetworkItem)
	if ok {
		isSelected := selectedItem.isSelected

		if isSelected {
			m.selectedNetworks.unselectNetworkInList(selectedItem.Network.ID)
		} else {
			m.selectedNetworks.selectNetworkInList(selectedItem.Network.ID, index)
		}

		selectedItem.isSelected = !isSelected
		m.list.SetItem(index, selectedItem)
	}
}

func (m *Model) handleToggleSelectionOfAll() {
	allAlreadySelected := true
	items := m.list.Items()

	for _, item := range items {
		if c, ok := item.(NetworkItem); ok {
			if _, selected := m.selectedNetworks.selections[c.Network.ID]; !selected {
				allAlreadySelected = false
				break
			}
		}
	}

	if allAlreadySelected {
		// Unselect all items
		m.selectedNetworks = newSelectedNetworks()

		for index, item := range m.list.Items() {
			item, ok := item.(NetworkItem)
			if ok {
				item.isSelected = false
				m.list.SetItem(index, item)
			}
		}
	} else {
		// Select all items
		m.selectedNetworks = newSelectedNetworks()

		for index, item := range m.list.Items() {
			item, ok := item.(NetworkItem)
			if ok {
				item.isSelected = true
				m.list.SetItem(index, item)
				m.selectedNetworks.selectNetworkInList(item.Network.ID, index)
			}
		}
	}
}

func (m Model) View() string {
	return m.style.Render(m.list.View())
}

func (m *Model) UpdateWindowDimensions(msg tea.WindowSizeMsg) {
	m.WindowWidth = msg.Width
	m.WindowHeight = msg.Height

	lm := shared.NewLayoutManager(msg.Width, msg.Height)
	dims := lm.CalculateFullscreen(m.style)

	m.style = m.style.Width(dims.Width).Height(dims.Height)
	m.list.SetWidth(dims.ContentWidth)
	m.list.SetHeight(dims.ContentHeight)
}
