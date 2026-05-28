//go:build js && wasm

package main

import (
	"context"
	"luminrate/internal/engine"
	"luminrate/internal/story"
)

var eng *engine.Engine
var ctx context.Context
var can context.CancelFunc

func main() {}

func Start() {
	ctx, can = context.WithCancel(context.Background())

	eng = story.NewEngine()
	go eng.Run(ctx)
}
