// Package images defines the images component.
package images

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type detailsKeybindings struct {
	Up     key.Binding
	Down   key.Binding
	Switch key.Binding
}

func newDetailsKeybindings() detailsKeybindings {
	return detailsKeybindings{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Switch: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
	}
}

type keybindings struct {
	toggleSelection      key.Binding
	toggleSelectionOfAll key.Binding
	remove               key.Binding
	switchTab            key.Binding
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
		remove: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "remove"),
		),
		switchTab: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "tab", "shift+tab"),
			key.WithHelp("1-4/tab", "switch tab"),
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

type sessionState int

const (
	viewMain sessionState = iota
	viewOverlay
)

const (
	focusList = iota
	focusDetails
)

type Model struct {
	shared.Component
	style          lipgloss.Style
	list           list.Model
	viewport       viewport.Model
	selectedImages *selectedImages
	keybindings    *keybindings

	// Overlay support
	sessionState sessionState
	// focusedView governs which panel is active (0: List, 1: Details)
	focusedView        int
	detailsKeybindings detailsKeybindings
	foreground         tea.Model
	overlayModel       *overlay.Model
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
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.Primary())
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.Primary())
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(colors.Primary())
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(colors.Primary())

	keybindings := newKeybindings()
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keybindings.toggleSelection,
			keybindings.toggleSelectionOfAll,
			keybindings.remove,
			keybindings.switchTab,
		}
	}

	vp := viewport.New(0, 0)

	m := Model{
		style:              style,
		list:               l,
		viewport:           vp,
		selectedImages:     newSelectedImages(),
		keybindings:        keybindings,
		sessionState:       viewMain,
		focusedView:        focusList,
		detailsKeybindings: newDetailsKeybindings(),
	}

	// Initialize overlay with nil content initially
	m.overlayModel = overlay.New(nil, m.list, overlay.Center, overlay.Center, 0, 0)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.sessionState {
	case viewOverlay:
		// If overlay is active, route messages there
		fg, cmd := m.foreground.Update(msg)
		m.foreground = fg
		cmds = append(cmds, cmd)

		// Check if overlay wants to close or confirm
		if _, ok := msg.(shared.CloseDialogMessage); ok {
			m.sessionState = viewMain
			m.foreground = nil
		} else if confirm, ok := msg.(shared.ConfirmationMessage); ok {
			// Handle confirmation
			if confirm.Action.Type == "DeleteImage" {
				id := confirm.Action.Payload.(string)
				// Perform deletion
				err := context.GetClient().RemoveImage(id)
				if err != nil {
					// TODO: Show error
				} else {
					// Remove from list
					// For now, just refresh (simple but ineffective for large lists)
					// m.Refresh() // Need to implement refresh or manually remove
				}
			}
			m.sessionState = viewMain
			m.foreground = nil
		}
	case viewMain:
		// Main loop handling

		// Handle focus switching
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "tab" && m.list.FilterState() != list.Filtering {
				if m.focusedView == focusList {
					m.focusedView = focusDetails
				} else {
					m.focusedView = focusList
				}
				return m, nil
			}
		}

		isKeyMsg := false
		if _, ok := msg.(tea.KeyMsg); ok {
			isKeyMsg = true
		}

		// Update list if focused or not a key message
		if !isKeyMsg || m.focusedView == focusList {
			switch msg := msg.(type) {
			case tea.WindowSizeMsg:
				m.UpdateWindowDimensions(msg)
			case tea.KeyMsg:
				if m.list.FilterState() == list.Filtering {
					break
				}

				switch {
				// allow passing through of tab switching keys
				case key.Matches(msg, m.keybindings.switchTab):
					return m, nil
				case key.Matches(msg, m.keybindings.toggleSelection):
					m.handleToggleSelection()
				case key.Matches(msg, m.keybindings.toggleSelectionOfAll):
					m.handleToggleSelectionOfAll()
				case key.Matches(msg, m.keybindings.remove):
					// Trigger remove dialog
					item := m.list.SelectedItem()
					if item != nil {
						if i, ok := item.(ImageItem); ok {
							// Check usage
							usedBy, _ := context.GetClient().GetContainersUsingImage(i.Image.ID)
							if len(usedBy) > 0 {
								// Show warning dialog with navigation option
								// For now just show warning
								dialog := shared.NewSmartDialog(
									fmt.Sprintf("Image %s is used by %d containers (%v).\nCannot delete.", i.Image.ID[:12], len(usedBy), usedBy),
									[]shared.DialogButton{
										{Label: "OK", IsSafe: true},
									},
								)
								m.foreground = dialog
								m.sessionState = viewOverlay
							} else {
								// Show confirmation
								dialog := shared.NewSmartDialog(
									fmt.Sprintf("Are you sure you want to delete image %s?", i.Image.ID[:12]),
									[]shared.DialogButton{
										{Label: "Cancel", IsSafe: true},
										{Label: "Delete", IsSafe: false, Action: shared.SmartDialogAction{Type: "DeleteImage", Payload: i.Image.ID}},
									},
								)
								m.foreground = dialog
								m.sessionState = viewOverlay
							}
						}
					}
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Check for selection change (if list updated)
		item := m.list.SelectedItem()
		if item != nil {
			if i, ok := item.(ImageItem); ok {
				content := fmt.Sprintf(
					"ID: %s\nSize: %d\nTags: %v",
					i.Image.ID, i.Image.Size, i.Image.RepoTags,
				)
				m.viewport.SetContent(content)
			}
		}

		// Update viewport if focused or not a key message
		if !isKeyMsg || m.focusedView == focusDetails {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Always update overlay model to sync state
	m.overlayModel.Foreground = m.foreground
	// m.overlayModel.Background = m.list // list is not a tea.Model, it's inside m which is
	// We need to render the background manually or pass a wrapper.
	// Actually overlayModel.Update handles internal state.
	// The View() calls overlayModel.View().
	// overlay.New takes (foreground, background).
	// We initialized it with m.list, but m.list changes.
	m.overlayModel.Background = m.list

	return m, tea.Batch(cmds...)
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
	if m.sessionState == viewOverlay && m.foreground != nil {
		return m.overlayModel.View()
	}

	// Get dimensions for master (list) and detail (inspect)
	lm := shared.NewLayoutManager(m.WindowWidth, m.WindowHeight)
	_, detail := lm.CalculateMasterDetail(lipgloss.NewStyle())

	// Render the list (background)
	listView := m.style.Render(m.list.View())

	borderColor := colors.Muted()
	if m.focusedView == focusDetails {
		borderColor = colors.Primary()
	}

	// Render the detail view (side pane)
	detailStyle := lipgloss.NewStyle().
		Width(detail.Width - 2). // Subtract 2 for border width compensation
		Height(detail.Height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1)

	var detailContent string
	if m.list.SelectedItem() != nil {
		detailContent = m.viewport.View()
	} else {
		detailContent = lipgloss.NewStyle().Foreground(colors.Muted()).Render("No image selected")
	}

	detailView := detailStyle.Render(detailContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
}

func (m *Model) UpdateWindowDimensions(msg tea.WindowSizeMsg) {
	m.WindowWidth = msg.Width
	m.WindowHeight = msg.Height

	lm := shared.NewLayoutManager(msg.Width, msg.Height)
	master, _ := lm.CalculateMasterDetail(m.style)

	m.style = m.style.Width(master.Width).Height(master.Height) // Update style dimensions

	// Update viewport size
	_, detail := lm.CalculateMasterDetail(lipgloss.NewStyle())
	width := detail.Width - 4
	height := detail.Height - 2
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	m.viewport.Width = width
	m.viewport.Height = height

	switch m.sessionState {
	case viewMain:
		if m.list.Width() != master.ContentWidth || m.list.Height() != master.ContentHeight {
			m.list.SetWidth(master.ContentWidth)
			m.list.SetHeight(master.ContentHeight)
		}
	case viewOverlay:
		if d, ok := m.foreground.(shared.SmartDialog); ok {
			d.UpdateWindowDimensions(msg)
			m.foreground = d
		}
	}
}

func (m Model) ShortHelp() []key.Binding {
	if m.focusedView == focusList {
		return m.list.ShortHelp()
	} else if m.focusedView == focusDetails {
		return []key.Binding{
			m.detailsKeybindings.Up,
			m.detailsKeybindings.Down,
			m.detailsKeybindings.Switch,
		}
	}
	return nil
}

func (m Model) FullHelp() [][]key.Binding {
	if m.focusedView == focusList {
		return m.list.FullHelp()
	} else if m.focusedView == focusDetails {
		return [][]key.Binding{
			{
				m.detailsKeybindings.Up,
				m.detailsKeybindings.Down,
				m.detailsKeybindings.Switch,
			},
		}
	}
	return nil
}
