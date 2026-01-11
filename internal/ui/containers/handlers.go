package containers

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/api/types/container"
)

func (cl *ContainerList) getSelectedContainerIDs() []string {
	selectedContainerIDs := make([]string, 0, len(cl.selectedContainers.selections))
	for id := range cl.selectedContainers.selections {
		selectedContainerIDs = append(selectedContainerIDs, id)
	}

	return selectedContainerIDs
}

func (cl *ContainerList) getSelectedContainerIndices() []int {
	selectedContainerIndices := make([]int, 0, len(cl.selectedContainers.selections))
	for _, index := range cl.selectedContainers.selections {
		selectedContainerIndices = append(selectedContainerIndices, index)
	}

	return selectedContainerIndices
}

func (cl *ContainerList) handlePauseContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIDs := cl.getSelectedContainerIDs()
		return PerformContainerOperation("pause", selectedContainerIDs)
	} else {
		selectedItem, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return PerformContainerOperation("pause", []string{selectedItem.ID})
		}
	}
	return nil
}

func (cl *ContainerList) handleUnpauseContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIDs := cl.getSelectedContainerIDs()
		return PerformContainerOperation("unpause", selectedContainerIDs)
	} else {
		selectedItem, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return PerformContainerOperation("unpause", []string{selectedItem.ID})
		}
	}
	return nil
}

func (cl *ContainerList) handleStartContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIDs := cl.getSelectedContainerIDs()
		return PerformContainerOperation("start", selectedContainerIDs)
	} else {
		selectedItem, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return PerformContainerOperation("start", []string{selectedItem.ID})
		}
	}
	return nil
}

func (cl *ContainerList) handleStopContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIDs := cl.getSelectedContainerIDs()
		return PerformContainerOperation("stop", selectedContainerIDs)
	} else {
		selectedItem, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return PerformContainerOperation("stop", []string{selectedItem.ID})
		}
	}
	return nil
}

func (cl *ContainerList) handleRemoveContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIndices := cl.getSelectedContainerIndices()

		var requestedContainersToDelete []*ContainerItem
		items := cl.list.Items()

		for _, index := range selectedContainerIndices {
			requestedContainer := items[index].(ContainerItem)
			requestedContainersToDelete = append(requestedContainersToDelete, &requestedContainer)
		}

		return func() tea.Msg {
			return MessageOpenDeleteConfirmationDialog{requestedContainersToDelete}
		}
	} else {
		item, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return func() tea.Msg {
				return MessageOpenDeleteConfirmationDialog{[]*ContainerItem{&item}}
			}
		}
	}

	return nil
}

func (cl *ContainerList) handleShowLogs() tea.Cmd {
	item, ok := cl.list.SelectedItem().(ContainerItem)
	if !ok {
		log.Print()
		return nil
	}

	if item.State != container.StateRunning {
		log.Printf("%s is not running...", item.Name)
		return nil
	}

	return OpenContainerLogs(&item)
}

func (cl *ContainerList) handleConfirmationOfRemoveContainers() tea.Cmd {
	if len(cl.selectedContainers.selections) > 0 {
		selectedContainerIDs := cl.getSelectedContainerIDs()
		return PerformContainerOperation("remove", selectedContainerIDs)
	} else {
		item, ok := cl.list.SelectedItem().(ContainerItem)
		if ok {
			return PerformContainerOperation("remove", []string{item.ID})
		}
	}
	return nil
}

func (cl *ContainerList) handleToggleSelection() {
	index := cl.list.Index()
	selectedItem, ok := cl.list.SelectedItem().(ContainerItem)
	if ok {
		isSelected := selectedItem.isSelected

		if isSelected {
			cl.selectedContainers.unselectContainerInList(selectedItem.ID)
		} else {
			cl.selectedContainers.selectContainerInList(selectedItem.ID, index)
		}

		selectedItem.isSelected = !isSelected
		cl.list.SetItem(index, selectedItem)
	}
}

func (cl *ContainerList) handleToggleSelectionOfAll() {
	allAlreadySelected := true
	items := cl.list.Items()

	for _, item := range items {
		if c, ok := item.(ContainerItem); ok {
			if _, selected := cl.selectedContainers.selections[c.ID]; !selected {
				allAlreadySelected = false
				break
			}
		}
	}

	if allAlreadySelected {
		// Unselect all items
		cl.selectedContainers = newSelectedContainers()

		for index, item := range cl.list.Items() {
			item, ok := item.(ContainerItem)
			if ok {
				item.isSelected = false
				cl.list.SetItem(index, item)
			}
		}
	} else {
		// Select all items
		cl.selectedContainers = newSelectedContainers()

		for index, item := range cl.list.Items() {
			item, ok := item.(ContainerItem)
			if ok {
				item.isSelected = true
				cl.list.SetItem(index, item)
				cl.selectedContainers.selectContainerInList(item.ID, index)
			}
		}
	}
}

func (cl *ContainerList) handleContainerOperationResult(msg MessageContainerOperationResult) {
	if msg.Error != nil {
		// For now, ignore errors to prevent UI freeze; TODO: show error in UI
		return
	}

	if msg.Operation == "remove" {
		// For remove, remove items from list
		items := cl.list.Items()
		// Collect indices to remove, in reverse order to avoid index shifting
		var indicesToRemove []int
		for i, item := range items {
			if c, ok := item.(ContainerItem); ok {
				for _, id := range msg.IDs {
					if c.ID == id {
						indicesToRemove = append([]int{i}, indicesToRemove...) // prepend to reverse
						break
					}
				}
			}
		}
		for _, index := range indicesToRemove {
			cl.list.RemoveItem(index)
		}
		return
	}

	var newState container.ContainerState
	switch msg.Operation {
	case "pause":
		newState = container.StatePaused
	case "unpause", "start":
		newState = container.StateRunning
	case "stop":
		newState = container.StateExited
	default:
		return
	}

	items := cl.list.Items()
	for i, item := range items {
		if c, ok := item.(ContainerItem); ok {
			for _, id := range msg.IDs {
				if c.ID == id {
					c.State = newState
					cl.list.SetItem(i, c)
					break
				}
			}
		}
	}
}
