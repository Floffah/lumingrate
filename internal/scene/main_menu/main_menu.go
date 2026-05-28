package mainmenu

import (
	"luminrate/internal/style"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

var (
	screenStyle         = style.Screen
	titleStyle          = style.MenuTitle
	buttonStyle         = style.MenuButton
	selectedButtonStyle = style.MenuButtonSelected
	panelStyle          = style.Panel
)

type screen int

const (
	menuScreen screen = iota
	settingsScreen
)

type StartGameMsg struct{}

type model struct {
	width    int
	height   int
	screen   screen
	selected int
}

func NewModel() tea.Model {
	return model{
		width:  defaultWidth,
		height: defaultHeight,
	}
}

func (m model) Init() tea.Cmd {
	return tea.RequestWindowSize
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch m.screen {
		case menuScreen:
			return m.updateMenu(msg)
		case settingsScreen:
			return m.updateSubscreen(msg)
		}
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		m.selected--
		if m.selected < 0 {
			m.selected = len(menuItems) - 1
		}
	case "down", "j":
		m.selected++
		if m.selected >= len(menuItems) {
			m.selected = 0
		}
	case "enter":
		switch menuItem(m.selected) {
		case startGameItem:
			return m, startGame
		case settingsItem:
			m.screen = settingsScreen
		case exitItem:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) updateSubscreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "backspace", "enter":
		m.screen = menuScreen
	}

	return m, nil
}

func (m model) View() tea.View {
	content := m.renderMenu()
	if m.screen == settingsScreen {
		content = m.renderSettings()
	}

	view := tea.NewView(m.place(content))
	view.AltScreen = true
	view.BackgroundColor = style.Dark
	view.ForegroundColor = style.White
	view.WindowTitle = "Lumingrate"
	return view
}

func startGame() tea.Msg {
	return StartGameMsg{}
}

func (m model) renderMenu() string {
	rows := []string{
		titleStyle.Render("LUMINGRATE"),
	}

	for i, item := range menuItems {
		style := buttonStyle
		if i == m.selected {
			style = selectedButtonStyle
		}
		rows = append(rows, style.Render(item.label))
	}

	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func (m model) place(content string) string {
	width := max(m.width, 1)
	height := max(m.height, 1)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceStyle(screenStyle),
	)
}

type menuItem int

const (
	startGameItem menuItem = iota
	settingsItem
	exitItem
)

var menuItems = []struct {
	label string
}{
	{label: "START GAME"},
	{label: "SETTINGS"},
	{label: "EXIT"},
}
