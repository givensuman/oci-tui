// Package ui implements the terminal user interface
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
)

// Model represents the application state
type Model struct {
	width  int
	height int
}

// NewModel creates a new model with default values
func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		context.SetWindowSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	// Create styles
	titleStyle := lipgloss.NewStyle().
		Foreground(colors.Primary()).
		Bold(true).
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(colors.Yellow()).
		Align(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(colors.Green()).
		Align(lipgloss.Center)

	// Content
	title := titleStyle.Render("ContainerTUI")
	subtitle := subtitleStyle.Render("A terminal UI for managing container lifecycles")
	help := helpStyle.Render("Press 'q' to quit")

	// Combine content
	content := strings.Join([]string{title, "", subtitle, "", help}, "\n")

	// Create a full-screen container
	container := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	return container.Render(content)
}

func Start() error {
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
