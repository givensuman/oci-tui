package containers

// MessageCloseOverlay indicates the overlay should display
// its background
type MessageCloseOverlay struct{}

// MessageOpenDeleteConfirmationDialog indicates the user
// has requested to delete an item in the ContainerList
type MessageOpenDeleteConfirmationDialog struct {
	item *ContainerItem
}

// MessageConfirmDelete indicates the user confirmed
// they wish to delete an item in the ContainerList
type MessageConfirmDelete struct {
	item *ContainerItem
}
