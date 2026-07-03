package story

import (
	"luminrate/internal/chapter"
	"luminrate/internal/engine"
)

const chapterSeparator = "\n---\n\n"

type chapterFactory func(engine.ChapterID) (engine.Chapter, bool)

type Story struct {
	initial engine.ChapterID
	factory chapterFactory

	currentID engine.ChapterID
	current   engine.Chapter
}

func New() *Story {
	return newWithFactory(chapter.InitialID, chapter.New)
}

func NewWithInitialChapter(initial engine.ChapterID) (*Story, bool) {
	if _, ok := chapter.New(initial); !ok {
		return nil, false
	}

	return newWithFactory(initial, chapter.New), true
}

func newWithFactory(initial engine.ChapterID, factory chapterFactory) *Story {
	return &Story{
		initial: initial,
		factory: factory,
	}
}

func NewEngine() *engine.Engine {
	return engine.NewEngine(New())
}

func NewEngineWithInitialChapter(initial engine.ChapterID) (*engine.Engine, bool) {
	story, ok := NewWithInitialChapter(initial)
	if !ok {
		return nil, false
	}

	return engine.NewEngine(story), true
}

func (s *Story) Begin(runtime engine.Runtime) {
	s.switchChapter(runtime, s.initial)
}

func (s *Story) HandleCommand(runtime engine.Runtime, command string) {
	if s.current == nil {
		s.switchChapter(runtime, s.initial)
	}
	if s.current == nil {
		return
	}

	s.current.HandleCommand(s.wrapRuntime(runtime), command)
}

func (s *Story) HandleTab(runtime engine.Runtime, input string) {
	if s.current == nil {
		s.switchChapter(runtime, s.initial)
	}
	if s.current == nil {
		return
	}

	s.current.HandleTab(s.wrapRuntime(runtime), input)
}

func (s *Story) switchChapter(runtime engine.Runtime, id engine.ChapterID) {
	next, ok := s.factory(id)
	if !ok {
		runtime.EmitAside("Chapter " + string(id) + " is not available.")
		return
	}

	if s.current != nil {
		runtime.Emit(engine.Event{Kind: engine.EventRaw, Text: chapterSeparator})
	}

	s.currentID = id
	s.current = next
	s.current.Begin(s.wrapRuntime(runtime))
}

func (s *Story) wrapRuntime(runtime engine.Runtime) engine.Runtime {
	return storyRuntime{
		Runtime: runtime,
		story:   s,
	}
}

type storyRuntime struct {
	engine.Runtime
	story *Story
}

func (r storyRuntime) SwitchChapter(id engine.ChapterID) {
	r.story.switchChapter(r.Runtime, id)
}
