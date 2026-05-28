package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"luminrate/internal/scene/app"
)

func main() {
	if _, err := tea.NewProgram(app.NewModel()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "lumingrate: %v\n", err)
		os.Exit(1)
	}
}
