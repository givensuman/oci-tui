// Package ui implements the terminal user interface
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	overlay "github.com/rmhubbert/bubbletea-overlay"

	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/containers"
	"github.com/givensuman/containertui/internal/ui/notifications"
)

// Model represents the top-level Bubbletea UI model.
// Contains terminal dimensions and the containers model (main TUI view).
type Model struct {
	width              int                 // current terminal width
	height             int                 // current terminal height
	containersModel    containers.Model    // main containers view/model
	notificationsModel notifications.Model // notifications view/model
	overlayModel       *overlay.Model      // global overlay for notifications
}

func NewModel() Model {
	width, height := context.GetWindowSize()

	conts := containers.New()
	notifs := notifications.New()

	ov := overlay.New(notifs, conts, overlay.Right, overlay.Top, 0, 0)

	return Model{
		width:              width,
		height:             height,
		containersModel:    conts,
		notificationsModel: notifs,
		overlayModel:       ov,
	}
}

// Init performs any initial commands for the Bubbletea app
// (no async initialization needed here)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all Bubbletea update messages dispatched by the tea runtime.
// Manages window resizing and quit keys, delegates other updates to containersModel.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update local size and global context
		m.width = msg.Width
		m.height = msg.Height
		context.SetWindowSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Handle quit signals (Ctrl-C, Ctrl-D)
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	// Update Notifications
	notifsModel, cmd := m.notificationsModel.Update(msg)
	m.notificationsModel = notifsModel.(notifications.Model)
	cmds = append(cmds, cmd)

	// Forward non-window/non-quit messages to containers view model
	containersModel, cmd := m.containersModel.Update(msg)
	m.containersModel = containersModel.(containers.Model)
	cmds = append(cmds, cmd)

	// Sync overlay
	m.overlayModel.Foreground = m.notificationsModel
	m.overlayModel.Background = m.containersModel

	// Update overlay (mostly for resizing if needed, though we handle window size manually too)
	// We ignore the returned model here as we hold the reference in m.overlayModel pointer,
	// but the library might return a new one or update internal state.
	updatedOverlay, cmd := m.overlayModel.Update(msg)
	if ov, ok := updatedOverlay.(*overlay.Model); ok {
		m.overlayModel = ov
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the terminal as a string (delegated to containersModel).
func (m Model) View() string {
	return m.overlayModel.View()
}

// Start the Bubbletea UI main loop.
// Returns error if Bubbletea program terminates abnormally.
func Start() error {
	model := NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
