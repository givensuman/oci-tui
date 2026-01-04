// Package ui provides the main user interface for the application.
package ui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/containers"
)

type Model struct {
	containersList containers.ListModel
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			context.SetWindowSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.containersList, cmd = m.containersList.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	v := tea.NewView(m.containersList.View())
	v.AltScreen = true

	return v
}

func Start() {
	m := Model{
		containersList: containers.NewListModel(),
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
