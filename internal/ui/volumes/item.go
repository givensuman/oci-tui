package volumes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type VolumeItem struct {
	Volume     client.Volume
	isSelected bool
}

var (
	_ list.Item        = (*VolumeItem)(nil)
	_ list.DefaultItem = (*VolumeItem)(nil)
)

func newDefaultDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d = shared.ChangeDelegateStyles(d)

	return d
}

func (v VolumeItem) getIsSelectedIcon() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		switch v.isSelected {
		case true:
			return "[x]"
		case false:
			return "[ ]"
		}
	case false: // Use nerd fonts
		switch v.isSelected {
		case true:
			return " "
		case false:
			return " "
		}
	}

	return "[ ]"
}

func (v VolumeItem) getTitleOrnament() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		return ""
	case false: // Use nerd fonts
		return " "
	}

	return ""
}

func (v VolumeItem) Title() string {
	titleOrnament := v.getTitleOrnament()

	title := fmt.Sprintf("%s %s (%s)", titleOrnament, v.Volume.Name, v.Volume.Driver)
	title = lipgloss.NewStyle().
		Foreground(colors.Gray()).
		Render(title)

	statusIcon := v.getIsSelectedIcon()
	var isSelectedColor lipgloss.Color
	switch v.isSelected {
	case true:
		isSelectedColor = colors.Blue()
	case false:
		isSelectedColor = colors.White()
	}
	statusIcon = lipgloss.NewStyle().
		Foreground(isSelectedColor).
		Render(statusIcon)

	return fmt.Sprintf("%s %s", statusIcon, title)
}

func (v VolumeItem) Description() string {
	return fmt.Sprintf("Mountpoint: %s", v.Volume.Mountpoint)
}

func (v VolumeItem) FilterValue() string {
	return v.Title()
}
