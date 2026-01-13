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
	style       lipgloss.Style
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

	lm := shared.NewLayoutManager(w, h)
	dims := lm.CalculateLargeOverlay(cl.style)

	cl.style = cl.style.Width(dims.Width).Height(dims.Height)
	cl.viewport.Width = dims.ContentWidth
	cl.viewport.Height = dims.ContentHeight
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
	// TODO: Improve this! Better bar design, scroll handling, etc.
	followStatus := " [Tail: ON] "
	if !cl.isFollowing {
		followStatus = " [Tail: OFF] "
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		cl.style.Render(cl.viewport.View(), lipgloss.NewStyle().Foreground(colors.Primary()).Render(followStatus)),
	)
}
