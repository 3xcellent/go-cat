package game

import (
	"github.com/3xcellent/go-cat/animals"
	"github.com/veandco/go-sdl2/sdl"
)

const defaultBackgroundColor = uint32(0x5090ee)

type Game struct {
	renderer               *sdl.Renderer
	running                bool
	HasUpPressed           bool
	HasDownPressed         bool
	HasLeftPressed         bool
	HasRightPressed        bool
	DisplayWidth           int32
	DisplayHeight          int32
	currentBackgroundColor uint32
	cat                    *animals.Cat
}

func CreateGame(width int32, height int32, renderer *sdl.Renderer) *Game {
	return &Game{
		DisplayWidth:           width,
		DisplayHeight:          height,
		currentBackgroundColor: defaultBackgroundColor,
		renderer:               renderer,
		cat:                    animals.CreateCat(renderer),
	}
}

func (g *Game) Start() {
	g.running = true
}

func (g *Game) IsRunning() bool {
	return g.running
}

func (g *Game) ClearRenderer() {
	err := g.renderer.Clear()
	if err != nil {
		panic(err)
	}
}

func (g *Game) DrawBackground() {
	err := g.renderer.SetDrawColor(155, 155, 255, 255)
	if err != nil {
		panic(err)
	}
}

func (g *Game) HandleKeyboardEvent(event *sdl.KeyboardEvent) {
	//  Sample: &sdl.KeyboardEvent{
	//  	Type:0x300,
	//  	Timestamp:0x132d,
	//  	WindowID:0x1,
	//  	State:0x1,
	//  	Repeat:0x1,
	//  	Keysym:sdl.Keysym{
	//  		Scancode:0x4f,
	//  		Sym:1073741903,
	//  		Mod:0x0,
	// 		}
	// 	}
	switch event.Keysym.Scancode {
	case sdl.SCANCODE_ESCAPE:
		println("Quit")
		g.running = false
		break
	case sdl.SCANCODE_UP:
		g.HasUpPressed = event.Type == sdl.KEYDOWN
	case sdl.SCANCODE_DOWN:
		g.HasDownPressed = event.Type == sdl.KEYDOWN
	case sdl.SCANCODE_LEFT:
		g.HasLeftPressed = event.Type == sdl.KEYDOWN
	case sdl.SCANCODE_RIGHT:
		g.HasRightPressed = event.Type == sdl.KEYDOWN
	default:
	}
}

func (g *Game) Stop() {
	g.running = false
}
func (g *Game) Arrange() {
	g.cat.Move(g.HasUpPressed, g.HasDownPressed, g.HasLeftPressed, g.HasRightPressed, g.DisplayWidth-200, g.DisplayHeight-200)
}
func (g *Game) Draw() {
	g.DrawBackground()
	g.cat.Draw()
	g.renderer.Present()
	g.ClearRenderer()
}
