package containers

import (
	tea "github.com/charmbracelet/bubbletea"
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
		return "Confirm"
	case decline:
		return "Decline"
	}

	return "Unknown"
}

type DeleteConfirmation struct {
	item *ContainerItem
}

func NewDeleteConfirmation(item *ContainerItem) DeleteConfirmation {
	return DeleteConfirmation{item}
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
	return dc, nil
}

func (dc DeleteConfirmation) View() string {
	return "Hello World!"
}
