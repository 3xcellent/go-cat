package main

import (
	"time"

	"github.com/3xcellent/go-cat/game"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var targetFramerate time.Duration

func init() {
	targetFramerate = time.Second / time.Duration(30)
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	display, err := sdl.GetDisplayBounds(0)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Cat", 100, 100, display.W-200, display.H-200, sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.Clear()
	img.Init(img.INIT_JPG | img.INIT_PNG)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	g := game.CreateGame(display.W, display.H, renderer)
	g.Start()
	for g.IsRunning() {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.KeyboardEvent:
				handleEvent := event.(*sdl.KeyboardEvent)
				g.HandleKeyboardEvent(handleEvent)
			case *sdl.QuitEvent:
				println("Quit")
				g.Stop()
				break
			}
		}

		g.Arrange()
		g.Draw()

		waitTime := targetFramerate - time.Now().Sub(frameStart)
		time.Sleep(waitTime)
	}
}
