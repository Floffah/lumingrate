package chapter

import (
	"context"
	"luminrate/internal/engine"
	"luminrate/internal/style"
	"strings"
	"time"
)

type prologuePhase int

const (
	prologuePhaseInitial prologuePhase = iota
	prologuePhaseGotUp
	prologuePhaseOpenedDoor
)

// Prologue Chapter 1
type Prologue struct {
	phase        prologuePhase
	gettingUpCan context.CancelFunc

	// prologuePhaseInitial
	toldToGetUp bool

	// prologuePhaseOpenedDoor
}

func NewAwaken() *Prologue {
	return &Prologue{phase: prologuePhaseInitial}
}

func (c *Prologue) Begin(runtime engine.Runtime) {
	runtime.Emit(engine.Event{Kind: engine.EventAside, Text: "Some tips:"})
	runtime.Emit(engine.Event{Kind: engine.EventAside, Text: " - If you don't know what to type, pressing tab will give you some hints."})
	runtime.Emit(engine.Event{Kind: engine.EventAside, Text: " - There is no AI here, so we won't be able to understand you if you type something unexpected, but the hints should help you figure out what to do."})
	runtime.Emit(engine.Event{Kind: engine.EventAside, Text: " - The prompt will only be hidden when there's nothing to do. But when it's there, that doesn't mean the narration is waiting for you to type something, it means you can interrupt."})
	runtime.Emit(engine.Event{Kind: engine.EventNarration})

	c.BeginPhaseInitial(runtime)
}

func (c *Prologue) HandleCommand(runtime engine.Runtime, cmd string) {
	switch c.phase {
	case prologuePhaseInitial:
		if strings.Contains(cmd, "get") {
			c.BeginPhaseGotUp(runtime)
			return
		}
	}

	runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "No such action"})
}

func (c *Prologue) HandleTab(runtime engine.Runtime, input string) {
	switch c.phase {
	case prologuePhaseInitial:
		if c.toldToGetUp {
			runtime.Emit(engine.Event{Kind: engine.EventInsertInput, Text: "get up"})
			runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "I have to do everything myself..."})
			return
		} else {
			runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: style.Aside.Bold(true).Render("Get") + style.Aside.Render(" up you lazy sod!")})
			c.toldToGetUp = true
			return
		}
	}

	runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "No hint available :)"})
}

func (c *Prologue) BeginPhaseInitial(runtime engine.Runtime) {
	c.phase = prologuePhaseInitial

	runtime.Emit(engine.Event{Kind: engine.EventHideInput})
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You awaken in your usual spot on the couch to three knocks at the door. The TV still on and empty pizza boxes on the coffee table."})
	time.Sleep(2 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: style.Text.Render("Looking around, ") + style.TextItalic.Render("Holy shit I need a job") + style.Text.Render(", you think to yourself.")})
	time.Sleep(2 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: style.TextItalic.Render("Who the hell could that be?")})
	runtime.Emit(engine.Event{Kind: engine.EventShowInput})

	getUpCtx, getUpCan := context.WithCancel(context.Background())
	c.gettingUpCan = getUpCan

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			case <-getUpCtx.Done():
				return
			}
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "Three more knocks at the door"})
		}
	}()
}

func (c *Prologue) BeginPhaseGotUp(runtime engine.Runtime) {
	c.gettingUpCan()
	c.phase = prologuePhaseGotUp

	runtime.Emit(engine.Event{Kind: engine.EventHideInput})

	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You get up and walk to the door."})
	time.Sleep(2 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "Three more knocks at the door"})
	time.Sleep(1 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: style.TextItalic.Render("Yeah yeah I'm coming...")})
	time.Sleep(2 * time.Second)

	runtime.Emit(engine.Event{Kind: engine.EventNarration})
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The man at the door stands straight, though not rigidly. Dark coat. No visible insignia beyond a small silver mark pinned near the collar, circular, crossed by a single vertical line. Military, perhaps. Or adjacent to it."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "Rain beads across the shoulders of the coat."})
	//time.Sleep(3 * time.Second)

	runtime.BlockUntilContinue()
	c.BeginPhaseOpenedDoor(runtime)
}

func (c *Prologue) BeginPhaseOpenedDoor(runtime engine.Runtime) {
	c.phase = prologuePhaseOpenedDoor

	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"Mr Paints?\""})
	time.Sleep(2 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "His voice is calm, practiced."})
	time.Sleep(4 * time.Second)

	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"You don't know who I am. That's expected.\""})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"My name is Elias Vale. I’m part of an outreach delegation representing an independent human civilisation currently operating beyond Earth jurisdiction.\""})

	time.Sleep(7 * time.Second)
}
