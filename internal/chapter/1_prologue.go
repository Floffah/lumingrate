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
	prologuePhaseCredibility
)

// Prologue Chapter 1
type Prologue struct {
	phase        prologuePhase
	gettingUpCan context.CancelFunc

	// prologuePhaseInitial
	toldToGetUp bool

	// prologuePhaseIntroduced
	finishedIntroductions bool

	// prologuePhaseCredibility
	didSlamDoor bool
	didTakeCard bool
}

func NewPrologue() *Prologue {
	return &Prologue{phase: prologuePhaseInitial}
}

func (c *Prologue) Begin(runtime engine.Runtime) {
	runtime.EmitAside("Some tips:")
	runtime.EmitAside(" - If you don't know what to type, pressing tab will give you some hints.")
	runtime.EmitAside(" - There is no AI here, so we won't be able to understand you if you type something unexpected, but the hints should help you figure out what to do.")
	runtime.EmitAside(" - The prompt will only be hidden when there's nothing to do. But when it's there, that doesn't mean the narration is waiting for you to type something, it means you can interrupt.")
	runtime.EmitNarration("")

	c.BeginPhaseInitial(runtime)
}

func (c *Prologue) HandleCommand(runtime engine.Runtime, cmd string) {
	switch c.phase {
	case prologuePhaseInitial:
		if strings.Contains(cmd, "get") {
			// label: get up
			c.BeginPhaseGotUp(runtime)
			return
		}
	case prologuePhaseGotUp:
		if strings.Contains(cmd, "greet") || strings.Contains(cmd, "hello") || strings.Contains(cmd, "hi") {
			// label: greet Elias
			runtime.HideInput()
			runtime.EmitNarration("\"Hello, do I know you?\"")
			runtime.Sleep(3 * time.Second)
			runtime.EmitNarration("The man accepts the greeting with the smallest nod, but continues as if he had expected exactly that.")
			runtime.Sleep(3 * time.Second)
			c.BeginPhaseOpenedDoor(runtime)
			return
		} else if strings.TrimSpace(cmd) == "" || strings.Contains(cmd, "stay") || strings.Contains(cmd, "silent") {
			// label: stay silent
			runtime.HideInput()
			runtime.EmitNarration("You stay silent.")
			runtime.Sleep(3 * time.Second)
			runtime.EmitNarration("The man doesn't seem to acknowledge your silence.")
			runtime.Sleep(3 * time.Second)
			c.BeginPhaseOpenedDoor(runtime)
			return
		}
	case prologuePhaseOpenedDoor:
		if strings.Contains(cmd, "question") || strings.Contains(cmd, "who") || strings.Contains(cmd, "what") {
			// label: ask who Elias is
			runtime.HideInput()
			runtime.EmitNarration("You ask him who he is and what he wants, because people in dark coats with custom lapel pins rarely sell broadband.")
			runtime.Sleep(3 * time.Second)
			c.BeginPhaseIntroduced(runtime)
			return
		} else if strings.TrimSpace(cmd) == "" || strings.Contains(cmd, "stay") || strings.Contains(cmd, "silent") {
			// label: let Elias continue
			runtime.HideInput()
			runtime.EmitNarration("You stay silent.")
			runtime.Sleep(3 * time.Second)
			c.BeginPhaseIntroduced(runtime)
			return
		}
	case prologuePhaseIntroduced:
		if strings.Contains(cmd, "question") || strings.Contains(cmd, "proof") || strings.Contains(cmd, "what") {
			// label: ask for proof
			runtime.HideInput()
			runtime.EmitNarration("You ask him what proof he has that he's not just some weirdo who showed up at your door. Maybe his Wifi is \"out of this world!\".")
			runtime.Sleep(5 * time.Second)
			runtime.EmitNarration("He doesn't seem offended by the question. He doesn't even seem surprised. He offers the card again.")
			runtime.Sleep(5 * time.Second)
			c.BeginPhaseCredibility(runtime)
			return
		} else if strings.Contains(cmd, "take") || strings.Contains(cmd, "card") {
			// label: take the card
			runtime.HideInput()
			runtime.EmitNarration("You take the card. It's cold to the touch, and the silver mark is slightly textured and raised.")
			c.didTakeCard = true
			runtime.Sleep(5 * time.Second)
			c.BeginPhaseCredibility(runtime)
			return
		} else if strings.Contains(cmd, "refuse") || strings.Contains(cmd, "listen") || strings.Contains(cmd, "no") || strings.Contains(cmd, "slam") || strings.Contains(cmd, "shut") || strings.Contains(cmd, "door") {
			// label: refuse to listen
			runtime.HideInput()
			runtime.EmitNarration("You slam the door in his face. The sound echoes through the quiet neighborhood.")
			c.didSlamDoor = true
			runtime.Sleep(5 * time.Second)
			c.BeginPhaseCredibility(runtime)
			return
		}
	}

	runtime.ShowInputMessage("No such action")
}

