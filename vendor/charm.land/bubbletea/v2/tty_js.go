//go:build js
// +build js

package tea

import (
	"errors"
)

const suspendSupported = false

func (p *Program) initInput() error {
	return nil
}

func suspendProcess() {
	panic(errors.New("suspend is not supported on the web"))
}
