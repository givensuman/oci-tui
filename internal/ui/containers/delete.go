package containers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/moby/moby/api/types/container"
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
	item                *ContainerItem
	hoveredButtonOption buttonOption
}

func newDeleteConfirmation(item *ContainerItem) DeleteConfirmation {
	style := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder(), true, true).
		BorderForeground(colors.Red())

	return DeleteConfirmation{
		style:               style,
		item:                item,
		hoveredButtonOption: decline,
	}
}

func (dc *DeleteConfirmation) Delete() {
	if dc.item.State == container.StateRunning {
		context.GetClient().StopContainer(dc.item.ID)
	}

	context.GetClient().RemoveContainer(dc.item.ID)
}

func (dc DeleteConfirmation) Init() tea.Cmd {
	return nil
}

func (dc DeleteConfirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
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
				dc.Delete()
				cmds = append(cmds, func() tea.Msg { return MessageConfirmDelete{dc.item} })
			}

			cmds = append(cmds, func() tea.Msg { return MessageCloseOverlay{} })
		}
	}

	return dc, tea.Batch(cmds...)
}

func (dc DeleteConfirmation) View() string {
	hoveredButtonStyle := lipgloss.NewStyle()
	defaultButtonStyle := lipgloss.NewStyle()

	confirmButton := confirm.String()
	declineButton := decline.String()

	if dc.hoveredButtonOption == confirm {
		confirmButton = hoveredButtonStyle.Render(confirmButton)
		declineButton = defaultButtonStyle.Render(declineButton)
	} else {
		confirmButton = defaultButtonStyle.Render(confirmButton)
		declineButton = hoveredButtonStyle.Render(declineButton)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, declineButton, confirmButton)

	return dc.style.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("Are you sure you want to delete %s?", dc.item.Name),
		buttons,
	))
}
