package shared

import "github.com/charmbracelet/lipgloss"

type Dimensions struct {
	Width         int
	Height        int
	OffsetX       int
	OffsetY       int
	ContentWidth  int
	ContentHeight int
}

type LayoutManager struct {
	windowWidth  int
	windowHeight int
}

func NewLayoutManager(width, height int) LayoutManager {
	return LayoutManager{
		windowWidth:  width,
		windowHeight: height,
	}
}

func (lm *LayoutManager) UpdateDimensions(width, height int) {
	lm.windowWidth = width
	lm.windowHeight = height
}

func (lm LayoutManager) Calculate(ratio WindowRatio, style lipgloss.Style) Dimensions {
	width := AdjustedWidth(lm.windowWidth, ratio)
	height := AdjustedHeight(lm.windowHeight, ratio)

	offsetX := (lm.windowWidth - width) / 2
	offsetY := (lm.windowHeight - height) / 2

	frameH := style.GetHorizontalFrameSize()
	frameV := style.GetVerticalFrameSize()

	return Dimensions{
		Width:         width,
		Height:        height,
		OffsetX:       offsetX,
		OffsetY:       offsetY,
		ContentWidth:  width - frameH,
		ContentHeight: height - frameV,
	}
}

func (lm LayoutManager) CalculateFullscreen(style lipgloss.Style) Dimensions {
	return lm.Calculate(RatioFullscreen, style)
}

func (lm LayoutManager) CalculateModal(style lipgloss.Style) Dimensions {
	return lm.Calculate(RatioModal, style)
}

func (lm LayoutManager) CalculateLargeOverlay(style lipgloss.Style) Dimensions {
	return lm.Calculate(RatioLargeOverlay, style)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
