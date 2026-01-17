// Package containers defines the containers component.
package containers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/ui/shared"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type sessionState int

const (
	viewMain sessionState = iota
	viewOverlay
)

// Model is the main containers screen for the TUI.
// Manages two views: either normal container list or overlay (logs, delete dialog).
type Model struct {
	shared.Component
	// sessionState governs whether we're in main or overlay view
	sessionState sessionState
	// foreground is the currently active overlay model (logs, confirmation, etc)
	foreground tea.Model
	// background is the container list model
	background tea.Model
	// overlayModel manages overlay transitions and rendering
	overlayModel *overlay.Model
}

var (
	_ tea.Model             = (*Model)(nil)
	_ shared.ComponentModel = (*Model)(nil)
)

func New() Model {
	deleteConfirmation := newDeleteConfirmation(nil)
	containerList := newContainerList()

	return Model{
		sessionState: viewMain,
		foreground:   deleteConfirmation,
		background:   containerList,
		overlayModel: overlay.New(
			deleteConfirmation,
			containerList,
			overlay.Center,
			overlay.Center,
			0,
			0,
		),
	}
}

// UpdateWindowDimensions updates all sub-models with latest terminal size.
// Called during resize events or window size changes.
func (m *Model) UpdateWindowDimensions(msg tea.WindowSizeMsg) {
	m.WindowWidth = msg.Width
	m.WindowHeight = msg.Height

	switch m.sessionState {
	case viewMain:
		if cl, ok := m.background.(ContainerList); ok {
			cl.UpdateWindowDimensions(msg)
			m.background = cl
		}
	case viewOverlay:
		// Forward dimension updates to correct overlay model
		switch fg := m.foreground.(type) {
		case *ContainerLogs:
			fg.setDimensions(msg.Width, msg.Height)
			m.foreground = fg
		case DeleteConfirmation:
			fg.UpdateWindowDimensions(msg)
			m.foreground = fg
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all Bubbletea messages relevant to this containers screen.
// Manages both main view logic and overlay/dialog/confirmation sub-models.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Update foreground/background model based on session mode.
	switch m.sessionState {
	case viewMain:
		bg, cmd := m.background.Update(msg)
		m.background = bg.(ContainerList)
		cmds = append(cmds, cmd)
	case viewOverlay:
		fg, cmd := m.foreground.Update(msg)
		m.foreground = fg
		cmds = append(cmds, cmd)
	}

	// Handle special message types (resize, open/close dialogs, etc)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateWindowDimensions(msg)

	case MessageCloseOverlay:
		m.sessionState = viewMain

	case MessageOpenDeleteConfirmationDialog:
		m.foreground = newDeleteConfirmation(msg.requestedContainersToDelete...)
		m.sessionState = viewOverlay
		cmds = append(cmds, m.foreground.Init())

		// case MessageOpenContainerLogs:
		// 	m.foreground = newContainerLogs(msg.container)
		// 	m.sessionState = viewOverlay
		// 	cmds = append(cmds, m.foreground.Init())
	}

	// Always update the overlay model and sync submodels
	m.overlayModel.Foreground = m.foreground
	m.overlayModel.Background = m.background

	updatedOverlayModel, cmd := m.overlayModel.Update(msg)
	overlayModel, ok := updatedOverlayModel.(*overlay.Model)
	if ok {
		m.overlayModel = overlayModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI for containers.
// If an overlay is open, delegates to overlay model’s view,
// Otherwise renders the main container background view.
func (m Model) View() string {
	if m.sessionState == viewOverlay && m.foreground != nil {
		return m.overlayModel.View()
	}

	// Get dimensions for master (list) and detail (inspect)
	lm := shared.NewLayoutManager(m.WindowWidth, m.WindowHeight)
	_, detail := lm.CalculateMasterDetail(lipgloss.NewStyle())

	// Render the list (background)
	listView := m.background.View()

	// Ensure the list view is constrained to its allocated width
	// If the list view is shorter than the allocated space, we might need to pad it?
	// But lipgloss usually handles this. The important part is that the background.View()
	// returns a string that respects the updated dimensions.
	// Since we called UpdateWindowDimensions on the background model, it should be fine.

	// Render the detail view (side pane)
	// For now, we'll just render a placeholder or simple text
	// TODO: Implement actual inspect view
	detailStyle := lipgloss.NewStyle().
		Width(detail.Width - 2). // Subtract 2 for border width compensation if needed, though Box sizing usually handles it. Let's rely on content width.
		Height(detail.Height).   // Ensure height matches master list
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Muted()).
		Padding(1)

	var detailContent string
	if cl, ok := m.background.(ContainerList); ok {
		item := cl.list.SelectedItem()
		if item != nil {
			if c, ok := item.(ContainerItem); ok {
				detailContent = fmt.Sprintf(
					"ID: %s\nName: %s\nImage: %s\nState: %s",
					c.ID, c.Name, c.Image, c.State,
				)
			}
		}
	}

	detailView := detailStyle.Render(detailContent)

	// Ensure heights match exactly before joining
	// The list view might be shorter if empty or filtered, so we rely on JoinHorizontal's Top alignment
	// but using Place to ensure exact box sizes is safer.

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
}
