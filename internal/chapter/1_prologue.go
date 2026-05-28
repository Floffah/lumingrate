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
)

// Prologue Chapter 1
type Prologue struct {
	phase        prologuePhase
	gettingUpCan context.CancelFunc
	toldToGetUp  bool
}

func NewAwaken() *Prologue {
	return &Prologue{phase: prologuePhaseInitial}
}

func (c *Prologue) Begin(runtime engine.Runtime) {
	runtime.Emit(engine.Event{Kind: engine.EventAside, Text: "Tip: if you don't know what to type, pressing tab will give you some hints."})
	runtime.Emit(engine.Event{Kind: engine.EventNarration})
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
		ticker := time.NewTicker(3 * time.Second)
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

func (c *Prologue) HandleCommand(runtime engine.Runtime, cmd string) {
	switch c.phase {
	case prologuePhaseInitial:
		if strings.HasPrefix(cmd, "get") {
			c.gettingUpCan()
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You get up and walk to the door."})
			time.Sleep(2 * time.Second)
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "Three more knocks at the door"})
			c.phase = prologuePhaseGotUp
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
		} else {
			runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "Get up you lazy sod!"})
			c.toldToGetUp = true
		}
	}
}
