package gamescene

import (
	"context"
	"luminrate/internal/engine"
	"luminrate/internal/story"
	"luminrate/internal/style"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
	prompt        = "> "
)

var (
	screenStyle  = style.Screen
	logStyle     = style.Text
	inputStyle   = style.Input
	cursorStyle  = style.Cursor
	commandStyle = style.Command
	asideStyle   = style.Aside
)

type entryKind int

const (
	entryNarration entryKind = iota
	entryCommand
	entryAside
	entryRaw
)

type entry struct {
	kind entryKind
	text string
}

type Model struct {
	width    int
	height   int
	cancel   context.CancelFunc
	commands chan<- engine.Command
	events   <-chan engine.Event

	log    []entry
	input  string
	cursor int

	inputVisible bool
	inputMessage string
}

func NewModel() Model {
	return NewModelWithSize(defaultWidth, defaultHeight)
}

func NewModelWithSize(width, height int) Model {
	ctx, cancel := context.WithCancel(context.Background())
	gameEngine := story.NewEngine()
	go gameEngine.Run(ctx)

	return Model{
		width:    max(width, 1),
		height:   max(height, 1),
		cancel:   cancel,
		commands: gameEngine.Commands(),
		events:   gameEngine.Events(),

		inputVisible: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, waitForEngineEvent(m.events))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case engineEventMsg:
		m.applyEvent(engine.Event(msg))
		return m, waitForEngineEvent(m.events)
	case engineStoppedMsg:
		return m, nil
	case tea.KeyPressMsg:
		return m.updateInput(msg)
	}

	return m, nil
}

func (m Model) updateInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	msgStr := msg.String()

	if msgStr == "ctrl+c" {
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	}

	if !m.inputVisible {
		return m, nil
	}

	switch msgStr {
	case "enter":
		return m.submitInput()
	case "tab":
		m.inputMessage = ""
		return m, sendEngineCommand(m.commands, engine.Command{Kind: engine.CommandTab, Text: m.input})
	case "backspace":
		if m.cursor > 0 {
			m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
			m.cursor--
			m.inputMessage = ""
		}
	case "delete":
		if m.cursor < len(m.input) {
			m.input = m.input[:m.cursor] + m.input[m.cursor+1:]
			m.inputMessage = ""
		}
	case "left":
		if m.cursor > 0 {
			m.cursor--
		}
	case "right":
		if m.cursor < len(m.input) {
			m.cursor++
		}
	case "home":
		m.cursor = 0
	case "end":
		m.cursor = len(m.input)
	default:
		if text := msg.Key().Text; text != "" {
			m.input = m.input[:m.cursor] + text + m.input[m.cursor:]
			m.cursor += len(text)
			m.inputMessage = ""
		}
	}

	return m, nil
}

func (m Model) submitInput() (tea.Model, tea.Cmd) {
	command := strings.TrimSpace(m.input)

	if command != "" {
		m.log = append(m.log, entry{kind: entryCommand, text: command})
	}
	m.input = ""
	m.cursor = 0
	m.inputMessage = ""

	return m, sendEngineCommand(m.commands, engine.Command{Kind: engine.CommandSubmit, Text: command})
}

func (m *Model) applyEvent(event engine.Event) {
	switch event.Kind {
	case engine.EventInputMessage:
		m.inputMessage = event.Text
	case engine.EventHideInput:
		m.inputVisible = false
	case engine.EventShowInput:
		m.inputVisible = true
	case engine.EventInsertInput:
		m.insertInput(event.Text)
	default:
		m.log = append(m.log, entryFromEvent(event))
	}
}

func (m *Model) insertInput(text string) {
	if text == "" {
		return
	}

	m.input = m.input[:m.cursor] + text + m.input[m.cursor:]
	m.cursor += len(text)
	m.inputMessage = ""
}

