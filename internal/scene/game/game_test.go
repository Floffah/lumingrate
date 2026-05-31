package gamescene

import (
	"reflect"
	"strings"
	"testing"

	"luminrate/internal/engine"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestEntryFromEventMapsRawEvents(t *testing.T) {
	entry := entryFromEvent(engine.Event{Kind: engine.EventRaw, Text: "plain"})

	if entry.kind != entryRaw {
		t.Fatalf("entry.kind = %v, want %v", entry.kind, entryRaw)
	}
	if entry.text != "plain" {
		t.Fatalf("entry.text = %q, want %q", entry.text, "plain")
	}
}

func TestRenderEntryLeavesRawTextUnstyled(t *testing.T) {
	model := Model{}
	rendered := visibleText(model.renderEntry(entry{kind: entryRaw, text: "plain text"}, 20))

	if rendered != "plain text" {
		t.Fatalf("rendered = %q, want %q", rendered, "plain text")
	}
}

func TestRenderLogLinesJoinsRawEntries(t *testing.T) {
	model := Model{
		log: []entry{
			{kind: entryRaw, text: "load"},
			{kind: entryRaw, text: "ing"},
			{kind: entryRaw, text: "..."},
		},
	}

	lines := model.renderLogLines(80)
	lines = visibleLines(lines)
	want := []string{"loading..."}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %#v, want %#v", lines, want)
	}
}

func TestRenderLogLinesUsesRawNewlines(t *testing.T) {
	model := Model{
		log: []entry{
			{kind: entryRaw, text: "one\n"},
			{kind: entryRaw, text: "two\nthree"},
		},
	}

	lines := model.renderLogLines(80)
	lines = visibleLines(lines)
	want := []string{"one", "two", "three"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %#v, want %#v", lines, want)
	}
}

func TestRenderInputUsesStyledCursorWhenVisible(t *testing.T) {
	model := Model{
		input:         "look",
		cursor:        2,
		cursorVisible: true,
	}

	rendered := model.renderInput()
	want := "lo" + cursorStyle.Render("o") + "k"
	if rendered != want {
		t.Fatalf("rendered = %q, want %q", rendered, want)
	}
}

func TestRenderInputUsesPlainCursorCellWhenHidden(t *testing.T) {
	model := Model{
		input:         "look",
		cursor:        2,
		cursorVisible: false,
	}

	rendered := model.renderInput()
	want := "look"
	if rendered != want {
		t.Fatalf("rendered = %q, want %q", rendered, want)
	}
}

func TestCursorBlinkTogglesWhileInputVisible(t *testing.T) {
	model := Model{
		inputVisible:  true,
		cursorVisible: true,
	}

	next, _ := model.Update(cursorBlinkMsg{})
	updated := next.(*Model)

	if updated.cursorVisible {
		t.Fatal("cursorVisible = true, want false")
	}
}

func visibleLines(lines []string) []string {
	visible := make([]string, len(lines))
	for i, line := range lines {
		visible[i] = visibleText(line)
	}
	return visible
}

func visibleText(text string) string {
	lines := strings.Split(ansi.Strip(text), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

func TestCursorBlinkStaysHiddenWhileInputHidden(t *testing.T) {
	model := Model{
		inputVisible:  false,
		cursorVisible: true,
	}

	next, _ := model.Update(cursorBlinkMsg{})
	updated := next.(*Model)

	if updated.cursorVisible {
		t.Fatal("cursorVisible = true, want false")
	}
}

func TestSpacePressWhileInputHiddenEnablesFastForward(t *testing.T) {
	tests := []struct {
		name string
		key  tea.Key
	}{
		{name: "code", key: tea.Key{Code: tea.KeySpace}},
		{name: "base code", key: tea.Key{BaseCode: tea.KeySpace}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameEngine := engine.NewEngine(nil)
			model := Model{
				engine:       gameEngine,
				inputVisible: false,
			}

			next, _ := model.Update(tea.KeyPressMsg(tt.key))
			updated := next.(*Model)

			if !updated.skipHeld {
				t.Fatal("skipHeld = false, want true")
			}
			if !gameEngine.IsFastForward() {
				t.Fatal("IsFastForward() = false, want true")
			}
		})
	}
}

func TestSpaceReleaseClearsFastForward(t *testing.T) {
	gameEngine := engine.NewEngine(nil)
	gameEngine.SetFastForward(true)
	model := Model{
		engine:       gameEngine,
		inputVisible: false,
		skipHeld:     true,
	}

	next, _ := model.Update(tea.KeyReleaseMsg(tea.Key{Code: tea.KeySpace}))
	updated := next.(*Model)

	if updated.skipHeld {
		t.Fatal("skipHeld = true, want false")
	}
	if gameEngine.IsFastForward() {
		t.Fatal("IsFastForward() = true, want false")
	}
}

func TestSpacePressWhileInputVisibleInsertsSpace(t *testing.T) {
	gameEngine := engine.NewEngine(nil)
	model := Model{
		engine:       gameEngine,
		inputVisible: true,
	}

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeySpace, Text: " "}))
	updated := next.(*Model)

	if updated.skipHeld {
		t.Fatal("skipHeld = true, want false")
	}
	if gameEngine.IsFastForward() {
		t.Fatal("IsFastForward() = true, want false")
	}
	if updated.input != " " {
		t.Fatalf("input = %q, want space", updated.input)
	}
}

func TestPromptVisibilityResyncsFastForward(t *testing.T) {
	gameEngine := engine.NewEngine(nil)
	model := Model{
		engine:       gameEngine,
		inputVisible: true,
		skipHeld:     true,
	}

	model.applyEvent(engine.Event{Kind: engine.EventHideInput})
	if !gameEngine.IsFastForward() {
		t.Fatal("IsFastForward() = false after hide, want true")
	}

	model.applyEvent(engine.Event{Kind: engine.EventShowInput})
	if model.skipHeld {
		t.Fatal("skipHeld = true after show, want false")
	}
	if gameEngine.IsFastForward() {
		t.Fatal("IsFastForward() = true after show, want false")
	}
}

func TestViewRequestsKeyReleaseEvents(t *testing.T) {
	model := Model{width: defaultWidth, height: defaultHeight}
	view := model.View()

	if !view.KeyboardEnhancements.ReportEventTypes {
		t.Fatal("ReportEventTypes = false, want true")
	}
}
