package containers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
)

type buttonOption int

const (
	confirm buttonOption = iota
	decline
)

func (bo buttonOption) String() string {
	switch bo {
	case confirm:
		return "Delete"
	case decline:
		return "Cancel"
	}

	return "Unknown"
}

type DeleteConfirmation struct {
	style               lipgloss.Style
	requestedContainers []*ContainerItem
	hoveredButtonOption buttonOption
}

func newDeleteConfirmation(requestedContainers ...*ContainerItem) DeleteConfirmation {
	style := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder(), true, true).
		BorderForeground(colors.Red())

	return DeleteConfirmation{
		style:               style,
		requestedContainers: requestedContainers,
		hoveredButtonOption: decline,
	}
}

func (dc DeleteConfirmation) Init() tea.Cmd {
	return nil
}

func (dc DeleteConfirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEscape.String(), tea.KeyEsc.String():
			cmds = append(cmds, func() tea.Msg { return MessageCloseOverlay{} })

		case tea.KeyTab.String():
			switch dc.hoveredButtonOption {
			case confirm:
				dc.hoveredButtonOption = decline
			case decline:
				dc.hoveredButtonOption = confirm
			}

		case tea.KeyEnter.String():
			if dc.hoveredButtonOption == confirm {
				cmds = append(cmds, func() tea.Msg { return MessageConfirmDelete{} })
			}

			cmds = append(cmds, func() tea.Msg { return MessageCloseOverlay{} })
		}
	}

	return dc, tea.Batch(cmds...)
}

func (dc DeleteConfirmation) View() string {
	hoveredButtonStyle := lipgloss.NewStyle().
		Background(colors.Primary()).
		Bold(true).
		Foreground(colors.Black()).
		Padding(0, 1)

	defaultButtonStyle := lipgloss.NewStyle().
		Background(colors.Gray()).
		Foreground(colors.White()).
		Padding(0, 1)

	confirmButton := confirm.String()
	declineButton := decline.String()

	if dc.hoveredButtonOption == confirm {
		confirmButton = hoveredButtonStyle.Render(confirmButton)
		declineButton = defaultButtonStyle.Render(declineButton)
	} else {
		confirmButton = defaultButtonStyle.Render(confirmButton)
		declineButton = hoveredButtonStyle.Render(declineButton)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		declineButton,
		"   ",
		confirmButton,
	)

	var message string
	if len(dc.requestedContainers) == 1 {
		message = fmt.Sprintf("Are you sure you want to delete %s?", dc.requestedContainers[0].Name)
	} else {
		message = fmt.Sprintf("Are you sure you want to delete the %d selected containers?", len(dc.requestedContainers))
	}

	return dc.style.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"   ",
		buttons,
	))
}
