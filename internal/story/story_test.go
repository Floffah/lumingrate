package story

import (
	"testing"
	"time"

	"luminrate/internal/chapter"
	"luminrate/internal/engine"
)

const (
	testChapterOne engine.ChapterID = "test_one"
	testChapterTwo engine.ChapterID = "test_two"
	testMissing    engine.ChapterID = "missing"
)

type fakeChapter struct {
	id    engine.ChapterID
	began *[]engine.ChapterID
}

func (c fakeChapter) Begin(engine.Runtime) {
	*c.began = append(*c.began, c.id)
}

func (c fakeChapter) HandleCommand(runtime engine.Runtime, command string) {
	if command == "next" {
		runtime.SwitchChapter(testChapterTwo)
	}
	if command == "missing" {
		runtime.SwitchChapter(testMissing)
	}
}

func (c fakeChapter) HandleTab(engine.Runtime, string) {}

func TestStoryBeginsInitialChapter(t *testing.T) {
	var began []engine.ChapterID
	story := newWithFactory(testChapterOne, testFactory(&began))
	runtime := &noopRuntime{}

	story.Begin(runtime)

	if story.currentID != testChapterOne {
		t.Fatalf("currentID = %q, want %q", story.currentID, testChapterOne)
	}
	if len(began) != 1 || began[0] != testChapterOne {
		t.Fatalf("began = %#v, want [%q]", began, testChapterOne)
	}
	if len(runtime.events) != 0 {
		t.Fatalf("events = %#v, want no initial chapter separator", runtime.events)
	}
}

func TestNewWithInitialChapterUsesRequestedChapter(t *testing.T) {
	story, ok := NewWithInitialChapter(chapter.TheCardID)
	if !ok {
		t.Fatal("NewWithInitialChapter returned ok = false, want true")
	}
	if story.initial != chapter.TheCardID {
		t.Fatalf("initial = %q, want %q", story.initial, chapter.TheCardID)
	}
}

func TestNewWithInitialChapterRejectsMissingChapter(t *testing.T) {
	story, ok := NewWithInitialChapter(testMissing)
	if ok {
		t.Fatal("NewWithInitialChapter returned ok = true, want false")
	}
	if story != nil {
		t.Fatalf("story = %#v, want nil", story)
	}
}

func TestStorySwitchesChapterFromRuntime(t *testing.T) {
	var began []engine.ChapterID
	story := newWithFactory(testChapterOne, testFactory(&began))
	runtime := &noopRuntime{}
	story.Begin(runtime)

	story.HandleCommand(runtime, "next")

	if story.currentID != testChapterTwo {
		t.Fatalf("currentID = %q, want %q", story.currentID, testChapterTwo)
	}
	if len(runtime.events) != 1 {
		t.Fatalf("events = %#v, want one chapter separator event", runtime.events)
	}
	if runtime.events[0].Kind != engine.EventRaw || runtime.events[0].Text != chapterSeparator {
		t.Fatalf("separator event = %#v, want raw %q", runtime.events[0], chapterSeparator)
	}
	want := []engine.ChapterID{testChapterOne, testChapterTwo}
	if len(began) != len(want) {
		t.Fatalf("began = %#v, want %#v", began, want)
	}
	for i := range want {
		if began[i] != want[i] {
			t.Fatalf("began = %#v, want %#v", began, want)
		}
	}
}

func TestStoryKeepsCurrentChapterWhenSwitchTargetIsMissing(t *testing.T) {
	var began []engine.ChapterID
	runtime := &noopRuntime{}
	story := newWithFactory(testChapterOne, testFactory(&began))
	story.Begin(runtime)

	story.HandleCommand(runtime, "missing")

	if story.currentID != testChapterOne {
		t.Fatalf("currentID = %q, want %q", story.currentID, testChapterOne)
	}
	if len(runtime.asides) != 1 {
		t.Fatalf("asides = %#v, want one unavailable message", runtime.asides)
	}
	if len(runtime.events) != 0 {
		t.Fatalf("events = %#v, want no separator for missing chapter", runtime.events)
	}
}

func testFactory(began *[]engine.ChapterID) chapterFactory {
	return func(id engine.ChapterID) (engine.Chapter, bool) {
		switch id {
		case testChapterOne, testChapterTwo:
			return fakeChapter{id: id, began: began}, true
		default:
			return nil, false
		}
	}
}

type noopRuntime struct {
	asides []string
	events []engine.Event
}

func (r *noopRuntime) Sleep(time.Duration)                      {}
func (r *noopRuntime) Emit(event engine.Event) bool             { r.events = append(r.events, event); return true }
func (r *noopRuntime) EmitNarration(string)                     {}
func (r *noopRuntime) EmitNarrationAndWait(string)              {}
func (r *noopRuntime) EmitAside(text string)                    { r.asides = append(r.asides, text) }
func (r *noopRuntime) IsFastForward() bool                      { return false }
func (r *noopRuntime) ShowInput()                               {}
func (r *noopRuntime) HideInput()                               {}
func (r *noopRuntime) ShowInputMessage(string)                  {}
func (r *noopRuntime) EmitCommandsHint(string, []string)        {}
func (r *noopRuntime) EmitTypingText(string)                    {}
func (r *noopRuntime) BlockUntilContinue() string               { return "" }
func (r *noopRuntime) BlockUntilContinueWithText(string) string { return "" }
func (r *noopRuntime) ConsumeCommandsNow()                      {}
func (r *noopRuntime) ConsumeCommandsInline() []engine.Command  { return nil }
func (r *noopRuntime) SwitchChapter(engine.ChapterID)           {}
