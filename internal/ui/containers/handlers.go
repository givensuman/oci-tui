package containers

import (
	"github.com/givensuman/containertui/internal/context"
	"github.com/moby/moby/api/types/container"
)

func (m Model) getSelectedContainerIDs() []string {
	selectedContainerIDs := make([]string, len(m.selectedContainers))
	for id := range m.selectedContainers {
		selectedContainerIDs = append(selectedContainerIDs, id)
	}

	return selectedContainerIDs
}

func (m Model) getSelectedContainerIndices() []int {
	selectedContainerIndices := make([]int, len(m.selectedContainers))
	for _, index := range m.selectedContainers {
		selectedContainerIndices = append(selectedContainerIndices, index)
	}

	return selectedContainerIndices
}

func (m Model) handlePauseContainers() Model {
	if len(m.selectedContainers) > 0 {
		selectedContainerIDs := m.getSelectedContainerIDs()
		selectedContainerIndices := m.getSelectedContainerIndices()

		context.GetClient().PauseContainers(selectedContainerIDs)
		items := m.list.Items()
		for _, index := range selectedContainerIndices {
			item := items[index].(ContainerItem)
			item.State = container.StatePaused
			m.list.SetItem(index, item)
		}
	} else {
		selectedItem, ok := m.list.SelectedItem().(ContainerItem)
		if ok {
			context.GetClient().PauseContainer(selectedItem.ID)
			selectedItem.State = container.StatePaused
			m.list.SetItem(m.list.Index(), selectedItem)
		}
	}

	return m
}

func (m Model) handleUnpauseContainers() Model {
	if len(m.selectedContainers) > 0 {
		selectedContainerIDs := m.getSelectedContainerIDs()
		selectedContainerIndices := m.getSelectedContainerIndices()

		context.GetClient().UnpauseContainers(selectedContainerIDs)
		items := m.list.Items()
		for _, index := range selectedContainerIndices {
			item := items[index].(ContainerItem)
			item.State = container.StateRunning
			m.list.SetItem(index, item)
		}
	} else {
		selectedItem, ok := m.list.SelectedItem().(ContainerItem)
		if ok {
			context.GetClient().UnpauseContainer(selectedItem.ID)
			selectedItem.State = container.StateRunning
			m.list.SetItem(m.list.Index(), selectedItem)
		}
	}

	return m
}

func (m Model) handleStartContainers() Model {
	if len(m.selectedContainers) > 0 {
		selectedContainerIDs := m.getSelectedContainerIDs()
		selectedContainerIndices := m.getSelectedContainerIndices()

		context.GetClient().StartContainers(selectedContainerIDs)
		items := m.list.Items()
		for _, index := range selectedContainerIndices {
			item := items[index].(ContainerItem)
			item.State = container.StateRunning
			m.list.SetItem(index, item)
		}
	} else {
		selectedItem, ok := m.list.SelectedItem().(ContainerItem)
		if ok {
			context.GetClient().StartContainer(selectedItem.ID)
			selectedItem.State = container.StateRunning
			m.list.SetItem(m.list.Index(), selectedItem)
		}
	}

	return m
}

func (m Model) handleStopContainers() Model {
	if len(m.selectedContainers) > 0 {
		selectedContainerIDs := m.getSelectedContainerIDs()
		selectedContainerIndices := m.getSelectedContainerIndices()

		context.GetClient().StopContainers(selectedContainerIDs)
		items := m.list.Items()
		for _, index := range selectedContainerIndices {
			item := items[index].(ContainerItem)
			item.State = container.StateExited
			m.list.SetItem(index, item)
		}
	} else {
		selectedItem, ok := m.list.SelectedItem().(ContainerItem)
		if ok {
			context.GetClient().StopContainer(selectedItem.ID)
			selectedItem.State = container.StateExited
			m.list.SetItem(m.list.Index(), selectedItem)
		}
	}

	return m
}
