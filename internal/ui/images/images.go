// Package images defines the images component.
package images

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

// selectedImages is map of an image's ID to
// its index in the list
type selectedImages struct {
	selections map[string]int
}

func newSelectedImages() *selectedImages {
	return &selectedImages{
		selections: make(map[string]int),
	}
}

func (si *selectedImages) selectImageInList(id string, index int) {
	si.selections[id] = index
}

func (si selectedImages) unselectImageInList(id string) {
	delete(si.selections, id)
}

type Model struct {
	shared.Component
	style          lipgloss.Style
	list           list.Model
	selectedImages *selectedImages
	keybindings    *keybindings
}

var (
	_ tea.Model             = (*Model)(nil)
	_ shared.ComponentModel = (*Model)(nil)
)

func New() Model {
	images, err := context.GetClient().GetImages()
	if err != nil {
		images = []client.Image{}
	}
	items := make([]list.Item, 0, len(images))
	for _, img := range images {
		items = append(items, ImageItem{Image: img})
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
		style:          style,
		list:           l,
		selectedImages: newSelectedImages(),
		keybindings:    keybindings,
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
	selectedItem, ok := m.list.SelectedItem().(ImageItem)
	if ok {
		isSelected := selectedItem.isSelected

		if isSelected {
			m.selectedImages.unselectImageInList(selectedItem.Image.ID)
		} else {
			m.selectedImages.selectImageInList(selectedItem.Image.ID, index)
		}

		selectedItem.isSelected = !isSelected
		m.list.SetItem(index, selectedItem)
	}
}

func (m *Model) handleToggleSelectionOfAll() {
	allAlreadySelected := true
	items := m.list.Items()

	for _, item := range items {
		if c, ok := item.(ImageItem); ok {
			if _, selected := m.selectedImages.selections[c.Image.ID]; !selected {
				allAlreadySelected = false
				break
			}
		}
	}

	if allAlreadySelected {
		// Unselect all items
		m.selectedImages = newSelectedImages()

		for index, item := range m.list.Items() {
			item, ok := item.(ImageItem)
			if ok {
				item.isSelected = false
				m.list.SetItem(index, item)
			}
		}
	} else {
		// Select all items
		m.selectedImages = newSelectedImages()

		for index, item := range m.list.Items() {
			item, ok := item.(ImageItem)
			if ok {
				item.isSelected = true
				m.list.SetItem(index, item)
				m.selectedImages.selectImageInList(item.Image.ID, index)
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
