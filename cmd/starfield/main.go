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
	"github.com/skycoin/cx-game/starmap"
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

var (
	// wx, wy, wz      float64 = 0, 0, -10
	// size            float64 = 1
	// spriteId        int
	// stars           []*Star
	backgroundStars []*Star

	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}

	//cli options
	background int = 1 //0 is black, 1 is rgb
	starAmount int = 5
	width      int = 800
	height     int = 600
)

func main() {
	//parse command line arguments and flags
	initArgs()

	// initialize both glfw and gl libraries, setting up the window and shader program
	win := render.NewWindow(height, width, true)
	defer glfw.Terminate()
	window := win.Window
	window.SetKeyCallback(keyCallback)
	program := win.Program
	gl.UseProgram(program)

	if background == 1 {
		//expects to create and bind vertex array object for some reason, otherwise will fall
		vao := makeVao()
		gl.BindVertexArray(vao)

		starmap.Init(&win)
		starmap.Generate(256, 0.08, 3)
	}

	//randomize stars
	initStarField(&win)

	//main loop
	for !window.ShouldClose() {
		//clearing buffers
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if background == 1 {
			starmap.Draw()
		}

		//has to be after starmap, otherwise starmap will be drawn over the stars
		drawStarField(&win)

		//polling events and swapping window buffers
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

//callback function to register key events
func keyCallback(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, m glfw.ModifierKey) {
	if a != glfw.Press {
		return
	}
	if k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	switch k {
	// case glfw.KeyUp:
	// 	wy += 1
	// case glfw.KeyDown:
	// 	wy -= 1
	// case glfw.KeyLeft:
	// 	wx -= 1
	// case glfw.KeyRight:
	// 	wx += 1
	// case glfw.KeyX:
	// 	size -= 0.1
	// case glfw.KeyZ:
	// 	size += 0.1
	// case glfw.KeyL:
	// 	spriteId += 1
	// case glfw.KeyJ:
	// 	spriteId -= 1
	case glfw.KeyTab:
		shuffle()
	}
	// log.Printf("wx: %f, wy: %f", wx, wy)

}

//function to parse cli flags
func initArgs() {

	app := cli.NewApp()
	app.Name = "starfield-cli"
	app.Description = "starfield example"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "background",
			Aliases:     []string{"bg", "b"},
			Usage:       "background to use",
			Value:       0,
			Destination: &background,
		},
		&cli.IntFlag{
			Name:        "stars",
			Aliases:     []string{"star"},
			Usage:       "number of stars to draw",
			Value:       5,
			Destination: &starAmount,
		},
		&cli.IntFlag{
			Name:        "width",
			Usage:       "Resolution width",
			Value:       800,
			Destination: &width,
		},
		&cli.IntFlag{
			Name:        "height",
			Usage:       "Resolution height",
			Value:       600,
			Destination: &height,
		},
	}
	app.After = func(c *cli.Context) error {
		command := c.Args().First()
		if command == "help" {
			os.Exit(0)
		}
		return nil
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
	// fmt.Printf("%f %f %f\n", wx, wy, wz)
}

//silly function to shuffle stars on the background
func shuffle() {
	for _, star := range backgroundStars {
		star.SpriteId = spriteloader.GetSpriteIdByName(fmt.Sprintf("background-stars-%d", rand.Intn(15)))
		star.Size = getSize()
	}
}

//create vertex array object
func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(sprite), gl.Ptr(sprite), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}

//create random stars
func initStarField(win *render.Window) {
	//spriteloader init
	spriteloader.InitSpriteloader(win)
	starSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/starfield_test_16x16_tiles_8x8_tile_grid_128x128.png")
	// galaxySheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/galaxy_256x256.png")

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			spriteloader.LoadSprite(starSheetId,
				fmt.Sprintf("background-stars-%d", y*4+x),
				x, y)
		}
	}
	// //load all sprites from spritesheet
	// for y := 0; y < 6; y++ {
	// 	for x := 0; x < 8; x++ {

	// 		//for stars
	// 		spriteloader.LoadSprite(galaxySheetId,
	// 			fmt.Sprintf("galaxy-stars-%d", y*8+x),
	// 			x, y,
	// 		)
	// 	}
	// }
	for x := 0; x < win.Width/60; x++ {
		for y := 0; y < win.Height/60; y++ {
			backgroundStars = append(backgroundStars, &Star{
				X:    float32(x - win.Width/120),
				Y:    float32(y - win.Height/120),
				Size: getSize(),
				// Size:     1,
				SpriteId: spriteloader.GetSpriteIdByName(fmt.Sprintf("background-stars-%d", rand.Intn(16))),
			})
		}
	}

	// for i := 0; i < starAmount; i++ {
	// 	stars = append(stars, &Star{
	// 		X:        rand.Float32()*15 - 8,
	// 		Y:        rand.Float32()*8 - 5,
	// 		Size:     1,
	// 		SpriteId: spriteloader.GetSpriteIdByName(fmt.Sprintf("galaxy-stars-%d", rand.Intn(48))),
	// 	})
	// }

}

func drawStarField(win *render.Window) {

	//draw star background

	for _, star := range backgroundStars {
		spriteloader.DrawSpriteQuad(star.X, star.Y, star.Size, star.Size, star.SpriteId)
	}

	// for _, star := range stars {
	// 	spriteloader.DrawSpriteQuad(star.X, star.Y, 1, 1, star.SpriteId)
	// }

}

func getSize() float32 {
	size := rand.Float32()/2 + 0.75
	if size > 0.5 && size < 0.75 {
		size = rand.Float32() / 4
	}
	return size
}

//TODO add shader 1d texture gradient
