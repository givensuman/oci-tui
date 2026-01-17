package networks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type NetworkItem struct {
	Network    client.Network
	isSelected bool
}

var (
	_ list.Item        = (*NetworkItem)(nil)
	_ list.DefaultItem = (*NetworkItem)(nil)
)

func newDefaultDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d = shared.ChangeDelegateStyles(d)

	return d
}

func (n NetworkItem) getIsSelectedIcon() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		switch n.isSelected {
		case true:
			return "[x]"
		case false:
			return "[ ]"
		}
	case false: // Use nerd fonts
		switch n.isSelected {
		case true:
			return " "
		case false:
			return " "
		}
	}

	return "[ ]"
}

func (n NetworkItem) getTitleOrnament() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		return ""
	case false: // Use nerd fonts
		return " "
	}

	return ""
}

func (n NetworkItem) Title() string {
	titleOrnament := n.getTitleOrnament()

	title := fmt.Sprintf("%s %s (%s)", titleOrnament, n.Network.Name, n.Network.Driver)
	title = lipgloss.NewStyle().
		Foreground(colors.Muted()).
		Render(title)

	statusIcon := n.getIsSelectedIcon()
	var isSelectedColor lipgloss.Color
	switch n.isSelected {
	case true:
		isSelectedColor = colors.Selected()
	case false:
		isSelectedColor = colors.Text()
	}
	statusIcon = lipgloss.NewStyle().
		Foreground(isSelectedColor).
		Render(statusIcon)

	return fmt.Sprintf("%s %s", statusIcon, title)
}

func (n NetworkItem) Description() string {
	return fmt.Sprintf("ID: %s | Scope: %s", n.Network.ID[:12], n.Network.Scope)
}

func (n NetworkItem) FilterValue() string {
	return n.Title()
}
