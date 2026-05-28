package mainmenu

import "strings"

func (m model) renderSettings() string {
	return panelStyle.Render(strings.Join([]string{
		"SETTINGS",
		"",
		"No settings are available yet.",
		"",
		"Press Enter to return.",
	}, "\n"))
}
