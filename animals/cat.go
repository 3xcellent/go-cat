package animals

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const defaultColor = uint32(0x44ff99)
const crashColor = uint32(0xff0000)

type Cat struct {
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	spriteWidth  int32
	spriteHeight int32
	frameWidth   int32
	frameHeight  int32
	angle        float64

	currentFrameIdx  int32
	currentActionIdx int32
	catActionFrames  map[string][]Frame

	direction sdl.RendererFlip
	maxSpeed  int32

	currentColor uint32

	PosX   int32
	PosY   int32
	Height int32
	Width  int32

	Velocity Velocity
}

type Velocity struct {
	Up    int32
	Down  int32
	Left  int32
	Right int32
}

type Frame struct {
	X int32
	Y int32
}

const (
	Right = sdl.FLIP_HORIZONTAL
	Left  = sdl.FLIP_NONE
)

func CreateCat(renderer *sdl.Renderer) *Cat {
	img, err := img.Load("assets/cat.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		os.Exit(4)
	}

	imgWidth := img.W
	imgHeight := img.H

	texture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(5)
	}
	img.Free()

	cat := &Cat{
		maxSpeed:         8,
		renderer:         renderer,
		texture:          texture,
		spriteWidth:      imgWidth,
		spriteHeight:     imgHeight,
		currentColor:     defaultColor,
		frameWidth:       imgWidth / 4,   // there are 4 columns of frames
		frameHeight:      imgHeight / 13, // there are 13 columns of frames
		direction:        Right,
		currentActionIdx: 3,
		PosX:             100,
		PosY:             150,
		Height:           20,
		Width:            15,
	}
	// TODO instead of having LoadAnimationFrames set the frames on the object after the fact, it should
	//  be set above and the the func returns the map
	err = cat.LoadAnimationFrames()
	if err != nil {
		panic(fmt.Sprintf("error loading frames: %v", err))
	}
	return cat
}

func (c *Cat) LoadAnimationFrames() error {
	// walking
	c.catActionFrames = map[string][]Frame{}
	numFrames := int32(12)
	walkingFrames := make([]Frame, 0)
	for i := int32(0); i < numFrames; i++ {
		walkingFrames = append(walkingFrames, Frame{X: 0, Y: i * c.frameHeight})
	}
	c.catActionFrames["walking"] = walkingFrames

	// sitting
	numFrames = int32(6)
	sittingFrames := make([]Frame, 0)
	for i := int32(0); i < numFrames; i++ {
		sittingFrames = append(sittingFrames, Frame{X: c.frameWidth, Y: i * c.frameHeight})
	}
	c.catActionFrames["sitting"] = sittingFrames

	// transitioning
	numFrames = int32(12)
	transitioningFrames := make([]Frame, 0)
	for i := int32(0); i < numFrames; i++ {
		transitioningFrames = append(transitioningFrames, Frame{X: 2 * c.frameWidth, Y: i * c.frameHeight})
	}
	c.catActionFrames["transitioning"] = transitioningFrames

	// running
	numFrames = int32(13)
	runningFrames := make([]Frame, 0)
	for i := int32(0); i < numFrames; i++ {
		runningFrames = append(runningFrames, Frame{X: 3 * c.frameWidth, Y: i * c.frameHeight})
	}
	c.catActionFrames["running"] = runningFrames

	return nil
}

func (c *Cat) Action() string {
	switch c.currentActionIdx {
	case 0:
		return "walking"
	case 1:
		return "sitting"
	case 2:
		return "transitioning"
	default:
		return "running"
	}
}

