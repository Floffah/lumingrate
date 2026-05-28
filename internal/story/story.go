package story

import (
	"luminrate/internal/chapter"
	"luminrate/internal/engine"
)

func NewEngine() *engine.Engine {
	return engine.NewEngine(chapter.NewInitialChapter())
}
