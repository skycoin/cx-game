package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/urfave/cli/v2"
)

//Press TAB to shuffle stars

func init() {
	// seed rand so stars will be random each program run
	rand.Seed(time.Now().UnixNano())
	//lock thread so drawing will be only in main thread, otherwise there will be errors
	runtime.LockOSThread()
}

type Star struct {
	// Drawable uint32
	X        float32
	Y        float32
	Size     float32
	SpriteId int
	Depth    float32
}

const (
	width  = 800
	height = 600
)

var (
	wx, wy, wz float64 = 0, 0, -10
	size       float64 = 1
	spriteId   int
	stars      []*Star
)

func main() {
	//parse command line arguments and flags
	// initArgs()

	win := render.NewWindow(height, width, true)
	defer glfw.Terminate()
	window := win.Window
	program := win.Program

	window.SetKeyCallback(keyCallback)

	//spriteloader init
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/starfield_test_16x16_tiles_8x8_tile_grid_128x128.png")

	//load all sprites from spritesheet
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			//load sprites
			spriteloader.LoadSprite(spriteSheetId,
				fmt.Sprintf("stars-%d", x*4+y),
				x, y)
		}
	}

	for x := 0; x < win.Width/50; x++ {
		for y := 0; y < win.Height/50; y++ {

			// create star structs
			stars = append(stars, &Star{
				X:    float32(x - 5),
				Y:    float32(y - 4),
				Size: rand.Float32()/2 + 0.5,
				// Size:     1,
				SpriteId: spriteloader.GetSpriteIdByName(fmt.Sprintf("stars-%d", rand.Intn(16))),
			})
		}
	}

	gl.UseProgram(program)
	//main loop
	for !window.ShouldClose() {
		//clearing buffers
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// drawing stars
		for _, star := range stars {
			spriteloader.DrawSpriteQuad(star.X, star.Y, 1, 1, star.SpriteId)
		}
		// spriteloader.DrawSpriteQuad(float32(wx), float32(wy), 1, 1, 0)

		//polling events and swapping window buffers
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func keyCallback(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, m glfw.ModifierKey) {
	if a != glfw.Press {
		return
	}
	if k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	switch k {
	case glfw.KeyUp:
		wy += 1
	case glfw.KeyDown:
		wy -= 1
	case glfw.KeyLeft:
		wx -= 1
	case glfw.KeyRight:
		wx += 1
	case glfw.KeyX:
		size -= 0.1
	case glfw.KeyZ:
		size += 0.1
	case glfw.KeyL:
		spriteId += 1
	case glfw.KeyJ:
		spriteId -= 1
	case glfw.KeyTab:
		shuffle()
	}
	// log.Printf("wx: %f, wy: %f", wx, wy)

}

func initArgs() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.Float64Flag{
			Name:        "x",
			Usage:       "x position for a star",
			Destination: &wx,
		},
		&cli.Float64Flag{
			Name:        "y",
			Usage:       "y position for a star",
			Destination: &wy,
		},
		&cli.Float64Flag{
			Name:        "size",
			Usage:       "size for a star (default = 1)",
			Value:       1,
			Destination: &size,
		},
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
	fmt.Printf("%f %f %f\n", wx, wy, wz)
}

func shuffle() {
	for _, star := range stars {
		star.SpriteId = spriteloader.GetSpriteIdByName(fmt.Sprintf("stars-%d", rand.Intn(15)))
	}
}
