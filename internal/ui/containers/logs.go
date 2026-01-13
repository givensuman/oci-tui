package containers

import (
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/client"
	"github.com/givensuman/containertui/internal/colors"
	"github.com/givensuman/containertui/internal/context"
	"github.com/givensuman/containertui/internal/ui/shared"
)

type (
	LogChunkMsg []byte
	LogErrorMsg error
)

type ContainerLogs struct {
	shared.Component
	container   *ContainerItem
	logs        client.Logs
	viewport    viewport.Model
	isFollowing bool
	logBuffer   strings.Builder

	style lipgloss.Style
}

func waitForLogs(reader io.Reader) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, 2048)
		n, err := reader.Read(buf)
		if err != nil {
			return LogErrorMsg(err)
		}
		return LogChunkMsg(buf[:n])
	}
}

func newContainerLogs(container *ContainerItem) *ContainerLogs {
	// 1. Get initial size from context to avoid the "Initializing" hang
	w, h := context.GetWindowSize()

	cl := &ContainerLogs{
		container:   container,
		logs:        context.GetClient().OpenLogs(container.ID),
		isFollowing: true,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Primary()).
			Padding(0, 1),
	}

	cl.viewport = viewport.New(0, 0)
	cl.setDimensions(w, h)

	return cl
}

func (cl *ContainerLogs) setDimensions(w, h int) {
	cl.WindowWidth = w
	cl.WindowHeight = h

	// Define log window size (80% width, 70% height)
	width := int(float64(w) * 0.8)
	height := int(float64(h) * 0.7)

	cl.style = cl.style.Width(width).Height(height)

	cl.viewport.Width = width - cl.style.GetHorizontalFrameSize()
	cl.viewport.Height = height - cl.style.GetVerticalFrameSize()
}

func (cl *ContainerLogs) Init() tea.Cmd {
	return waitForLogs(cl.logs)
}

func (cl *ContainerLogs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case LogChunkMsg:
		cl.logBuffer.Write(msg)
		cl.viewport.SetContent(cl.logBuffer.String())

		if cl.isFollowing {
			cl.viewport.GotoBottom()
		}
		cmds = append(cmds, waitForLogs(cl.logs))

	case tea.WindowSizeMsg:
		cl.setDimensions(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "f": // Toggle Follow/Tail
			cl.isFollowing = !cl.isFollowing
			if cl.isFollowing {
				cl.viewport.GotoBottom()
			}
		case "esc", "q":
			_ = cl.logs.Close()
			return cl, CloseOverlay()
		}
	}

	var cmd tea.Cmd
	cl.viewport, cmd = cl.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return cl, tea.Batch(cmds...)
}

func (cl *ContainerLogs) View() string {
	// Status bar to show if Tailing is on
	followStatus := " [Tail: ON] "
	if !cl.isFollowing {
		followStatus = " [Tail: OFF] "
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		cl.viewport.View(),
		lipgloss.NewStyle().Foreground(colors.Primary()).Render(followStatus),
	)
}
