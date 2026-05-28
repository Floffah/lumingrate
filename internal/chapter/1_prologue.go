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
	prologuePhaseIntroduced
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
	case prologuePhaseGotUp:
		if strings.Contains(cmd, "greet") || strings.Contains(cmd, "hello") || strings.Contains(cmd, "hi") {
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"Hello, do I know you?\""})
			time.Sleep(3 * time.Second)
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The man accepts the greeting with the smallest nod, but continues as if he had expected exactly that."})
			time.Sleep(3 * time.Second)
			c.BeginPhaseOpenedDoor(runtime)
			return
		} else if strings.TrimSpace(cmd) == "" || strings.Contains(cmd, "stay") || strings.Contains(cmd, "silent") {
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You stay silent."})
			time.Sleep(3 * time.Second)
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The man doesn't seem to acknowledge your silence."})
			time.Sleep(3 * time.Second)
			c.BeginPhaseOpenedDoor(runtime)
			return
		}
	case prologuePhaseOpenedDoor:
		if strings.Contains(cmd, "question") || strings.Contains(cmd, "who") || strings.Contains(cmd, "what") {
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You ask him who he is and what he wants. You ask him who he is and what he wants, because people in dark coats with custom lapel pins rarely sell broadband."})
			time.Sleep(3 * time.Second)
			c.BeginPhaseIntroduced(runtime)
			return
		} else if strings.TrimSpace(cmd) == "" || strings.Contains(cmd, "stay") || strings.Contains(cmd, "silent") {
			runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "You stay silent."})
			time.Sleep(3 * time.Second)
			c.BeginPhaseIntroduced(runtime)
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
			runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: style.AsideKeyword.Render("Get") + style.Aside.Render(" up you lazy sod!")})
			c.toldToGetUp = true
			return
		}
	case prologuePhaseGotUp:
		runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "You could " + style.AsideKeyword.Render("stay silent") + style.Aside.Render(" or ") + style.AsideKeyword.Render("greet") + style.Aside.Render(" him")})
		return
	case prologuePhaseOpenedDoor:
		runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "You could " + style.AsideKeyword.Render("stay silent") + style.Aside.Render(" or ") + style.AsideKeyword.Render("question") + style.Aside.Render(" him")})
		return
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

	runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "Hit enter or type to continue..."})
	runtime.Emit(engine.Event{Kind: engine.EventShowInput})
}

func (c *Prologue) BeginPhaseOpenedDoor(runtime engine.Runtime) {
	c.phase = prologuePhaseOpenedDoor

	runtime.Emit(engine.Event{Kind: engine.EventHideInput})
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"Mr Paints?\""})
	time.Sleep(2 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "His voice is calm, practiced."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"My name is Elias Vale. I represent the Luminous Outreach Department.\""})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration})
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "There is no grin, no theatricality, no attempt to sell his statement."})

	runtime.Emit(engine.Event{Kind: engine.EventInputMessage, Text: "Hit enter or type to continue..."})
	runtime.Emit(engine.Event{Kind: engine.EventShowInput})
}

func (c *Prologue) BeginPhaseIntroduced(runtime engine.Runtime) {
	c.phase = prologuePhaseIntroduced

	runtime.Emit(engine.Event{Kind: engine.EventHideInput})
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"You are under no obligation to continue this conversation. You are not in danger. And you are not being accused of anything.\""})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The man’s eyes are steady, but not unkind. He seems to be giving you the benefit of the doubt, as if he’s not sure whether you’re going to be cooperative or not."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The rain taps softly against the path outside."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"We were given your name by someone familiar with your work in astrodynamics. Your paper on some unexplained noise in low-orbit survey data reached someone in our office. A lot of people file these irregularities as sensor issues and noise. However, you did not.\""})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "He reaches into his coat slowly enough to avoid alarm, producing a slim matte card."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "The card has no government seal nor any indication of its origin or national flag, but the same circular silver mark is embossed in the corner."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"There are things about the present state of humanity that you have not been told yet. We believe you are someone capable of understanding them responsibly.\""})

	runtime.BlockUntilContinue()
	runtime.Emit(engine.Event{Kind: engine.EventHideInput})

	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "His eyes move briefly past you towards the dark house behind."})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"We’re not here because you’re special, Mr Paints. We’re here because you’re competent.\""})
	time.Sleep(5 * time.Second)
	runtime.Emit(engine.Event{Kind: engine.EventNarration, Text: "\"And because, at some point in your life, you continued spending your evenings looking at things that others wouldn’t even know are there.\""})
}
