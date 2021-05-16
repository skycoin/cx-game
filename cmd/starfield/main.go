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
	Pixel_Size              int
	Star_Separation_Enabled int
	Gaussian_Percentage     int
	Gaussian_Angle          int
	Gaussian_Offset_X       float32
	Gaussian_Offset_Y       float32
	Gaussian_Sigma_X        float32
	Gaussian_Sigma_Y        float32
	Gaussian_Constant       float32
}

//via cli set bg, star_amount, window width and height
type cliSettings struct {
	Background       int
	StarAmount       int
	Width            int
	Height           int
	Starfield_Width  int
	Starfield_Height int
}

type Star struct {
	X             float32
	Y             float32
	Size          float32
	SpriteId      int
	GradientValue float32
	GradientId    int
	Depth         float32
	IsGaussian    bool
}

var (
	stars []*Star

	perlinMap [][]float32
	//cli options
	cliConfig *cliSettings = &cliSettings{
		StarAmount:       500,
		Background:       0,
		Width:            800,
		Height:           600,
		Starfield_Width:  1200,
		Starfield_Height: 900,
	}
	gaussianAmount int
	//star options (pixelsize)
	starConfig *starSettings = &starSettings{
		// Pixel_Size:            1,
		// StarSeparationEnabled: false,
		// Gaussian_Percentage:   25,
		// Gaussian_Angle:        45,
		// Gaussian_Offset_X:     0,
		// Gaussian_Offset_Y:     0,
		// Gaussian_Sigma_X:      0.3,
		// Gaussian_Sigma_Y:      0.2,
		// Gaussian_Constant:     1,
	}

	//perlin options
	noiseConfig *noiseSettings = &noiseSettings{
		// Size:     1024,
		// Scale:    0.04,
		// Levels:   8,
		// Contrast: 1.0,

		// Seed:        1,
		// X:           512,
		// Xs:          4,
		// Gradmax:     256,
		// Persistance: 0.5,
		// Lacunarity:  2,
		// Octaves:     8,
	}

	dt, lastFrame float32

	starConfigReloaded   chan struct{} = make(chan struct{})
	perlinConfigReloaded chan struct{} = make(chan struct{})

	done chan struct{} = make(chan struct{})

	//stop auto moving stars
	autoMove bool

	//variable for star move speed
	speed float32 = 0.25
)

