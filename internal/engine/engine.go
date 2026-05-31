package engine

import (
	"context"
	"sync"
	"time"
)

type EventKind int

const (
	EventNarration EventKind = iota
	EventAside
	EventInputMessage
	EventHideInput
	EventShowInput
	EventInsertInput
	EventRaw
)

type Event struct {
	Kind EventKind
	Text string
}

type Command struct {
	Kind CommandKind
	Text string
}

type CommandKind int

const (
	CommandSubmit CommandKind = iota
	CommandTab
)

type Runtime interface {
	Sleep(time.Duration)
	Emit(Event) bool
	EmitNarration(string)
	EmitAside(string)
	IsFastForward() bool
	ShowInput()
	HideInput()
	ShowInputMessage(string)
	EmitCommandsHint(string, []string)
	EmitTypingText(string)
	BlockUntilContinue() string
	BlockUntilContinueWithText(string) string
	ConsumeCommandsNow()
	ConsumeCommandsInline() []Command
}

type Chapter interface {
	Begin(Runtime)
	HandleCommand(Runtime, string)
	HandleTab(Runtime, string)
}

type Engine struct {
	ctx             context.Context
	commands        chan Command
	events          chan Event
	chapter         Chapter
	fastForwardMu   sync.Mutex
	fastForward     bool
	fastForwardWake chan struct{}
}

func NewEngine(chapter Chapter) *Engine {
	return &Engine{
		commands:        make(chan Command),
		events:          make(chan Event, 16),
		chapter:         chapter,
		fastForwardWake: make(chan struct{}),
	}
}

func (e *Engine) Commands() chan<- Command {
	return e.commands
}

func (e *Engine) Events() <-chan Event {
	return e.events
}

func (e *Engine) SetFastForward(enabled bool) {
	e.fastForwardMu.Lock()
	defer e.fastForwardMu.Unlock()

	if e.fastForward == enabled {
		return
	}

	e.fastForward = enabled
	if enabled {
		close(e.fastForwardWake)
		return
	}

	e.fastForwardWake = make(chan struct{})
}

func (e *Engine) IsFastForward() bool {
	e.fastForwardMu.Lock()
	defer e.fastForwardMu.Unlock()

	return e.fastForward
}

func (e *Engine) Run(ctx context.Context) {
	defer close(e.events)

	e.ctx = ctx
	if e.chapter != nil {
		e.chapter.Begin(e)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case command := <-e.commands:
			e.handleCommand(command)
		}
	}
}

func (e *Engine) handleCommand(command Command) {
	switch command.Kind {
	case CommandTab:
		e.handleTab(command.Text)
	default:
		e.handleSubmit(command.Text)
	}
}

func (e *Engine) handleSubmit(command string) {
	if e.chapter == nil {
		e.EmitAside("There is no chapter loaded.")
		return
	}

	e.chapter.HandleCommand(e, NormaliseCommand(command))
}

func (e *Engine) handleTab(input string) {
	if e.chapter == nil {
		e.ShowInputMessage("There is no chapter loaded.")
		return
	}

	e.chapter.HandleTab(e, NormaliseCommand(input))
}

func matchingCommands(commands []string, prefix string) []string {
	matches := make([]string, 0, len(commands))
	for _, command := range commands {
		normalised := NormaliseCommand(command)
		if prefix == "" || len(normalised) >= len(prefix) && normalised[:len(prefix)] == prefix {
			matches = append(matches, normalised)
		}
	}
	return matches
}

func (e *Engine) Emit(event Event) bool {
	select {
	case <-e.ctx.Done():
		return false
	case e.events <- event:
		return true
	}
}
