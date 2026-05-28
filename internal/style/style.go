package style

import "charm.land/lipgloss/v2"

const (
	ButtonWidth = 22
	PanelWidth  = 48
)

var (
	Dark  = lipgloss.Color("#000000")
	White = lipgloss.Color("#F7F7F2")
	Muted = lipgloss.Color("#777777")

	Screen = lipgloss.NewStyle().
		Background(Dark).
		Foreground(White)

	Text = lipgloss.NewStyle().
		Background(Dark).
		Foreground(White)

	TextItalic = Text.Italic(true)

	Input = lipgloss.NewStyle().
		Background(Dark).
		Foreground(White)

	Cursor = lipgloss.NewStyle().
		Foreground(Dark).
		Background(White)

	Command = lipgloss.NewStyle().
		Background(Dark).
		Foreground(White).
		Bold(true)

	Aside = lipgloss.NewStyle().
		Background(Dark).
		Foreground(Muted)

	MenuTitle = lipgloss.NewStyle().
			Foreground(White).
			Background(Dark).
			Bold(true).
			Align(lipgloss.Center).
			Width(ButtonWidth).
			PaddingBottom(1)

	MenuButton = lipgloss.NewStyle().
			Foreground(White).
			Background(Dark).
			Align(lipgloss.Center).
			Width(ButtonWidth).
			Padding(0, 2)

	MenuButtonSelected = MenuButton.
				Foreground(Dark).
				Background(White)

	Panel = lipgloss.NewStyle().
		Foreground(White).
		Background(Dark).
		Align(lipgloss.Center).
		Width(PanelWidth)
)
