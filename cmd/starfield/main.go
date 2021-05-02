package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/starmap"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
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
	Gradient float32
	Depth    float32
}
type Config struct {
	PixelSize float32
}

var (
	// wx, wy, wz      float64 = 0, 0, -10
	// size            float64 = 1
	// spriteId        int
	stars           []*Star
	backgroundStars []*Star

	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}

	gradValue = rand.Float32()
	//cli options
	background int = 1 //0 is black, 1 is rgb
	starAmount int = 15
	width      int = 800
	height     int = 600

	config *Config = &Config{1}
)

func main() {
	//parse command line arguments and flags
	initArgs()

	// initialize both glfw and gl libraries, setting up the window and shader program
	win := render.NewWindow(height, width, true)
	defer glfw.Terminate()
	window := win.Window
	window.SetKeyCallback(keyCallback)
	program1 := win.Program
	program2 := render.InitOpenGLCustom("./cmd/starfield/shaders/")

	if background == 1 {
		// vao := makeVao()
		// gl.BindVertexArray(vao)

		starmap.Init(&win)
		starmap.Generate(256, 0.08, 3)
	}

	go checkAndReload()
	//randomize stars
	initStarField(&win)

	//get gradient
	tex_2d := gen2DTexture()
	tex_1d := gen1DTexture(6)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_1D, tex_1d)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, tex_2d)
	// gl.ActiveTexture(gl.TEXTURE0)
	// vavao := spriteloader.MakeQuadVao()
	//main loop
	for !window.ShouldClose() {
		//clearing buffers
		if background == 2 {
			gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
		} else {
			gl.ClearColor(0.2, 0.3, 0.4, 1)
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if background == 1 {
			starmap.Draw()
		}
		//has to be after starmap, otherwise starmap will be drawn over the stars
		gl.UseProgram(program1)
		drawBackGroundStars()
		gl.UseProgram(program2)

		drawStars(program2)
		// world := mgl32.Ident4()
		// projection := mgl32.Ident4()
		// gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("projection\x00")), 1, false, &projection[0])
		// gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("world\x00")), 1, false, &world[0])
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_1D, tex_1d)
		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("ourTexture\x00")), 0)

		gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("texture_2d\x00")), 2)
		gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("texture_1d\x00")), 0)

		for _, star := range stars {
			gl.Uniform1f(gl.GetUniformLocation(program2, gl.Str("mixvalue\x00")), float32(0.3))
			// gl.ProgramUniform1f(program2, gl.GetUniformLocation(program2, gl.Str("gradvalue\x00")), gradValue)
			spriteloader.DrawSpriteQuadCustom(star.X, star.Y, 1, 1, star.SpriteId, program2)
		}

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		// // drawStars(program2)
		// gl.BindVertexArray(vavao)
		// gl.DrawArrays(gl.TRIANGLES, 0, 6)

		// gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("texture1\x00")), 0)

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
	case glfw.KeyTab:
		shuffle()
	}

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
}

//function to shuffle stars on the background
func shuffle() {
	for _, star := range backgroundStars {
		star.SpriteId = spriteloader.GetSpriteIdByName(fmt.Sprintf("background-stars-%d", rand.Intn(15)))
		star.Size = getSize()
	}
}

//create random stars
func initStarField(win *render.Window) {
	//spriteloader init
	spriteloader.InitSpriteloader(win)
	backgroundStarsheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/starfield_test_16x16_tiles_8x8_tile_grid_128x128.png")
	// starSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")
	// starSheetId := spriteloader.LoadSpriteSheet("./cmd/starfield/stars_2.png")
	galaxySheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/galaxy_256x256.png")

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			spriteloader.LoadSprite(backgroundStarsheetId,
				fmt.Sprintf("background-stars-%d", y*4+x),
				x, y)
		}
	}
	//load all sprites from spritesheet
	for y := 0; y < 6; y++ {
		for x := 0; x < 8; x++ {
			//for stars
			spriteloader.LoadSprite(galaxySheetId,
				fmt.Sprintf("stars-%d", y*8+x),
				x, y,
			)
		}
	}
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

	for i := 0; i < starAmount; i++ {
		stars = append(stars, &Star{
			X: rand.Float32()*8 - 4,
			Y: rand.Float32()*8 - 4,
			// X:        1,
			// Y:        1,
			Size:     1,
			SpriteId: spriteloader.GetSpriteIdByName(fmt.Sprintf("stars-%d", rand.Intn(48))),
			Gradient: rand.Float32(),
		})
	}

}

func drawBackGroundStars() {
	for _, star := range backgroundStars {
		// spriteloader.DrawSpriteQuad(star.X, star.Y, star.Size, star.Size, star.SpriteId)
		spriteloader.DrawSpriteQuad(star.X, star.Y, star.Size*(1+config.PixelSize/10), star.Size*(1+config.PixelSize/10), star.SpriteId)
	}
}
func drawStars(program uint32) {
	gl.UseProgram(program)
	for _, star := range stars {
		// gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("gradValue\x00")), star.Gradient)
		spriteloader.DrawSpriteQuadCustom(star.X, star.Y, 1, 1, star.SpriteId, program)
	}

}

func getSize() float32 {
	size := rand.Float32()/2 + 0.75
	if size > 0.5 && size < 0.75 {
		size = rand.Float32() / 4
	}
	return size
}

func checkAndReload() {
	configFilename := "./cmd/starfield/config.yaml"
	fileStat, err := os.Stat(configFilename)
	if err != nil {
		log.Panic(err)
	}

	for {
		newFileStat, err := os.Stat(configFilename)
		if err != nil {
			log.Panic(err)
		}
		//check if file is changed
		if newFileStat.ModTime() != fileStat.ModTime() || newFileStat.Size() != fileStat.Size() {
			data, err := ioutil.ReadFile(configFilename)
			if err != nil {
				log.Panic(err)
			}
			yaml.Unmarshal(data, config)
			fmt.Println(config)
			fileStat = newFileStat
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func gen1DTexture(gradientNumber uint) uint32 {
	var tex uint32

	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_1D, tex)

	result, img, _ := spriteloader.LoadPng(filepath.Join("./assets/starfield/gradients", fmt.Sprintf("heightmap_gradient_%02d.png", gradientNumber)))
	if result != 0 {
		log.Panic("Could not load picture!")
	}

	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA, int32(img.Rect.Size().X), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return tex
}

func gen2DTexture() uint32 {
	var tex uint32

	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	result, img, _ := spriteloader.LoadPng("./assets/sprite.png")
	if result != 0 {
		log.Panic("Could not load picture")
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X), int32(img.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return tex
}
