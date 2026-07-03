package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"luminrate/internal/engine"
	"luminrate/internal/scene/app"
)

func main() {
	chapterFlag := flag.String("chapter", "", "start directly at a chapter ID")
	flag.Parse()

	model := app.NewModel()
	if *chapterFlag != "" {
		var ok bool
		model, ok = app.NewModelWithInitialChapter(engine.ChapterID(*chapterFlag))
		if !ok {
			fmt.Fprintf(os.Stderr, "lumingrate: unknown chapter %q\n", *chapterFlag)
			os.Exit(2)
		}
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "lumingrate: %v\n", err)
		os.Exit(1)
	}
}
