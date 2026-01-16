// Package ui implements the terminal user interface
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"

	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/containers"
	"github.com/givensuman/containertui/internal/ui/images"
	"github.com/givensuman/containertui/internal/ui/networks"
	"github.com/givensuman/containertui/internal/ui/notifications"
	"github.com/givensuman/containertui/internal/ui/tabs"
	"github.com/givensuman/containertui/internal/ui/volumes"
)

// Model represents the top-level Bubbletea UI model.
// Contains terminal dimensions and the containers model (main TUI view).
type Model struct {
	width              int                 // current terminal width
	height             int                 // current terminal height
	tabsModel          tabs.Model          // tabs component
	containersModel    containers.Model    // main containers view/model
	imagesModel        images.Model        // images view/model
	volumesModel       volumes.Model       // volumes view/model
	networksModel      networks.Model      // networks view/model
	notificationsModel notifications.Model // notifications view/model
	overlayModel       *overlay.Model      // global overlay for notifications
}

func NewModel() Model {
	width, height := context.GetWindowSize()

	tabsM := tabs.New()
	conts := containers.New()
	imgs := images.New()
	vols := volumes.New()
	nets := networks.New()
	notifs := notifications.New()

	ov := overlay.New(notifs, conts, overlay.Right, overlay.Top, 0, 0)

	return Model{
		width:              width,
		height:             height,
		tabsModel:          tabsM,
		containersModel:    conts,
		imagesModel:        imgs,
		volumesModel:       vols,
		networksModel:      nets,
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

	// Message to pass to the overlay (and content)
	// By default it's the original message, but for resize events we adjust it.
	overlayMsg := msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update local size and global context
		m.width = msg.Width
		m.height = msg.Height
		context.SetWindowSize(msg.Width, msg.Height)

		// Update tabs with full width
		newTabs, _ := m.tabsModel.Update(msg)
		m.tabsModel = newTabs.(tabs.Model)

		// Calculate content height (subtract tabs height which is usually 1-2 lines)
		contentHeight := msg.Height - 2
		if contentHeight < 0 {
			contentHeight = 0
		}

		contentMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: contentHeight,
		}
		overlayMsg = contentMsg

		// Update sub-models with adjusted height
		newConts, _ := m.containersModel.Update(contentMsg)
		m.containersModel = newConts.(containers.Model)

		newImgs, _ := m.imagesModel.Update(contentMsg)
		m.imagesModel = newImgs.(images.Model)

		newVols, _ := m.volumesModel.Update(contentMsg)
		m.volumesModel = newVols.(volumes.Model)

		newNets, _ := m.networksModel.Update(contentMsg)
		m.networksModel = newNets.(networks.Model)

	case tea.KeyMsg:
		// Handle quit signals (Ctrl-C, Ctrl-D)
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}

		// Delegate tab switching to tabs model
		// This now handles tab/shift+tab and 1-4
		newTabs, cmd := m.tabsModel.Update(msg)
		m.tabsModel = newTabs.(tabs.Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update Notifications
	notifsModel, cmd := m.notificationsModel.Update(msg)
	m.notificationsModel = notifsModel.(notifications.Model)
	cmds = append(cmds, cmd)

	// Update active view based on selected tab
	var activeView tea.Model
	switch m.tabsModel.ActiveTab {
	case tabs.Containers:
		activeView = m.containersModel
		// Forward non-window/non-quit messages to containers view model
		// (WindowSizeMsg is handled specifically above to adjust for tab height)
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			containersModel, cmd := m.containersModel.Update(msg)
			m.containersModel = containersModel.(containers.Model)
			cmds = append(cmds, cmd)
			activeView = m.containersModel
		}
	case tabs.Images:
		activeView = m.imagesModel
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			imagesModel, cmd := m.imagesModel.Update(msg)
			m.imagesModel = imagesModel.(images.Model)
			cmds = append(cmds, cmd)
			activeView = m.imagesModel
		}
	case tabs.Volumes:
		activeView = m.volumesModel
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			volumesModel, cmd := m.volumesModel.Update(msg)
			m.volumesModel = volumesModel.(volumes.Model)
			cmds = append(cmds, cmd)
			activeView = m.volumesModel
		}
	case tabs.Networks:
		activeView = m.networksModel
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			networksModel, cmd := m.networksModel.Update(msg)
			m.networksModel = networksModel.(networks.Model)
			cmds = append(cmds, cmd)
			activeView = m.networksModel
		}
	}

	// Sync overlay
	m.overlayModel.Foreground = m.notificationsModel
	m.overlayModel.Background = activeView

	// Update overlay (mostly for resizing if needed, though we handle window size manually too)
	// We ignore the returned model here as we hold the reference in m.overlayModel pointer,
	// but the library might return a new one or update internal state.
	updatedOverlay, cmd := m.overlayModel.Update(overlayMsg)
	if ov, ok := updatedOverlay.(*overlay.Model); ok {
		m.overlayModel = ov
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the terminal as a string (delegated to containersModel).
func (m Model) View() string {
	// Render tabs at the top
	tabsView := m.tabsModel.View()

	// Render content below
	contentView := m.overlayModel.View()

	return lipgloss.JoinVertical(lipgloss.Top, tabsView, contentView)
}

// Start the Bubbletea UI main loop.
// Returns error if Bubbletea program terminates abnormally.
func Start() error {
	model := NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
