package engine

import "time"

func (e *Engine) EmitTypingText(text string) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for _, c := range text {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.Emit(Event{Kind: EventRaw, Text: string(c)})
		}
	}
}

func (e *Engine) EmitCommandsHint(input string, commands []string) {
	command := NormaliseCommand(input)
	matches := matchingCommands(commands, command)
	switch len(matches) {
	case 0:
		e.Emit(Event{Kind: EventInputMessage, Text: "No hints available :)"})
	case 1:
		e.Emit(Event{Kind: EventInsertInput, Text: matches[0][len(command):]})
	default:
		e.Emit(Event{Kind: EventInputMessage, Text: "Hint: " + JoinForDisplay(matches) + ""})
	}
}

func (e *Engine) BlockUntilContinueWithText(message string) string {
	e.Emit(Event{Kind: EventInputMessage, Text: message})
	e.Emit(Event{Kind: EventShowInput})

	defer func() {
		e.Emit(Event{Kind: EventInputMessage, Text: ""})
	}()

	for {
		select {
		case <-e.ctx.Done():
			return ""
		case command, ok := <-e.commands:
			if !ok {
				return ""
			}
			if command.Kind == CommandSubmit {
				return NormaliseCommand(command.Text)
			} else if command.Kind == CommandTab {
				e.handleTab(command.Text)
			} else {
				e.Emit(Event{Kind: EventInputMessage, Text: message})
			}
		}
	}
}

func (e *Engine) BlockUntilContinue() string {
	return e.BlockUntilContinueWithText("Press enter to continue...")
}

func (e *Engine) ConsumeCommandsNow() {
	for {
		select {
		case <-e.ctx.Done():
			return
		case command, ok := <-e.commands:
			if !ok {
				return
			}
			if command.Kind == CommandSubmit {
				e.handleSubmit(command.Text)
			}
		default:
			return
		}
	}
}
