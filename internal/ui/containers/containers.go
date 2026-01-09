// Package containers defines the containers list component
package containers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/givensuman/containertui/internal/context"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type sessionState int

const (
	viewMain sessionState = iota
	viewOverlay
)

type Model struct {
	sessionState sessionState
	width        int
	height       int
	foreground   tea.Model
	background   tea.Model
	overlayModel *overlay.Model
}

func New() Model {
	width, height := context.GetWindowSize()

	deleteConfirmation := NewDeleteConfirmation(nil)
	containerList := newContainerList()

	return Model{
		sessionState: viewMain,
		width:        width,
		height:       height,
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

func (m *Model) ToggleOverlay() {
	switch m.sessionState {
	case viewMain:
		m.sessionState = viewOverlay
	case viewOverlay:
		m.sessionState = viewMain
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.sessionState {
	case viewMain:
		bg, cmd := m.background.Update(msg)
		m.background = bg
		cmds = append(cmds, cmd)
	case viewOverlay:
		fg, cmd := m.foreground.Update(msg)
		m.foreground = fg
		cmds = append(cmds, cmd)
	}

	m.overlayModel = overlay.New(
		m.foreground,
		m.background,
		overlay.Center,
		overlay.Center,
		0,
		0,
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.sessionState == viewOverlay {
		return m.overlayModel.View()
	}

	return m.background.View()
}