func entryFromEvent(event engine.Event) entry {
	switch event.Kind {
	case engine.EventAside:
		return entry{kind: entryAside, text: event.Text}
	case engine.EventRaw:
		return entry{kind: entryRaw, text: event.Text}
	}
	return entry{kind: entryNarration, text: event.Text}
}

func (m Model) View() tea.View {
	view := tea.NewView(m.renderLogInput())
	view.AltScreen = true
	view.BackgroundColor = style.ColourDark
	view.ForegroundColor = style.ColourWhite
	view.WindowTitle = "Lumingrate"
	return view
}

func (m Model) renderLogInput() string {
	width := max(m.width, 1)
	height := max(m.height, 1)
	messageHeight := 1
	inputHeight := 1
	logHeight := max(height-inputHeight-messageHeight, 1)

	log := m.renderLog(width, logHeight)
	input := m.renderInputLine(width)
	message := m.renderInputMessage(width)

	return lipgloss.JoinVertical(lipgloss.Top, log, input, message)
}

func (m Model) renderLog(width, height int) string {
	heightWithPadding := height - 2
	lines := m.renderLogLines(width)

	if len(lines) > heightWithPadding {
		lines = lines[len(lines)-heightWithPadding:]
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Left,
		lipgloss.Bottom,
		strings.Join(m.padLines(lines, 1, width), "\n"),
		lipgloss.WithWhitespaceStyle(screenStyle),
	)
}

func (m Model) renderLogLines(width int) []string {
	lines := make([]string, 0, len(m.log))
	rawLine := ""
	for _, item := range m.log {
		if item.kind == entryRaw {
			parts := strings.Split(item.text, "\n")
			rawLine += parts[0]
			for _, part := range parts[1:] {
				lines = append(lines, style.Background.Width(width).Render(rawLine))
				rawLine = part
			}
			continue
		}

		if rawLine != "" {
			lines = append(lines, style.Background.Width(width).Render(rawLine))
			rawLine = ""
		}

		rendered := m.renderEntry(item, width)
		lines = append(lines, strings.Split(rendered, "\n")...)
	}

	if rawLine != "" {
		lines = append(lines, style.Background.Width(width).Render(rawLine))
	}

	return lines
}

func (m Model) padLines(lines []string, size int, fullWidth int) []string {
	padded := make([]string, len(lines)+2)
	padded[0] = screenStyle.Height(size).Width(fullWidth).Render("")
	for i, line := range lines {
		padded[i+1] = screenStyle.Padding(0, size).Render(line)
	}
	padded[len(padded)-1] = screenStyle.Height(size).Width(fullWidth).Render("")
	return padded
}

func (m Model) renderEntry(item entry, width int) string {
	entryStyle := logStyle
	switch item.kind {
	case entryCommand:
		entryStyle = commandStyle
	case entryAside:
		entryStyle = asideStyle
	case entryRaw:
		return style.Background.Width(width).Render(item.text)
	}

	return entryStyle.Width(width).Background(style.ColourBackground).Render(lipgloss.Wrap(item.text, width, " "))
}

func (m Model) renderInputLine(width int) string {
	content := ""
	if m.inputVisible {
		content = prompt + m.renderInput()
	}

	return inputStyle.Width(width).Padding(0, 1).Render(content)
}

func (m Model) renderInputMessage(width int) string {
	return asideStyle.Width(width).Padding(0, 1).Render(m.inputMessage)
}

func (m Model) renderInput() string {
	before := m.input[:m.cursor]
	after := m.input[m.cursor:]
	cursor := " "
	if after != "" {
		cursor = after[:1]
		after = after[1:]
	}

	return before + cursorStyle.Render(cursor) + after
}

type engineEventMsg engine.Event

type engineStoppedMsg struct{}

func waitForEngineEvent(events <-chan engine.Event) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-events
		if !ok {
			return engineStoppedMsg{}
		}
		return engineEventMsg(event)
	}
}

func sendEngineCommand(commands chan<- engine.Command, command engine.Command) tea.Cmd {
	return func() tea.Msg {
		commands <- command
		return nil
	}
}