func (c *Prologue) HandleTab(runtime engine.Runtime, input string) {
	switch c.phase {
	case prologuePhaseInitial:
		if c.toldToGetUp {
			runtime.Emit(engine.Event{Kind: engine.EventInsertInput, Text: "get up"})
			runtime.ShowInputMessage("I have to do everything myself...")
			return
		} else {
			runtime.ShowInputMessage(style.AsideKeyword.Render("Get") + style.Aside.Render(" up you lazy sod!"))
			c.toldToGetUp = true
			return
		}
	case prologuePhaseGotUp:
		runtime.ShowInputMessage("You could " + style.AsideKeyword.Render("stay silent") + style.Aside.Render(" or ") + style.AsideKeyword.Render("greet") + style.Aside.Render(" him"))
		return
	case prologuePhaseOpenedDoor:
		runtime.ShowInputMessage("You could " + style.AsideKeyword.Render("stay silent") + style.Aside.Render(" or ") + style.AsideKeyword.Render("question") + style.Aside.Render(" him"))
		return
	case prologuePhaseIntroduced:
		if c.finishedIntroductions {
			runtime.ShowInputMessage("You could " + style.AsideKeyword.Render("question") + style.Aside.Render(" him, ") + style.AsideKeyword.Render("take the card") + style.Aside.Render(", or ") + style.AsideKeyword.Render("refuse") + style.Aside.Render(" to listen"))
			return
		}
	}

	runtime.ShowInputMessage("No hint available :)")
}

func (c *Prologue) BeginPhaseInitial(runtime engine.Runtime) {
	// phase: initial
	c.phase = prologuePhaseInitial

	runtime.HideInput()
	runtime.EmitNarration("You awaken in your usual spot on the couch to three knocks at the door. The TV still on and empty pizza boxes on the coffee table.")
	runtime.Sleep(2 * time.Second)
	runtime.EmitNarration(style.Text.Render("Looking around, ") + style.TextItalic.Render("Holy shit I need a job") + style.Text.Render(", you think to yourself."))
	runtime.Sleep(2 * time.Second)
	runtime.EmitNarration(style.TextItalic.Render("Who the hell could that be?"))
	runtime.ShowInput()

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
			runtime.EmitNarration("Three more knocks at the door")
		}
	}()
}

