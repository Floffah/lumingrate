package engine

import (
	"context"
	"luminrate/internal/util"
)

type EventKind int

const (
	EventNarration EventKind = iota
	EventAside
	EventInputMessage
	EventHideInput
	EventShowInput
	EventInsertInput
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
	Emit(Event) bool
	EmitCommandsHint(string, []string)
}

type Chapter interface {
	Begin(Runtime)
	HandleCommand(Runtime, string)
	HandleTab(Runtime, string)
}

type Engine struct {
	ctx      context.Context
	commands chan Command
	events   chan Event
	chapter  Chapter
}

func NewEngine(chapter Chapter) *Engine {
	return &Engine{
		commands: make(chan Command),
		events:   make(chan Event, 16),
		chapter:  chapter,
	}
}

func (e *Engine) Commands() chan<- Command {
	return e.commands
}

func (e *Engine) Events() <-chan Event {
	return e.events
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
		e.Emit(Event{Kind: EventAside, Text: "There is no chapter loaded."})
		return
	}

	e.chapter.HandleCommand(e, util.NormaliseCommand(command))
}

func (e *Engine) handleTab(input string) {
	if e.chapter == nil {
		e.Emit(Event{Kind: EventInputMessage, Text: "There is no chapter loaded."})
		return
	}

	e.chapter.HandleTab(e, util.NormaliseCommand(input))
}

func matchingCommands(commands []string, prefix string) []string {
	matches := make([]string, 0, len(commands))
	for _, command := range commands {
		normalised := util.NormaliseCommand(command)
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

func (e *Engine) EmitCommandsHint(input string, commands []string) {
	command := util.NormaliseCommand(input)
	matches := matchingCommands(commands, command)
	switch len(matches) {
	case 0:
		e.Emit(Event{Kind: EventInputMessage, Text: "No hints available :)"})
	case 1:
		e.Emit(Event{Kind: EventInsertInput, Text: matches[0][len(command):]})
	default:
		e.Emit(Event{Kind: EventInputMessage, Text: "Hint: " + util.JoinForDisplay(matches) + ""})
	}
}
