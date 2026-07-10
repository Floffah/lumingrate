package chapter

import (
	"luminrate/internal/engine"
	"strings"
)

type thecardPhase int

const (
	thecardPhaseInitial thecardPhase = iota
	thecardPhaseSleep
)

// TheCard Chapter 2
type TheCard struct {
	phase thecardPhase

	didLookAround     bool
	amtInitialActions int
}

func NewTheCard() *TheCard {
	return &TheCard{
		phase:             thecardPhaseInitial,
		didLookAround:     false,
		amtInitialActions: 0,
	}
}

func (c *TheCard) Begin(runtime engine.Runtime) {
	c.BeginPhaseInitial(runtime)
}

func (c *TheCard) HandleCommand(runtime engine.Runtime, cmd string) {
	switch c.phase {
	case thecardPhaseInitial:
		if c.amtInitialActions >= 4 && c.amtInitialActions < 6 {
			runtime.HideInput()
			defer runtime.ShowInput()
			runtime.EmitNarrationAndWait("A yawn takes you over. A feeling of self-consciousness briefly washes over you. Maybe it's time to sleep.")
		}

		if !c.didLookAround && strings.Contains(cmd, "look") {
			c.didLookAround = true
			c.amtInitialActions += 1
			runtime.EmitNarration("Looking around the room, pizza boxes paint the coffee table, the TV is murmuring some odd show, your bed is unmade and unslept in, rain patters at the windows, and the computer is still open on a browser tab about a card game you looked up earlier.")
			return
		} else if strings.Contains(cmd, "search") {
			if strings.Contains(cmd, "luminous") || strings.Contains(cmd, "card") {
				c.amtInitialActions += 1
				runtime.EmitNarration("You open a browser tab and search for 'luminous' and 'card'. All you find is companies specialising in business cards, and a few card games. Maybe I could search about Elias?")
				return
			} else if strings.Contains(cmd, "elias") || strings.Contains(cmd, "vale") {
				c.amtInitialActions += 1
				runtime.EmitNarration("You search for 'Elias Vale' and find a few social media profiles, but none with profile pictures that match the man you met. I guess I won't find out any more unless I go... when was the train again?")
				return
			}
		} else if strings.Contains(cmd, "paper") {
			c.amtInitialActions += 1
			runtime.HideInput()
			runtime.EmitNarrationAndWait("You open the old paper. The ignored one. Forty-three pages of careful modelling, yet it got several rejections, and some reviewers who clearly stopped reading halfway through.")
			runtime.EmitNarration("You remember telling yourself it didn't matter. Apparently that was optimistic.")
			runtime.ShowInput()
			return
		} else if strings.Contains(cmd, "train") {
			c.amtInitialActions += 1
			runtime.EmitNarration("You check the train times, there is in fact a train departing at 9am. Rush hour, really? Maybe I should get some sleep. Maybe I should check the card again, just to be sure.")
			return
		} else if strings.Contains(cmd, "card") {
			c.amtInitialActions += 1
			runtime.HideInput()
			runtime.EmitNarrationAndWait("The card says exactly what it said before. LUMINOUS OUTREACH DEPARTMENT. EDINBURGH WAVERLEY. 09:00. BRING THIS CARD.")
			runtime.EmitNarration("You turn it over again. Still blank. Well, apparently a civilisation beyond Earth doesn't believe in phone numbers.")
			runtime.ShowInput()
			return
		} else if strings.Contains(cmd, "sleep") || c.amtInitialActions >= 6 {
			runtime.HideInput()
			if c.amtInitialActions >= 6 {
				runtime.EmitNarrationAndWait("After dwelling on this situation for a while, you decide it's time to sleep.")
			}
			c.BeginPhaseSleep(runtime)
			return
		}
	}

	runtime.ShowInputMessage("No such action")
}

func (c *TheCard) HandleTab(runtime engine.Runtime, input string) {
	switch c.phase {
	case thecardPhaseInitial:
		if !c.didLookAround {
			runtime.ShowInputMessage("Look around")
		} else {
			runtime.ShowInputMessage("Tips will be relaxed from now on, there's more you can do than hinted or in prose :) (get creative)")
		}
		return
	}

	runtime.ShowInputMessage("No hint available :)")
}

func (c *TheCard) BeginPhaseInitial(runtime engine.Runtime) {
	runtime.HideInput()
	c.phase = thecardPhaseInitial
	runtime.EmitNarrationAndWait("Elias is long gone into the night. The TV is still murmuring, rain is pattering at the windows. You're left with the card.")
	runtime.ShowInput()
}

func (c *TheCard) BeginPhaseSleep(runtime engine.Runtime) {
	c.phase = thecardPhaseSleep
	runtime.HideInput()
	runtime.EmitNarrationAndWait("You drop the card into the bin.")
	runtime.EmitNarrationAndWait("For quite a few seconds, this feels decisive.")
	runtime.EmitNarrationAndWait("Then you take it back out. Of course you do.")
	runtime.EmitNarrationAndWait("You set an alarm for 07:30.")
	runtime.EmitNarrationAndWait("The card remains on the desk, catching the weak blue light from the TV.")
	runtime.EmitNarrationAndWait("You go to bed, and as you drift off to sleep, you feel a sense of unease. The card is still on your desk, and you can't shake the feeling that something about this might just be true.")
}
