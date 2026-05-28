package style

import "charm.land/lipgloss/v2"

const (
	ButtonWidth = 22
	PanelWidth  = 48
)

var (
	ColourDark  = lipgloss.Color("#000000")
	ColourWhite = lipgloss.Color("#F7F7F2")
	ColourMuted = lipgloss.Color("#777777")

	ColourBackground      = ColourDark
	ColourForeground      = ColourWhite
	ColourMutedForeground = ColourMuted

	Background = lipgloss.NewStyle().
			Background(ColourBackground)
	Foreground = lipgloss.NewStyle().
			Foreground(ColourForeground)
	Muted = lipgloss.NewStyle().
		Foreground(ColourMuted)

	Screen = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground)

	Text = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground)

	TextItalic = Text.Italic(true)

	Input = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground)

	Cursor = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground)

	Command = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground).
		Bold(true)

	Aside = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourMuted)

	MenuTitle = lipgloss.NewStyle().
			Background(ColourBackground).
			Foreground(ColourForeground).
			Bold(true).
			Align(lipgloss.Center).
			Width(ButtonWidth).
			PaddingBottom(1)

	MenuButton = lipgloss.NewStyle().
			Background(ColourBackground).
			Foreground(ColourForeground).
			Align(lipgloss.Center).
			Width(ButtonWidth).
			Padding(0, 2)

	MenuButtonSelected = MenuButton.
				Background(ColourForeground).
				Foreground(ColourBackground)

	Panel = lipgloss.NewStyle().
		Background(ColourBackground).
		Foreground(ColourForeground).
		Align(lipgloss.Center).
		Width(PanelWidth)
)