func (c *Prologue) BeginPhaseGotUp(runtime engine.Runtime) {
	// phase: got up
	c.gettingUpCan()
	c.phase = prologuePhaseGotUp

	runtime.HideInput()

	runtime.EmitNarration("You get up and walk to the door.")
	runtime.Sleep(2 * time.Second)
	runtime.EmitNarration("Three more knocks at the door")
	runtime.Sleep(1 * time.Second)
	runtime.EmitNarration(style.TextItalic.Render("Yeah yeah I'm coming..."))
	runtime.Sleep(2 * time.Second)

	runtime.EmitNarration("")
	runtime.EmitNarration("The man at the door stands straight, though not rigidly. Dark coat. No visible insignia beyond a small silver mark pinned near the collar, circular, crossed by a single vertical line. Military, perhaps. Or adjacent to it.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("Rain beads across the shoulders of the coat.")
	//runtime.Sleep(3 * time.Second)

	runtime.ShowInputMessage("Hit enter or type to continue...")
	runtime.ShowInput()
}

func (c *Prologue) BeginPhaseOpenedDoor(runtime engine.Runtime) {
	// phase: opened door
	c.phase = prologuePhaseOpenedDoor

	runtime.HideInput()
	runtime.EmitNarration("\"Mr Paints?\"")
	runtime.Sleep(2 * time.Second)
	runtime.EmitNarration("His voice is calm, practiced.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("\"My name is Elias Vale. I represent the Luminous Outreach Department.\"")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("")
	runtime.EmitNarration("There is no grin, no theatricality, no attempt to sell his statement.")

	runtime.ShowInputMessage("Hit enter or type to continue...")
	runtime.ShowInput()
}

func (c *Prologue) BeginPhaseIntroduced(runtime engine.Runtime) {
	// phase: introduced
	c.phase = prologuePhaseIntroduced

	runtime.HideInput()
	runtime.EmitNarration("\"You are under no obligation to continue this conversation. You are not in danger. And you are not being accused of anything.\"")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("The man’s eyes are steady, but not unkind. He seems to be giving you the benefit of the doubt, as if he’s not sure whether you’re going to be cooperative or not.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("The rain taps softly against the path outside.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("\"We were given your name by someone familiar with your work in astrodynamics. Your research related to some unexplained noise in quantitative orbital survey data reached someone in our office. A lot of people file these irregularities as sensor malfunction. However, you did not.\"")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("He reaches into his coat slowly enough to avoid alarm, producing a slim matte card.")
	runtime.Sleep(4 * time.Second)
	runtime.EmitNarration("The card has no government seal nor any indication of its origin or national flag, but the same circular silver mark is embossed in the corner.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("\"There are things about the present state of humanity that you have not been told yet. We believe you are someone capable of understanding them responsibly.\"")

	runtime.BlockUntilContinue()

	runtime.EmitNarration("His eyes move briefly past you towards the dark house behind.")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("\"We’re not here because you’re special, Mr Paints. We’re here because you’re competent.\"")
	runtime.Sleep(5 * time.Second)
	runtime.EmitNarration("\"And because, at some point in your life, you continued spending your evenings looking at things that others wouldn’t even know are there.\"")

	c.finishedIntroductions = true
	runtime.EmitNarration("")
	runtime.ShowInput()
}

func (c *Prologue) BeginPhaseCredibility(runtime engine.Runtime) {
	// phase: credibility
	c.phase = prologuePhaseCredibility

	runtime.HideInput()
	if !c.didSlamDoor {
		// label: did not slam door
		if !c.didTakeCard {
			// label: sceptical without card
			runtime.EmitNarration("\"You don't believe me.\"")
			runtime.Sleep(3 * time.Second)
			runtime.EmitNarration("You reply, \"Too right, this is nonsense. You come to my door claiming to represent some secret civilisation and expect me to nod along?\"")
			runtime.Sleep(3 * time.Second)
			runtime.EmitNarration("\"I wouldn't believe me either.\", he replies.")
			runtime.Sleep(5 * time.Second)
		}
		// label: proof and interview
		runtime.EmitNarration("\"Unfortunately, everything important about this conversation sounds ridiculous when spoken aloud.\"")
		runtime.BlockUntilContinue()

		runtime.EmitNarration("")
		runtime.EmitNarration("You stir a little at his acknowledgement of the absurdity of the situation.")
		runtime.Sleep(5 * time.Second)
		runtime.EmitNarration("\"We are human,\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"We are not a government,\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"We are not a corporation,\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"We are not affiliated with any nation on Earth,\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"And we are not asking you for money.\"")
		runtime.Sleep(5 * time.Second)
		runtime.EmitNarration("You reply, \"Okay, if you’re not a government or a corporation, then who are you?\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"A society. A civilisation.\"")
		runtime.BlockUntilContinue()
		runtime.EmitNarration("\"So what, you've built some kind of utopia?\"")
		runtime.Sleep(2 * time.Second)
		runtime.EmitNarration("\"No.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"Then what is it?\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"A place.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("A pause.")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"You have three published papers.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"Two are cited regularly.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"One was largely ignored.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"The ignored one is why I’m here.\"")

		runtime.BlockUntilContinue()

		runtime.EmitNarration("\"We're inviting you to an interview.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("You pause for a second. The word 'interview' sounds so normal, so mundane, that it almost makes you forget the surreal circumstances of this conversation. Wait, an interview? Like a job?")
		runtime.Sleep(6 * time.Second)
		runtime.EmitNarration("\"An interview?\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"An interview.\"")

		runtime.BlockUntilContinue()

		runtime.EmitNarration("")
		runtime.EmitNarration("\"There is a train leaving Waverley tomorrow morning.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"If you board it, someone will meet you.\"")
		runtime.Sleep(3 * time.Second)
		runtime.EmitNarration("\"If you don’t, you’ll never hear from us again.\"")

		runtime.BlockUntilContinue()

		if !c.didTakeCard {
			// label: Elias leaves card
			runtime.EmitNarration("He places the card on the doorstep, steps back, and turns away.")
			runtime.Sleep(5 * time.Second)
		} else {
			// label: already holding card
			runtime.EmitNarration("He steps back and turns away, leaving the card in your hand.")
			runtime.Sleep(5 * time.Second)
		}

		runtime.EmitNarration("You watch him walk away, his figure receding into the rain and darkness.")
		runtime.Sleep(5 * time.Second)

		if !c.didTakeCard {
			// label: pick up doorstep card
			runtime.EmitNarration("The card still lies on the doorstep, its silver mark catching the dim light.")
			runtime.Sleep(5 * time.Second)
			runtime.EmitNarration("You pick it up and examine it. It's cold to the touch, and the silver mark is slightly textured and raised.")
		}
	} else {
		// label: card through letterbox
		runtime.EmitNarration("The card comes through the letterbox a few seconds later, landing on the doormat with barely a sound. You pick it up and examine it. It's cold to the touch, and the silver mark is slightly textured and raised.")
	}
}
