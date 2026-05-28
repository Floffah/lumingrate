package gamescene

import (
	"reflect"
	"testing"

	"luminrate/internal/engine"
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
	rendered := (Model{}).renderEntry(entry{kind: entryRaw, text: "plain text"}, 5)

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
	want := []string{"one", "two", "three"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %#v, want %#v", lines, want)
	}
}
