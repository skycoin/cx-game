package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/starmap"
	"github.com/skycoin/cx-game/utility"
	"github.com/urfave/cli/v2"
)

//Press TAB to shuffle stars

func init() {
	// seed rand so stars will be random each program run
	rand.Seed(time.Now().UnixNano())
	//lock thread so drawing will be only in main thread, otherwise there will be errors
	runtime.LockOSThread()
}

type noiseSettings struct {
	Size     int
	Scale    float32
	Levels   uint8
	Contrast float32

	Seed        int64
	Gradmax     int
	X           int
	Xs          int
	Persistance float32
	Lacunarity  float32
	Octaves     int

	GradFile string
}

//via yaml set pixel_size
type starSettings struct {
	Pixel_Size int

	Gaussian_Percentage int
	Gaussian_Angle      int
	Gaussian_Offset_X   float32
	Gaussian_Offset_Y   float32
	Gaussian_Sigma_X    float32
	Gaussian_Sigma_Y    float32
	Gaussian_Constant   float32
}

//via cli set bg, star_amount, window width and height
type cliSettings struct {
	Background          int
	StarAmount          int
	Width               int
	Height              int
	Gaussian_Percentage int
}

type Star struct {
	X             float32
	Y             float32
	Size          float32
	SpriteId      int
	GradientValue float32
	GradientId    int32
	Depth         float32
	IsGaussian    bool
}

var (
	stars []*Star

	perlinMap = genPerlin(cliConfig.Width, cliConfig.Height, noiseConfig)
	//cli options
	cliConfig *cliSettings = &cliSettings{
		StarAmount:          500,
		Background:          0,
		Width:               800,
		Height:              600,
		Gaussian_Percentage: 45,
	}
	gaussianAmount int
	//star options (pixelsize)
	starConfig *starSettings = &starSettings{
		Pixel_Size:          1,
		Gaussian_Percentage: 25,
		Gaussian_Angle:      45,
		Gaussian_Offset_X:   0,
		Gaussian_Offset_Y:   0,
		Gaussian_Sigma_X:    0.3,
		Gaussian_Sigma_Y:    0.2,
		Gaussian_Constant:   1,
	}

	//perlin options
	noiseConfig *noiseSettings = &noiseSettings{
		Size:     1024,
		Scale:    0.04,
		Levels:   8,
		Contrast: 1.0,

		Seed:        1,
		X:           512,
		Xs:          4,
		Gradmax:     256,
		Persistance: 0.5,
		Lacunarity:  2,
		Octaves:     8,
	}

	starConfigReloaded   chan struct{}
	perlinConfigReloaded chan struct{}
)

