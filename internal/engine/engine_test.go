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
	expectEventKind(t, gameEngine.Events(), EventShowInput)

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

func expectEventKind(t *testing.T, events <-chan Event, want EventKind) {
	t.Helper()

	select {
	case event := <-events:
		if event.Kind != want {
			t.Fatalf("event.Kind = %v, want %v", event.Kind, want)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for event kind %v", want)
	}
}

func TestSleepReturnsAfterDuration(t *testing.T) {
	gameEngine := NewEngine(nil)
	done := make(chan struct{})

	go func() {
		gameEngine.Sleep(20 * time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sleep did not return after duration")
	}
}

func TestSleepReturnsWhenFastForwardEnabled(t *testing.T) {
	gameEngine := NewEngine(nil)
	done := make(chan struct{})

	go func() {
		gameEngine.Sleep(time.Hour)
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	gameEngine.SetFastForward(true)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sleep did not return after fast-forward was enabled")
	}
}

func TestEmitNarrationAndWaitEmitsNarration(t *testing.T) {
	gameEngine := NewEngine(nil)
	gameEngine.ctx = context.Background()

	gameEngine.EmitNarrationAndWait("")

	select {
	case event := <-gameEngine.Events():
		if event.Kind != EventNarration {
			t.Fatalf("event.Kind = %v, want %v", event.Kind, EventNarration)
		}
		if event.Text != "" {
			t.Fatalf("event.Text = %q, want empty text", event.Text)
		}
	default:
		t.Fatal("EmitNarrationAndWait did not emit an event")
	}
}

func TestNarrationWaitDurationUsesReadingSpeedFormula(t *testing.T) {
	text := make([]byte, averageNarrationWordLength*averageNarrationWordsPerMinute)

	if got := narrationWaitDuration(string(text)); got != time.Minute {
		t.Fatalf("narrationWaitDuration() = %v, want %v", got, time.Minute)
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
