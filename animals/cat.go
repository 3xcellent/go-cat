package animals

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const defaultColor = uint32(0x44ff99)
const crashColor = uint32(0xff0000)
const JumpFreq = time.Second / 4

const (
	Sitting       = 1
	Walking       = 0
	Transitioning = 2
	Running       = 3
)

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
	maxSpeed  float64
	JumpedAt  time.Time

	PosX   int32
	PosY   int32
	Height int32
	Width  int32

	Velocity Velocity
}

type Velocity struct {
	Up    float64
	Down  float64
	Left  float64
	Right float64
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
		maxSpeed:         8 * 2,
		renderer:         renderer,
		texture:          texture,
		spriteWidth:      imgWidth,
		spriteHeight:     imgHeight,
		frameWidth:       imgWidth / 4,   // there are 4 columns of frames
		frameHeight:      imgHeight / 13, // there are 13 columns of frames
		direction:        Right,
		currentActionIdx: Sitting,
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
	case Walking:
		return "walking"
	case Sitting:
		return "sitting"
	case Transitioning:
		return "transitioning"
	default:
		return "running"
	}
}

func (c *Cat) Move(hasUpPressed, hasDownPressed, hasLeftPressed, hasRightPressed, hasSpacebarPressed bool, displayWidth, displayHeight int32) {
	//fmt.Printf("hasUpPressed: %t, hasDownPressed: %t, hasLeftPressed: %t, hasRightPressed: %t, displayWidth: %d, displayHeight: %d\n", hasUpPressed, hasDownPressed, hasLeftPressed, hasRightPressed, displayWidth, displayHeight)

	// handle jump
	if hasSpacebarPressed {
		hasUpPressed = true
	}

	//handle 'gravity'
	c.Velocity.Down++

	isFalling := true
	//detect 'platform'
	if c.PosY+c.frameHeight/4 >= displayHeight-1 {
		isFalling = false
		c.Velocity.Down = 0
		c.PosY = displayHeight - c.frameHeight/4
	} else {
		// falling
		if c.Velocity.Down > c.maxSpeed/2 {
			c.currentFrameIdx = 3
			c.currentActionIdx = Running
		}
		if c.Velocity.Down > c.maxSpeed {
			c.currentFrameIdx = 2
		}
		c.PosY = c.PosY + int32(c.Velocity.Down)
	}

	// jumping
	if !isFalling && hasUpPressed && c.JumpedAt.Add(JumpFreq).Before(time.Now()) {
		c.currentActionIdx = Running
		c.currentFrameIdx = 11
		c.JumpedAt = time.Now()
		c.Velocity.Up = c.Velocity.Up + c.maxSpeed*2
	}
	if c.Velocity.Up > 0 {
		// the effect here is that holding jump will jump a little higher
		c.currentActionIdx = Running
		if hasUpPressed {
			c.Velocity.Up = c.Velocity.Up - 0.06
			c.currentFrameIdx = 12
		} else {
			c.Velocity.Up = c.Velocity.Up - c.Velocity.Down
			c.currentFrameIdx = 0
		}
		c.Velocity.Up = c.Velocity.Up / 1.2

		newY := c.PosY - int32(c.Velocity.Up)
		if newY < 0 {
			newY = 0
		}
		c.PosY = newY
	}

	// Pressing Down should not do anything now
	//if hasDownPressed {
	//	c.Velocity.Down = int32(math.Min(float64(c.Velocity.Down+1), float64(c.maxSpeed)))
	//}
	//if c.Velocity.Down > 0 {
	//	if !hasDownPressed {
	//		c.Velocity.Down--
	//	}
	//	newY := c.PosY + c.Velocity.Down
	//
	//	if newY > (displayHeight - c.Height) {
	//		newY = displayHeight - c.Height
	//		c.currentColor = crashColor
	//	} else {
	//		c.currentColor = defaultColor
	//	}
	//	 = newY
	//}

	if hasLeftPressed {
		c.Velocity.Left = math.Min(c.Velocity.Left+1, c.maxSpeed)
	}
	if c.Velocity.Left > 0 {
		if !hasLeftPressed {
			c.Velocity.Left--
		}
		newX := c.PosX - int32(c.Velocity.Left)
		if newX < 0 {
			newX = 0
		}
		c.PosX = newX
	}

	if hasRightPressed {
		c.Velocity.Right = math.Min(c.Velocity.Right+1, c.maxSpeed)
	}
	if c.Velocity.Right > 0 {
		if !hasRightPressed {
			c.Velocity.Right--
		}
		newX := c.PosX + int32(c.Velocity.Right)
		if newX > displayWidth-c.Width {
			newX = displayWidth - c.Width
		} else {
		}
		c.PosX = newX
	}

	// handle sprite direction/angle
	if c.Velocity.Left > c.Velocity.Right {
		c.direction = Left
	} else if c.Velocity.Right > c.Velocity.Left {
		c.direction = Right
	}

	// handle sprite action and frame
	if isFalling {
		if c.direction == Left {
			c.angle = float64(c.Velocity.Up*2 - c.Velocity.Down*2)
		} else {
			c.angle = float64(c.Velocity.Down*2 - c.Velocity.Up*2)
		}
		return
	}
	c.angle = 0

	if c.Velocity.Left == c.maxSpeed || c.Velocity.Right == c.maxSpeed {
		if c.currentActionIdx != Running {
			c.currentActionIdx = Running
			c.currentFrameIdx = 0
		}
		return
	}

	if c.Velocity.Left != c.Velocity.Right {
		if c.currentActionIdx != Walking {
			c.currentActionIdx = Walking
			c.currentFrameIdx = 0
		}
		return
	}

	if c.currentActionIdx != Sitting {
		c.currentActionIdx = Sitting
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
	err := c.renderer.CopyEx(c.texture, frameRect, &sdl.Rect{c.PosX, c.PosY, c.frameWidth / 4, c.frameHeight / 4}, c.angle, nil, c.direction)
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
