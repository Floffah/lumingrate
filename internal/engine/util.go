package engine

import "time"

const (
	averageNarrationWordLength     = 5
	averageNarrationWordsPerMinute = 220
)

func (e *Engine) EmitNarration(text string) {
	e.Emit(Event{Kind: EventNarration, Text: text})
}

func (e *Engine) EmitNarrationAndWait(text string) {
	e.EmitNarration(text)
	e.Sleep(narrationWaitDuration(text))
}

func narrationWaitDuration(text string) time.Duration {
	return time.Duration(len(text)) * time.Minute / (averageNarrationWordLength * averageNarrationWordsPerMinute)
}

func (e *Engine) EmitAside(text string) {
	e.Emit(Event{Kind: EventAside, Text: text})
}

func (e *Engine) ShowInput() {
	e.Emit(Event{Kind: EventShowInput})
}

func (e *Engine) HideInput() {
	e.Emit(Event{Kind: EventHideInput})
}

func (e *Engine) ShowInputMessage(text string) {
	e.Emit(Event{Kind: EventInputMessage, Text: text})
}

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
		e.ShowInputMessage("No hints available :)")
	case 1:
		e.Emit(Event{Kind: EventInsertInput, Text: matches[0][len(command):]})
	default:
		e.ShowInputMessage("Hint: " + JoinForDisplay(matches) + "")
	}
}

func (e *Engine) BlockUntilContinueWithText(message string) string {
	e.ShowInputMessage(message)
	e.ShowInput()

	defer func() {
		e.ShowInputMessage("")
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
				if NormaliseCommand(command.Text) == "" {
					e.HideInput()
					return ""
				}
				e.ShowInputMessage(message)
			} else if command.Kind == CommandTab {
				e.handleTab(command.Text)
			} else {
				e.ShowInputMessage(message)
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

func (e *Engine) ConsumeCommandsInline() []Command {
	var commands []Command
	for {
		select {
		case <-e.ctx.Done():
			return commands
		case command, ok := <-e.commands:
			if !ok {
				return commands
			}
			if command.Kind == CommandSubmit {
				commands = append(commands, command)
			} else if command.Kind == CommandTab {
				e.handleTab(command.Text)
			}
		default:
			return commands
		}
	}
}

func (e *Engine) SwitchChapter(id ChapterID) {
	e.EmitAside("Chapter " + string(id) + " is not available from this runtime.")
}

func (e *Engine) Sleep(d time.Duration) {
	if d <= 0 {
		return
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-e.done():
	case <-timer.C:
	case <-e.fastForwardDone():
	}
}

func (e *Engine) done() <-chan struct{} {
	if e.ctx == nil {
		return nil
	}

	return e.ctx.Done()
}

func (e *Engine) fastForwardDone() <-chan struct{} {
	e.fastForwardMu.Lock()
	defer e.fastForwardMu.Unlock()

	return e.fastForwardWake
}
