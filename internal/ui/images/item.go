package images

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type ImageItem struct {
	Image      client.Image
	isSelected bool
}

var (
	_ list.Item        = (*ImageItem)(nil)
	_ list.DefaultItem = (*ImageItem)(nil)
)

func newDefaultDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate = shared.ChangeDelegateStyles(delegate)

	return delegate
}

func (i ImageItem) getIsSelectedIcon() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		switch i.isSelected {
		case true:
			return "[x]"
		case false:
			return "[ ]"
		}
	case false: // Use nerd fonts
		switch i.isSelected {
		case true:
			return " "
		case false:
			return " "
		}
	}

	return "[ ]"
}

func (i ImageItem) getTitleOrnament() string {
	switch context.GetConfig().NoNerdFonts {
	case true: // Don't use nerd fonts
		return ""
	case false: // Use nerd fonts
		return " "
	}

	return ""
}

func (i ImageItem) FilterValue() string {
	return i.Title()
}

func (i ImageItem) Title() string {
	var repoTag string
	if len(i.Image.RepoTags) > 0 {
		repoTag = i.Image.RepoTags[0]
	} else {
		repoTag = "<none>"
	}

	titleOrnament := i.getTitleOrnament()
	sizeMB := float64(i.Image.Size) / 1024 / 1024

	title := fmt.Sprintf("%s %s (%.2fMB)", titleOrnament, repoTag, sizeMB)
	title = lipgloss.NewStyle().
		Foreground(colors.Gray()).
		Render(title)

	statusIcon := i.getIsSelectedIcon()
	var isSelectedColor lipgloss.Color
	switch i.isSelected {
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

func (i ImageItem) Description() string {
	shortID := i.Image.ID
	if len(shortID) > 12 {
		shortID = shortID[7:19] // Remove "sha256:" prefix and take first 12 chars
	}
	return fmt.Sprintf("ID: %s", shortID)
}
