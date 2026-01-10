package containers

import (
	tea "github.com/charmbracelet/bubbletea"
)

// MessageCloseOverlay indicates the overlay should display
// its background
type MessageCloseOverlay struct{}

func CloseOverlay() tea.Cmd {
	return func() tea.Msg {
		return MessageCloseOverlay{}
	}
}

// MessageOpenContainerLogs indicates the user
// has requested to view the logs of the selected container
type MessageOpenContainerLogs struct {
	container *ContainerItem
}

func OpenContainerLogs(container *ContainerItem) tea.Cmd {
	return func() tea.Msg {
		return MessageOpenContainerLogs{container}
	}
}

// MessageOpenDeleteConfirmationDialog indicates the user
// has requested to delete an item in the ContainerList
type MessageOpenDeleteConfirmationDialog struct {
	requestedContainersToDelete []*ContainerItem
}

// MessageConfirmDelete indicates the user confirmed
// they wish to delete an item in the ContainerList
type MessageConfirmDelete struct{}
