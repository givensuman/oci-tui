package shared

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
)

// DialogButton defines a button in the dialog
type DialogButton struct {
	Label  string
	Action SmartDialogAction // Empty action means Cancel/Close
	IsSafe bool              // True = Primary/Safe color, False = Danger/Red
}

type SmartDialog struct {
	Component
	style          lipgloss.Style
	message        string
	buttons        []DialogButton
	selectedButton int
	width          int
	height         int
}

var (
	_ tea.Model      = (*SmartDialog)(nil)
	_ ComponentModel = (*SmartDialog)(nil)
)

// NewSmartDialog creates a generic confirmation/warning dialog
func NewSmartDialog(message string, buttons []DialogButton) SmartDialog {
	width, height := context.GetWindowSize()

	style := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder(), true, true).
		BorderForeground(colors.Primary()).
		Align(lipgloss.Center)

	lm := NewLayoutManager(width, height)
	dims := lm.CalculateModal(style)
	style = style.Width(dims.Width).Height(dims.Height)

	// Ensure there is at least one button (Cancel)
	if len(buttons) == 0 {
		buttons = []DialogButton{{Label: "Cancel", IsSafe: true}}
	}

	return SmartDialog{
		style:          style,
		message:        message,
		buttons:        buttons,
		selectedButton: 0, // Default to first button (usually Cancel/Safe)
		width:          width,
		height:         height,
	}
}

func (d *SmartDialog) UpdateWindowDimensions(msg tea.WindowSizeMsg) {
	d.width = msg.Width
	d.height = msg.Height

	lm := NewLayoutManager(msg.Width, msg.Height)
	dims := lm.CalculateModal(d.style)
	d.style = d.style.Width(dims.Width).Height(dims.Height)
}

func (d SmartDialog) Init() tea.Cmd {
	return nil
}

func (d SmartDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.UpdateWindowDimensions(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return d, func() tea.Msg { return CloseDialogMessage{} }

		case "tab", "right", "l":
			d.selectedButton = (d.selectedButton + 1) % len(d.buttons)

		case "shift+tab", "left", "h":
			d.selectedButton = (d.selectedButton - 1 + len(d.buttons)) % len(d.buttons)

		case "enter":
			btn := d.buttons[d.selectedButton]
			if btn.Action.Type == "" {
				return d, func() tea.Msg { return CloseDialogMessage{} }
			}
			return d, func() tea.Msg { return ConfirmationMessage{Action: btn.Action} }
		}
	}

	return d, tea.Batch(cmds...)
}

func (d SmartDialog) View() string {
	var buttonViews []string

	defaultButtonStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Foreground(colors.Text()).
		Background(colors.Muted())

	activeSafeStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Bold(true).
		Foreground(colors.Text()).
		Background(colors.Primary())

	activeDangerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Bold(true).
		Foreground(colors.Text()).
		Background(colors.Error())

	for i, btn := range d.buttons {
		var btnStyle lipgloss.Style

		if i == d.selectedButton {
			if btn.IsSafe {
				btnStyle = activeSafeStyle
			} else {
				btnStyle = activeDangerStyle
			}
		} else {
			btnStyle = defaultButtonStyle
		}

		buttonViews = append(buttonViews, btnStyle.Render(btn.Label))
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, buttonViews...)

	// Change border color if danger button is selected
	currentBtn := d.buttons[d.selectedButton]
	renderStyle := d.style
	if !currentBtn.IsSafe {
		renderStyle = renderStyle.BorderForeground(colors.Error())
	}

	return renderStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		d.message,
		"", // Spacer
		buttons,
	))
}
