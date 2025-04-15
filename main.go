package main

import (
	"math/rand"
	"syscall/js"
)

type Application = js.Value
type Texture = js.Value
type Sprite = js.Value

const NUM_SPRITES = 20_000

var PIXI = js.Global().Get("PIXI")

var sprites = make([]Sprite, NUM_SPRITES)

var rotation = 0.0

func await(promise js.Value) js.Value {
	ch := make(chan js.Value, 1)

	then := js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args[0]
		return nil
	})
	promise.Call("then", then)

	// block until resolved
	result := <-ch
	then.Release()

	return result
}

func initPixiApp() Application {
	return await(js.Global().Call("initPixi"))
}

func newSprite(texture Texture) Sprite {
	sprite := PIXI.Get("Sprite").New(texture)
	sprite.Set("x", rand.Float32()*1600)
	sprite.Set("y", rand.Float32()*800)
	sprite.Set("anchor", map[string]any{"x": 0.5, "y": 0.5})
	return sprite
}

func onTick(this js.Value, args []js.Value) any {
	rotation += 0.01

	for _, sprite := range sprites {
		sprite.Set("rotation", rotation)
	}
	return nil
}

func main() {
	pixi := js.Global().Get("PIXI")

	// app := pixi.Get("Application").New(map[string]any{
	// 	"width":  800,
	// 	"height": 600,
	// })
	// js.Global().Get("document").Get("body").Call("appendChild", app.Get("canvas"))

	app := initPixiApp()

	texture := await(pixi.Get("Assets").Call("load", "https://pixijs.io/examples/examples/assets/bunny.png"))

	for i := range NUM_SPRITES {
		sprite := newSprite(texture)
		sprites[i] = sprite
		app.Get("stage").Call("addChild", sprite)
	}

	ticker := app.Get("ticker")
	ticker.Call("add", js.FuncOf(onTick))

	select {} // keep Go running
}