func main() {
	starConfigReloaded = make(chan struct{})
	perlinConfigReloaded = make(chan struct{})
	go utility.CheckAndReload("./cmd/starfield/config/config.yaml", starConfig, starConfigReloaded)
	go utility.CheckAndReload("./cmd/starfield/config/perlin.yaml", starConfig, perlinConfigReloaded)
	go func() {
		for {
			select {
			case <-starConfigReloaded:
				regenStarField()
			case <-perlinConfigReloaded:
				perlinMap = genPerlin(cliConfig.Width, cliConfig.Height, noiseConfig)
				regenStarField()
			}
		}
	}()
	//parse command line arguments and flags
	initArgs()
	gaussianAmount = cliConfig.StarAmount * cliConfig.Gaussian_Percentage / 100
	// initialize both glfw and gl libraries, setting up the window and shader program
	win := render.NewWindow(cliConfig.Height, cliConfig.Width, true)
	defer glfw.Terminate()
	window := win.Window
	window.SetKeyCallback(keyCallback)
	shader := utility.NewShader("./cmd/starfield/shaders/vertex.glsl", "./cmd/starfield/shaders/fragment.glsl")
	shader.Use()
	ortho := mgl32.Ortho2D(0, float32(cliConfig.Width), 0, float32(cliConfig.Height))
	shader.SetMat4("ortho", &ortho)
	//randomize stars
	initStarField(&win)

	if cliConfig.Background == 1 {
		starmap.Init(&win)
		starmap.Generate(256, 0.08, 3)
	}

	//bind gradient 1d textures
	for i := 1; i < 12; i++ {
		tex := getGradient(uint(i))
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_1D, tex)
	}

	shader.Use()

	//main loop
	for !window.ShouldClose() {
		//clearing buffers
		gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if cliConfig.Background == 1 {
			starmap.Draw()
		}
		drawStarField(shader)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

//create random stars
func initStarField(win *render.Window) {
	//spriteloader init
	spriteloader.InitSpriteloader(win)
	star1SpriteSheetId := spriteloader.LoadSpriteSheet("./cmd/starfield/stars_1.png")
	star2SpriteSheetId := spriteloader.LoadSpriteSheet("./cmd/starfield/stars_2.png")

	//load all sprites from spritesheet
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			//for stars
			spriteloader.LoadSprite(star1SpriteSheetId,
				fmt.Sprintf("stars-1-%d", y*4+x),
				x, y,
			)
			spriteloader.LoadSprite(star2SpriteSheetId,
				fmt.Sprintf("stars-2-%d", y+4*x),
				x, y,
			)
		}
	}
	for i := 0; i < cliConfig.StarAmount; i++ {
		spriteName := fmt.Sprintf("stars-%d-%d", rand.Intn(2)+1, rand.Intn(16))
		star := &Star{
			//bad generation position, TODO
			SpriteId: spriteloader.GetSpriteIdByName(spriteName),
		}
		if i < gaussianAmount {
			star.IsGaussian = true
			star.X, star.Y = getStarPosition(true)
			star.GradientValue = 0.9
		} else {
			star.X, star.Y = getStarPosition(false)
			star.GradientValue = rand.Float32()
			star.IsGaussian = false
		}
		star.Size = float32(starConfig.Pixel_Size) / 32 * (1 + rand.Float32()/2)
		stars = append(stars, star)
	}
}

func regenStarField() {
	for i, star := range stars {
		if i < gaussianAmount {
			star.X, star.Y = getStarPosition(true)
		} else {
			star.X, star.Y = getStarPosition(false)
		}
		star.Size = float32(starConfig.Pixel_Size) / 32 * (1 + rand.Float32()/2)
	}
}

