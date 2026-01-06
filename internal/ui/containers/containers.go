// Package containers defines the containers list component.
package containers

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keybindings struct {
	pauseContainer       key.Binding
	unpauseContainer     key.Binding
	startContainer       key.Binding
	stopContainer        key.Binding
	toggleSelection      key.Binding
	toggleSelectionOfAll key.Binding
}

func newKeybindings() *keybindings {
	return &keybindings{
		pauseContainer: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause container"),
		),
		unpauseContainer: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "unpause container"),
		),
		startContainer: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start container"),
		),
		stopContainer: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "stop container"),
		),
		toggleSelection: key.NewBinding(
			key.WithKeys(tea.KeySpace.String()),
			key.WithHelp("space", "toggle selection"),
		),
		toggleSelectionOfAll: key.NewBinding(
			key.WithKeys(tea.KeyCtrlA.String()),
			key.WithHelp("ctrl+a", "toggle selection of all"),
		),
	}
}

// selectedContainers is map of a container's ID to
// its index in the list
type selectedContainers map[string]int

func (sc selectedContainers) selectContainerInList(id string, index int) selectedContainers {
	sc[id] = index
	return sc
}

func (sc selectedContainers) unselectContainerInList(id string) selectedContainers {
	delete(sc, id)
	return sc
}
