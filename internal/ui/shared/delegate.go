package shared

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/givensuman/containertui/internal/colors"
)

func ChangeDelegateStyles(d list.DefaultDelegate) list.DefaultDelegate {
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		BorderLeftForeground(colors.Primary()).
		Foreground(colors.Primary())

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		BorderLeftForeground(colors.Primary()).
		Foreground(colors.Primary())

	d.Styles.DimmedDesc = d.Styles.DimmedDesc.
		Foreground(colors.Gray()).
		Bold(false)

	d.Styles.FilterMatch = d.Styles.FilterMatch.
		Foreground(colors.Primary()).
		Bold(true)

	return d
}
