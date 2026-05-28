//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"syscall/js"

	tea "charm.land/bubbletea/v2"

	"luminrate/internal/scene/app"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

var (
	stateMu sync.Mutex
	program *tea.Program
	input   *webInput
	cancel  context.CancelFunc

	registeredFuncs []js.Func
)

func main() {
	api := js.Global().Get("Object").New()

	start := js.FuncOf(startProgram)
	sendInput := js.FuncOf(sendProgramInput)
	resize := js.FuncOf(resizeProgram)
	stop := js.FuncOf(stopProgram)

	registeredFuncs = append(registeredFuncs, start, sendInput, resize, stop)
	api.Set("start", start)
	api.Set("input", sendInput)
	api.Set("resize", resize)
	api.Set("stop", stop)

	js.Global().Set("lumingrate", api)
	select {}
}

func startProgram(_ js.Value, args []js.Value) any {
	opts := parseStartOptions(args)
	if opts.write.Type() != js.TypeFunction {
		return jsError("lumingrate.start requires a write function")
	}

	stateMu.Lock()
	defer stateMu.Unlock()

	stopLocked()

	ctx, can := context.WithCancel(context.Background())
	in := newWebInput()
	out := &webOutput{write: opts.write}
	p := tea.NewProgram(
		app.NewModel(),
		tea.WithContext(ctx),
		tea.WithInput(in),
		tea.WithOutput(out),
		tea.WithWindowSize(opts.cols, opts.rows),
		tea.WithEnvironment(webEnvironment(opts.cols, opts.rows)),
		tea.WithoutSignalHandler(),
	)

	input = in
	cancel = can
	program = p

	go func() {
		_, err := p.Run()
		_ = in.Close()

		stateMu.Lock()
		if program == p {
			program = nil
			input = nil
			cancel = nil
		}
		stateMu.Unlock()

		if err != nil && ctx.Err() == nil {
			out.WriteString(fmt.Sprintf("lumingrate: %v\r\n", err))
		}
	}()

	return nil
}

func sendProgramInput(_ js.Value, args []js.Value) any {
	if len(args) == 0 {
		return nil
	}

	stateMu.Lock()
	in := input
	stateMu.Unlock()
	if in == nil {
		return nil
	}

	if err := in.Push(args[0].String()); err != nil {
		return jsError(err.Error())
	}
	return nil
}

func resizeProgram(_ js.Value, args []js.Value) any {
	cols := defaultWidth
	rows := defaultHeight
	if len(args) > 0 {
		cols = positiveInt(args[0], defaultWidth)
	}
	if len(args) > 1 {
		rows = positiveInt(args[1], defaultHeight)
	}

	stateMu.Lock()
	p := program
	stateMu.Unlock()
	if p != nil {
		p.Send(tea.WindowSizeMsg{Width: cols, Height: rows})
	}
	return nil
}

func stopProgram(_ js.Value, _ []js.Value) any {
	stateMu.Lock()
	defer stateMu.Unlock()
	stopLocked()
	return nil
}

func stopLocked() {
	if cancel != nil {
		cancel()
		cancel = nil
	}
	if input != nil {
		_ = input.Close()
		input = nil
	}
	if program != nil {
		program.Quit()
		program = nil
	}
}

type startOptions struct {
	write js.Value
	cols  int
	rows  int
}

func parseStartOptions(args []js.Value) startOptions {
	opts := startOptions{
		cols: defaultWidth,
		rows: defaultHeight,
	}
	if len(args) == 0 {
		return opts
	}

	if args[0].Type() == js.TypeFunction {
		opts.write = args[0]
		if len(args) > 1 {
			opts.cols = positiveInt(args[1], defaultWidth)
		}
		if len(args) > 2 {
			opts.rows = positiveInt(args[2], defaultHeight)
		}
		return opts
	}

	opts.write = args[0].Get("write")
	opts.cols = positiveInt(args[0].Get("cols"), defaultWidth)
	opts.rows = positiveInt(args[0].Get("rows"), defaultHeight)
	return opts
}

func positiveInt(value js.Value, fallback int) int {
	if value.IsUndefined() || value.IsNull() {
		return fallback
	}
	n := value.Int()
	if n <= 0 {
		return fallback
	}
	return n
}

func webEnvironment(cols, rows int) []string {
	return []string{
		"TERM=xterm-256color",
		"TERM_PROGRAM=ghostty",
		"COLORTERM=truecolor",
		"CLICOLOR_FORCE=1",
		fmt.Sprintf("COLUMNS=%d", cols),
		fmt.Sprintf("LINES=%d", rows),
	}
}

func jsError(message string) js.Value {
	return js.Global().Get("Error").New(message)
}

type webInput struct {
	mu     sync.Mutex
	cond   *sync.Cond
	buffer []byte
	closed bool
}

func newWebInput() *webInput {
	in := &webInput{}
	in.cond = sync.NewCond(&in.mu)
	return in
}

func (in *webInput) Read(p []byte) (int, error) {
	in.mu.Lock()
	defer in.mu.Unlock()

	for len(in.buffer) == 0 && !in.closed {
		in.cond.Wait()
	}
	if len(in.buffer) == 0 && in.closed {
		return 0, io.EOF
	}

	n := copy(p, in.buffer)
	in.buffer = in.buffer[n:]
	return n, nil
}

func (in *webInput) Push(data string) error {
	if data == "" {
		return nil
	}

	in.mu.Lock()
	defer in.mu.Unlock()
	if in.closed {
		return io.ErrClosedPipe
	}

	in.buffer = append(in.buffer, data...)
	in.cond.Broadcast()
	return nil
}

func (in *webInput) Close() error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if in.closed {
		return nil
	}

	in.closed = true
	in.cond.Broadcast()
	return nil
}

type webOutput struct {
	mu     sync.Mutex
	write  js.Value
	lastCR bool
}

func (out *webOutput) Write(p []byte) (int, error) {
	if out.write.Type() != js.TypeFunction {
		return len(p), nil
	}

	out.write.Invoke(out.normaliseNewlines(p))
	return len(p), nil
}

func (out *webOutput) WriteString(s string) {
	_, _ = out.Write([]byte(s))
}

func (out *webOutput) normaliseNewlines(p []byte) string {
	out.mu.Lock()
	defer out.mu.Unlock()

	buf := make([]byte, 0, len(p))
	lastCR := out.lastCR
	for _, b := range p {
		if b == '\n' && !lastCR {
			buf = append(buf, '\r')
		}
		buf = append(buf, b)
		lastCR = b == '\r'
	}
	out.lastCR = lastCR
	return string(buf)
}
