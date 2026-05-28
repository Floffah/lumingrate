package engine

import (
	"context"
	"testing"
	"time"
)

type blockingChapter struct {
	resumed chan struct{}
}

func (c blockingChapter) Begin(runtime Runtime) {
	runtime.BlockUntilContinue()
	close(c.resumed)
}

func (c blockingChapter) HandleCommand(Runtime, string) {}

func (c blockingChapter) HandleTab(Runtime, string) {}

func TestBlockUntilContinueResumesOnEmptySubmit(t *testing.T) {
	chapter := blockingChapter{resumed: make(chan struct{})}
	gameEngine := NewEngine(chapter)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gameEngine.Run(ctx)

	select {
	case event := <-gameEngine.Events():
		if event.Kind != EventInputMessage {
			t.Fatalf("event.Kind = %v, want %v", event.Kind, EventInputMessage)
		}
		if event.Text != "Press enter to continue..." {
			t.Fatalf("event.Text = %q, want %q", event.Text, "Press enter to continue...")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for continue prompt")
	}

	select {
	case gameEngine.Commands() <- Command{Kind: CommandSubmit}:
	case <-time.After(time.Second):
		t.Fatal("timed out sending continue command")
	}

	select {
	case <-chapter.resumed:
	case <-time.After(time.Second):
		t.Fatal("chapter did not resume after continue command")
	}
}

func TestBlockUntilContinueRepromptsUntilEmptySubmit(t *testing.T) {
	chapter := blockingChapter{resumed: make(chan struct{})}
	gameEngine := NewEngine(chapter)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gameEngine.Run(ctx)

	expectContinuePrompt(t, gameEngine.Events())

	select {
	case gameEngine.Commands() <- Command{Kind: CommandSubmit, Text: "wait"}:
	case <-time.After(time.Second):
		t.Fatal("timed out sending non-empty command")
	}

	expectContinuePrompt(t, gameEngine.Events())

	select {
	case <-chapter.resumed:
		t.Fatal("chapter resumed before empty submit")
	default:
	}

	select {
	case gameEngine.Commands() <- Command{Kind: CommandSubmit}:
	case <-time.After(time.Second):
		t.Fatal("timed out sending continue command")
	}

	select {
	case <-chapter.resumed:
	case <-time.After(time.Second):
		t.Fatal("chapter did not resume after continue command")
	}
}

func expectContinuePrompt(t *testing.T, events <-chan Event) {
	t.Helper()

	select {
	case event := <-events:
		if event.Kind != EventInputMessage {
			t.Fatalf("event.Kind = %v, want %v", event.Kind, EventInputMessage)
		}
		if event.Text != "Press enter to continue..." {
			t.Fatalf("event.Text = %q, want %q", event.Text, "Press enter to continue...")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for continue prompt")
	}
}