func drawStarField(shader *utility.Shader) {
	shader.Use()
	for _, star := range stars {
		if star.IsGaussian {
			shader.SetInt("texture_1d", 1)
		} else {
			shader.SetInt("texture_1d", 4)
			shader.SetFloat("gradValue", star.GradientValue)
		}
		spriteloader.DrawSpriteQuadCustom(star.X, star.Y, star.Size, star.Size, star.SpriteId, shader.ID)
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
			Value:       cliConfig.Background,
			Destination: &cliConfig.Background,
		},
		&cli.IntFlag{
			Name:        "stars",
			Aliases:     []string{"star"},
			Usage:       "number of stars to draw",
			Value:       cliConfig.StarAmount,
			Destination: &cliConfig.StarAmount,
		},
		&cli.IntFlag{
			Name:        "width",
			Usage:       "Resolution width",
			Value:       cliConfig.Width,
			Destination: &cliConfig.Width,
		},
		&cli.IntFlag{
			Name:        "height",
			Usage:       "Resolution height",
			Value:       cliConfig.Height,
			Destination: &cliConfig.Height,
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
	regenStarField()
}

func getSize() float32 {
	size := rand.Float32()/2 + 0.75
	if size > 0.5 && size < 0.75 {
		size = rand.Float32() / 4
	}
	return size
}

//get gradient file
func getGradient(gradientNumber uint) uint32 {
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

//todo
func getStarPosition(isGaussian bool) (float32, float32) {

	var xPos, yPos float32
	if isGaussian {
		for {
			xPos, yPos = rand.Float32(), rand.Float32()
			z := gaussianTheta(xPos, yPos)
			if 1-z > 0.5*rand.Float32() {
				continue
			}
			xPos, yPos = convertToSpriteCoords(xPos, yPos, 0, 1, 0, 1)
			break
		}
	} else {

		for {
			x := rand.Intn(cliConfig.Width)
			y := rand.Intn(cliConfig.Height)
			perlinProb := perlinMap[x][y]
			deleteChance := (1 - perlinProb)
			if deleteChance > 0.5 {
				continue
			}
			xPos, yPos = convertIntToSpriteCoords(x, y, 0, cliConfig.Width, 0, cliConfig.Height)
			break
		}
	}

	return xPos, yPos
	// xPos, yPos = convertIntCoords(x, y, 0, cliConfig.Width, 0, cliConfig.Height)

	// starGap := 0.7
	// xPos, yPos := rand.Intn(cliConfig.Width), rand.Intn(cliConfig.Height)
	// return float32(xPos), float32(yPos)
	//if too many stars
	// if cliConfig.StarAmount > 20 || starGap > 1.3 {
	// 	return xPos, yPos
	// }
	// for _, star := range stars {
	// 	if math.Abs(float64(xPos-star.X)) < float64(starGap) && math.Abs(float64(yPos-star.Y)) < float64(starGap) {
	// 		return getStarPosition()
	// 	}
	// }

}

//Convert from any coordinates specified min max to screen coordinates
func convertIntToSpriteCoords(xPos, yPos, minX, maxX, minY, maxY int) (float32, float32) {
	return convertToSpriteCoords(
		float32(xPos),
		float32(yPos),
		float32(minX),
		float32(maxX),
		float32(minY),
		float32(maxY),
	)
}

//Convert from any coordinates specified min max to screen coordinates
func convertToSpriteCoords(xPos, yPos, minX, maxX, minY, maxY float32) (float32, float32) {
	// for 800/600 aspect
	// bottom left = -5.4, -4
	// top left =  -5.4, 4
	// top right = 5.4, 4
	// bottom right = 5.4, -4
	diffX := maxX - minX
	diffY := maxY - minY
	b := (xPos - minX) / diffX
	c := (yPos - minY) / diffY
	x := 10.8*b - 5.4
	y := 8*c - 4

	return x, y
}
func genPerlin(width, height int, noiseConfig *noiseSettings) [][]float32 {
	grid := make([][]float32, 0)
	myPerlin := perlin.NewPerlin2D(
		noiseConfig.Seed,
		noiseConfig.X,
		noiseConfig.Xs,
		noiseConfig.Gradmax,
	)
	max := float32(math.Sqrt2 / (1.9 * noiseConfig.Contrast))
	min := float32(-math.Sqrt2 / (1.9 * noiseConfig.Contrast))
	for y := 0; y < width; y++ {
		grid = append(grid, []float32{})
		for x := 0; x < height; x++ {
			result := myPerlin.Noise(
				float32(x)*noiseConfig.Scale,
				float32(y)*noiseConfig.Scale,
				noiseConfig.Persistance,
				noiseConfig.Lacunarity,
				noiseConfig.Octaves,
			)
			result = clamp(result-min/(max-min), 0.0, 1.0)
			grid[y] = append(grid[y], result)
		}

	}

	return grid
}

func clamp(number, min, max float32) float32 {
	if number > max {
		return max
	}
	if number < min {
		return min
	}
	return number
}

func gaussianTheta(x32, y32 float32) float32 {
	x, y := float64(x32), float64(y32)
	var sigmaX, sigmaY, x0, y0, A, theta float64
	// var A float64 = 1.9
	// var theta float64
	// sigmaX = 0.1
	// sigmaY = 0.3
	// x0 = 0.5
	// y0 = 0.5
	A = float64(starConfig.Gaussian_Constant)
	theta = float64(mgl32.DegToRad(float32(starConfig.Gaussian_Angle)))
	sigmaX = float64(starConfig.Gaussian_Sigma_X)
	sigmaY = float64(starConfig.Gaussian_Sigma_Y)
	x0 = float64(starConfig.Gaussian_Offset_X)
	y0 = float64(starConfig.Gaussian_Offset_Y)

	a := math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaY, 2))
	b := -math.Sin(2*theta)/(4*math.Pow(sigmaX, 2)) + math.Sin(2*theta)/(4*math.Pow(sigmaY, 2))
	c := math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaY, 2))
	result := A * math.Exp(-(a*math.Pow(x-x0, 2) + 2*b*(x-x0)*(y-y0) + c*math.Pow(y-y0, 2)))
	return float32(result)
}

func DegToRad(angle float32) float32 {
	return math.Pi / 180 * angle
}
