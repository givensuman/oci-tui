package containers

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/givensuman/containertui/internal/client"
)

// ContainerItem wraps a client.Container for use in the list.
type ContainerItem struct {
	container client.Container
}

// FilterValue implements list.Item.
func (c ContainerItem) FilterValue() string {
	return c.container.Name
}

// Title returns the title for the list item.
func (c ContainerItem) Title() string {
	return c.container.Name
}

// Description returns the description for the list item.
func (c ContainerItem) Description() string {
	return fmt.Sprintf("%s - %s", c.container.Image, c.container.State)
}

// ListModel wraps the bubbles list model.
type ListModel struct {
	list list.Model
}
