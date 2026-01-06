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
	selectedContainers selectedContainers
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
			keybindings.toggleSelection,
			keybindings.toggleSelectionOfAll,
		}
	}

	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)

	selectedContainers := make(selectedContainers)

	return Model{list, selectedContainers, keybindings}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model := m

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		widthOffset, heightOffset := style.GetFrameSize()
		m.list.SetWidth(msg.Width - widthOffset)
		m.list.SetHeight(msg.Height - heightOffset)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keybindings.pauseContainer):
			model = m.handlePauseContainers()
		case key.Matches(msg, m.keybindings.unpauseContainer):
			model = m.handleUnpauseContainers()
		case key.Matches(msg, m.keybindings.startContainer):
			model = m.handleStartContainers()
		case key.Matches(msg, m.keybindings.stopContainer):
			model = m.handleStopContainers()
		case key.Matches(msg, m.keybindings.toggleSelection):
			index := m.list.Index()
			selectedItem, ok := m.list.SelectedItem().(ContainerItem)
			if ok {
				isSelected := selectedItem.isSelected

				if isSelected {
					m.selectedContainers = m.selectedContainers.unselectContainerInList(selectedItem.ID)
				} else {
					m.selectedContainers = m.selectedContainers.selectContainerInList(selectedItem.ID, index)
				}

				selectedItem.isSelected = !isSelected
				m.list.SetItem(index, selectedItem)
			}
		case key.Matches(msg, m.keybindings.toggleSelectionOfAll):
			allAlreadySelected := true
			items := m.list.Items()

			for _, item := range items {
				if c, ok := item.(ContainerItem); ok {
					if _, selected := m.selectedContainers[c.ID]; !selected {
						allAlreadySelected = false
						break
					}
				}
			}

			if allAlreadySelected {
				// Unselect all items
				model.selectedContainers = make(selectedContainers)

				for index, item := range m.list.Items() {
					item, ok := item.(ContainerItem)
					if ok {
						item.isSelected = false
						model.list.SetItem(index, item)
					}
				}
			} else {
				// Select all items
				model.selectedContainers = make(selectedContainers)

				for index, item := range m.list.Items() {
					item, ok := item.(ContainerItem)
					if ok {
						item.isSelected = true
						model.list.SetItem(index, item)
						model.selectedContainers = model.selectedContainers.selectContainerInList(item.ID, index)
					}
				}
			}
		}
	}

	list, cmd := m.list.Update(msg)
	model.list = list

	return model, cmd
}

func (m Model) View() string {
	return style.Render(m.list.View())
}