func main() {
	//parse command line arguments and flags
	initArgs()

	//goroutines to check reloaded files and update 	configurations
	go utility.CheckAndReload("./cmd/starfield/config/config.yaml", starConfig, starConfigReloaded)
	go utility.CheckAndReload("./cmd/starfield/config/perlin.yaml", noiseConfig, perlinConfigReloaded)
	go func() {
		var isStarConfigFirstLoad, isPerlinConfigFirstLoad bool = true, true
		for {
			select {
			case <-starConfigReloaded:
				gaussianAmount = cliConfig.StarAmount * starConfig.Gaussian_Percentage / 100
				if isStarConfigFirstLoad {
					done <- struct{}{}
					isStarConfigFirstLoad = false
				} else {
					regenStarField()
				}

			case <-perlinConfigReloaded:

				if isPerlinConfigFirstLoad {
					done <- struct{}{}
					isPerlinConfigFirstLoad = false
				} else {
					regenStarField()
				}
			}
		}
	}()

	// initialize both glfw and gl libraries, setting up the window and shader program
	win := render.NewWindow(cliConfig.Width, cliConfig.Height, true)
	defer glfw.Terminate()
	window := win.Window
	window.SetKeyCallback(keyCallback)
	shader := utility.NewShader("./cmd/starfield/shaders/main/vertex.glsl", "./cmd/starfield/shaders/main/fragment.glsl")
	shader.Use()
	// ortho := mgl32.Ortho2D(0, float32(cliConfig.Width), 0, float32(cliConfig.Height))
	// shader.SetMat4("ortho", &ortho)
	//init  stars
	initStarField(&win)

	if cliConfig.Background == 1 {
		starmap.SettingsFile = "./cmd/starfield/config/starmap.yaml"
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
		//deltatime
		currFrame := float32(glfw.GetTime())
		dt = currFrame - lastFrame
		lastFrame = currFrame

		//clearing buffers
		gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if cliConfig.Background == 1 {
			starmap.Draw()
		}
		updateStarField()
		drawStarField(shader)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

//create random stars
func initStarField(win *render.Window) {
	<-done
	<-done
	close(done)
	//spriteloader init
	perlinMap = genPerlin(cliConfig.Starfield_Width, cliConfig.Starfield_Height, noiseConfig)

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
	gaussianDepth := rand.Float32()
	for i := 0; i < cliConfig.StarAmount; i++ {
		spriteName := fmt.Sprintf("stars-%d-%d", rand.Intn(2)+1, rand.Intn(16))
		star := &Star{
			//bad generation position, TODO
			SpriteId: spriteloader.GetSpriteIdByName(spriteName),
		}

		star.Size = float32(starConfig.Pixel_Size) / 32 * (3 * rand.Float32() / 2)
		star.Depth = rand.Float32()
		if i < gaussianAmount {
			star.IsGaussian = true
			star.X, star.Y = getStarPosition(true)
			star.GradientValue = 0.9
			star.Depth = gaussianDepth
		} else {
			star.X, star.Y = getStarPosition(false)
			star.GradientValue = rand.Float32()
			star.IsGaussian = false
		}

		stars = append(stars, star)
	}
}

//regenerates positions of stars
func regenStarField() {
	gaussianAmount = cliConfig.StarAmount * starConfig.Gaussian_Percentage / 100
	for i, star := range stars {
		if i < gaussianAmount {
			star.X, star.Y = getStarPosition(true)
			star.IsGaussian = true
		} else {
			star.X, star.Y = getStarPosition(false)
			star.IsGaussian = false
		}
		star.Size = float32(starConfig.Pixel_Size) / 32 * (3 * rand.Float32() / 2)
	}
}

//updates position of each star at each game iteration
func updateStarField() {
	if autoMove {
		return
	}

	coords := float64(cliConfig.Starfield_Width) / float64(cliConfig.Width) * 5.333
	for _, star := range stars {
		// star.X = float32(math.Mod(float64(star.X+(speed*dt)), coords*2))
		star.X += speed * dt * star.Depth
		if star.X > float32(coords) {
			star.X = -float32(coords)
		}
		// star.X = star.X + dt
	}
}

//draws stars via drawquad function
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
func keyCallback(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press {
		switch k {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		case glfw.KeyTab:
			shuffle()
		case glfw.KeyS:
			if mk == glfw.ModShift {
				fmt.Println("saved!")
			}
		case glfw.KeyP:
			autoMove = !autoMove
		}
	} else {
		switch k {
		case glfw.KeyKPAdd:
			speed += 0.05
		case glfw.KeyKPSubtract:
			speed -= 0.05
		}
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
		&cli.IntFlag{
			Name:        "swidth",
			Aliases:     []string{"starfield_width", "starfieldw"},
			Usage:       "Starfield width",
			Value:       cliConfig.Starfield_Width,
			Destination: &cliConfig.Starfield_Width,
		},
		&cli.IntFlag{
			Name:        "sheight",
			Aliases:     []string{"starfield_height", "starfieldh"},
			Usage:       "Starfield height",
			Value:       cliConfig.Starfield_Width,
			Destination: &cliConfig.Starfield_Height,
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

//get gradient file given id from 1-11 (see assets folder)
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

// Get Random Position for a star
//  isGaussian - boolean to specify is star should be in a gaussian band
func getStarPosition(isGaussian bool) (float32, float32) {

	var xPos, yPos float32
	if isGaussian {
		for {
			xPos, yPos = rand.Float32(), rand.Float32()
			z := gaussianTheta(xPos, yPos)
			//randomness
			if 1-z > 0.5*rand.Float32() {
				continue
			}
			xPos, yPos = convertToWorldCoords(xPos, yPos, 0, 1, 0, 1)
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

	if starConfig.Star_Separation_Enabled > 0 {
		if starPositionInConflict(xPos, yPos) {
			// return getStarPosition(isGaussian)
		}
	}
	return xPos, yPos
}

// Checks whether the position is within minimum distance of other stars
var counter int

func starPositionInConflict(xPos, yPos float32) bool {
	mindistance := 0.3

	for _, star := range stars {
		if math.Abs(float64(star.X-xPos)) < mindistance ||
			math.Abs(float64(star.Y-yPos)) < mindistance {
			return true
		}
	}
	return false
}

//Convert from any coordinates specified min max to screen coordinates
func convertIntToSpriteCoords(xPos, yPos, minX, maxX, minY, maxY int) (float32, float32) {
	return convertToWorldCoords(
		float32(xPos),
		float32(yPos),
		float32(minX),
		float32(maxX),
		float32(minY),
		float32(maxY),
	)
}

//Convert from any coordinates specified min max to screen coordinates
//  xPos - X position in original coordinates
//  yPos - Y position
//  minX - lower X bound
//  maxX - upper X bound
//  minY - lower Y bound
//  maxY - upper Y bound
func convertToWorldCoords(xPos, yPos, minX, maxX, minY, maxY float32) (float32, float32) {
	// for 800/600 aspect
	// bottom left = -5.4`QQ, -4
	// top left =  -5.4, 4
	// top right = 5.4, 4
	// bottom right = 5.4, -4
	diffX := maxX - minX
	diffY := maxY - minY
	b := (xPos - minX) / diffX
	c := (yPos - minY) / diffY
	x := 10.8*b - 5.4
	y := 8*c - 4

	//multipliers for star
	xMul := float32(cliConfig.Starfield_Width) / float32(cliConfig.Width)
	yMul := float32(cliConfig.Starfield_Height) / float32(cliConfig.Height)
	return x * xMul, y * yMul
}

// Generate perlin 2d slice for each pixel on the window.
//  width - width of the window
//  height - height of the window
//  noiseConfig - config struct with parameters
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

// Clamp values to edges
//  number - value to clamp
//  min - lower bound
//  max - upper bound
func clamp(number, min, max float32) float32 {
	if number > max {
		return max
	}
	if number < min {
		return min
	}
	return number
}

// Gaussian bivariate distribution function
//  x32 - X position (float32)
//  y32 - Y position (float32)
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

// Convenient function to convert angle in degrees to radians
//  angle - angle in degrees
func DegToRad(angle float32) float32 {
	return math.Pi / 180 * angle
}
