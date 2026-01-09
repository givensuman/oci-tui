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
	background   ContainerList
	overlayModel *overlay.Model
}

func New() Model {
	width, height := context.GetWindowSize()

	deleteConfirmation := newDeleteConfirmation(nil)
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case MessageCloseOverlay:
		m.sessionState = viewMain
	case MessageOpenDeleteConfirmationDialog:
		m.foreground = newDeleteConfirmation(msg.item)
		m.sessionState = viewOverlay
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
	if m.sessionState == viewOverlay && m.foreground != nil {
		return m.overlayModel.View()
	}

	return m.background.View()
}
