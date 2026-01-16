// Package volumes defines the volumes component.
package volumes

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

// selectedVolumes is map of a volume's Name (volumes don't always have IDs like images/containers) to
// its index in the list. The client uses Name as the primary identifier for volumes usually.
// Wait, client.Volume struct doesn't have ID field, only Name.
type selectedVolumes struct {
	selections map[string]int
}

func newSelectedVolumes() *selectedVolumes {
	return &selectedVolumes{
		selections: make(map[string]int),
	}
}

func (sv *selectedVolumes) selectVolumeInList(name string, index int) {
	sv.selections[name] = index
}

func (sv selectedVolumes) unselectVolumeInList(name string) {
	delete(sv.selections, name)
}

type Model struct {
	shared.Component
	style           lipgloss.Style
	list            list.Model
	selectedVolumes *selectedVolumes
	keybindings     *keybindings
}

var (
	_ tea.Model             = (*Model)(nil)
	_ shared.ComponentModel = (*Model)(nil)
)

func New() Model {
	volumes, err := context.GetClient().GetVolumes()
	if err != nil {
		volumes = []client.Volume{}
	}
	items := make([]list.Item, 0, len(volumes))
	for _, vol := range volumes {
		items = append(items, VolumeItem{Volume: vol})
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
		style:           style,
		list:            l,
		selectedVolumes: newSelectedVolumes(),
		keybindings:     keybindings,
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
	selectedItem, ok := m.list.SelectedItem().(VolumeItem)
	if ok {
		isSelected := selectedItem.isSelected

		if isSelected {
			m.selectedVolumes.unselectVolumeInList(selectedItem.Volume.Name)
		} else {
			m.selectedVolumes.selectVolumeInList(selectedItem.Volume.Name, index)
		}

		selectedItem.isSelected = !isSelected
		m.list.SetItem(index, selectedItem)
	}
}

func (m *Model) handleToggleSelectionOfAll() {
	allAlreadySelected := true
	items := m.list.Items()

	for _, item := range items {
		if c, ok := item.(VolumeItem); ok {
			if _, selected := m.selectedVolumes.selections[c.Volume.Name]; !selected {
				allAlreadySelected = false
				break
			}
		}
	}

	if allAlreadySelected {
		// Unselect all items
		m.selectedVolumes = newSelectedVolumes()

		for index, item := range m.list.Items() {
			item, ok := item.(VolumeItem)
			if ok {
				item.isSelected = false
				m.list.SetItem(index, item)
			}
		}
	} else {
		// Select all items
		m.selectedVolumes = newSelectedVolumes()

		for index, item := range m.list.Items() {
			item, ok := item.(VolumeItem)
			if ok {
				item.isSelected = true
				m.list.SetItem(index, item)
				m.selectedVolumes.selectVolumeInList(item.Volume.Name, index)
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