func (c *Cat) Move(hasUpPressed, hasDownPressed, hasLeftPressed, hasRightPressed bool, displayWidth, displayHeight int32) {
	//fmt.Printf("hasUpPressed: %t, hasDownPressed: %t, hasLeftPressed: %t, hasRightPressed: %t, displayWidth: %d, displayHeight: %d\n", hasUpPressed, hasDownPressed, hasLeftPressed, hasRightPressed, displayWidth, displayHeight)
	if hasUpPressed {
		c.Velocity.Up = int32(math.Min(float64(c.Velocity.Up+1), float64(c.maxSpeed)))
	}
	if c.Velocity.Up > 0 {
		if !hasUpPressed {
			c.Velocity.Up--
		}
		newY := c.PosY - c.Velocity.Up
		if newY < 0 {
			newY = 0
			c.currentColor = crashColor
		} else {
			c.currentColor = defaultColor
		}
		c.PosY = newY
	}

	if hasDownPressed {
		c.Velocity.Down = int32(math.Min(float64(c.Velocity.Down+1), float64(c.maxSpeed)))
	}
	if c.Velocity.Down > 0 {
		if !hasDownPressed {
			c.Velocity.Down--
		}
		newY := c.PosY + c.Velocity.Down

		if newY > (displayHeight - c.Height) {
			newY = displayHeight - c.Height
			c.currentColor = crashColor
		} else {
			c.currentColor = defaultColor
		}
		c.PosY = newY
	}

	if hasLeftPressed {
		c.Velocity.Left = int32(math.Min(float64(c.Velocity.Left+1), float64(c.maxSpeed)))
	}
	if c.Velocity.Left > 0 {
		if !hasLeftPressed {
			c.Velocity.Left--
		}
		newX := c.PosX - c.Velocity.Left
		if newX < 0 {
			newX = 0
			c.currentColor = crashColor
		} else {
			c.currentColor = defaultColor
		}
		c.PosX = newX
	}

	if hasRightPressed {
		c.Velocity.Right = int32(math.Min(float64(c.Velocity.Right+1), float64(c.maxSpeed)))
	}
	if c.Velocity.Right > 0 {
		if !hasRightPressed {
			c.Velocity.Right--
		}
		newX := c.PosX + c.Velocity.Right
		if newX > displayWidth-c.Width {
			newX = displayWidth - c.Width
			c.currentColor = crashColor
		} else {
			c.currentColor = defaultColor
		}
		c.PosX = newX
	}

	if c.Velocity.Left > c.Velocity.Right {
		c.direction = Left
		c.angle = float64(c.Velocity.Up*4 - c.Velocity.Down*4)
		if c.Velocity.Left == c.maxSpeed {
			if c.currentActionIdx != 3 {
				c.currentActionIdx = 3
				c.currentFrameIdx = 0
			}
			return
		}
		if c.currentActionIdx != 0 {
			c.currentActionIdx = 0
			c.currentFrameIdx = 0
		}
		return
	}
	if c.Velocity.Right > c.Velocity.Left {
		c.direction = Right
		c.angle = float64(c.Velocity.Down*4 - c.Velocity.Up*4)
		if c.Velocity.Right == c.maxSpeed {
			if c.currentActionIdx != 3 {
				c.currentActionIdx = 3
				c.currentFrameIdx = 0
			}
			return
		}
		if c.currentActionIdx != 0 {
			c.currentActionIdx = 0
			c.currentFrameIdx = 0
		}
		return
	}
	if c.currentActionIdx != 1 {
		c.currentActionIdx = 1
		c.currentFrameIdx = 0
	}
}

func (c *Cat) Draw() {
	frameRect := &sdl.Rect{
		X: c.catActionFrames[c.Action()][c.currentFrameIdx].X,
		Y: c.catActionFrames[c.Action()][c.currentFrameIdx].Y,
		W: c.frameWidth,
		H: c.frameHeight,
	}

	// flip options: sdl.FLIP_NONE, sdl.FLIP_HORIZONTAL, sdl.SDL_FLIP_VERTICAL, sdl.FLIP_HORIZONTAL | sdl.SDL_FLIP_VERTICAL
	err := c.renderer.CopyEx(c.texture, frameRect, &sdl.Rect{c.PosX, c.PosY, c.frameWidth / 8, c.frameHeight / 8}, c.angle, nil, c.direction)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to renderer.CopyEx texture: %s\n", err)
		os.Exit(5)
	}

	if c.Action() == "sitting" && c.currentFrameIdx == 5 {
		return
	}
	c.currentFrameIdx++
	if c.currentFrameIdx >= int32(len(c.catActionFrames[c.Action()])) {
		c.currentFrameIdx = 0
	}
}
