package chapter

import "luminrate/internal/engine"

const (
	PrologueID engine.ChapterID = "prologue"
	TheCardID  engine.ChapterID = "the_card"
)

const InitialID = PrologueID

func New(id engine.ChapterID) (engine.Chapter, bool) {
	switch id {
	case PrologueID:
		return NewPrologue(), true
	case TheCardID:
		return NewTheCard(), true
	default:
		return nil, false
	}
}

func NewInitialChapter() engine.Chapter {
	chapter, _ := New(InitialID)
	return chapter
}
