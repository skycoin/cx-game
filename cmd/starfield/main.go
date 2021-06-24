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

// A/D for moving

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
	Intensity_Period        int
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
	Intensity     float32
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

	//stop auto moving stars
	pauseAutoMove bool

	//variable for star move speed
	speed        float32 = 25
	rightPressed bool
	leftPressed  bool
)

func main() {
	//parse command line arguments and flags
	initArgs()
	//start goroutines to check reloaded files and update 	configurations
	initReloadConfig()

	// initialize both glfw and gl libraries, setting up the window and shader program
	window := render.NewWindow(cliConfig.Width, cliConfig.Height, true)
	defer glfw.Terminate()
	glfwWindow := window.Window
	glfwWindow.SetKeyCallback(keyCallback)
	shader := utility.NewShader("./cmd/starfield/shaders/main/vertex.glsl", "./cmd/starfield/shaders/main/fragment.glsl")
	shader.Use()

	spriteloader.InitSpriteloader(&window)
	spriteloader.DEBUG = false

	//ENABLE BLENDING
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// ortho := mgl32.Ortho2D(0, float32(cliConfig.Width), 0, float32(cliConfig.Height))
	// shader.SetMat4("ortho", &ortho)
	//init  stars
	initStarField()

	if cliConfig.Background == 1 {
		starmap.SettingsFile = "./cmd/starfield/config/starmap.yaml"
		starmap.Init(&window)
		starmap.Generate(256, 0.08, 3)
	}

	//bind gradient 1d textures
	for i := 1; i < 12; i++ {
		tex := getGradient(uint(i))
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_1D, tex)
	}

	lastFrame := float32(glfw.GetTime())
	dt := float32(0)

	//set up its own projection
	ortho := mgl32.Ortho2D(0, float32(window.Width), 0, float32(window.Height))
	shader.SetMat4("projection", &ortho)

	//main loop
	for !glfwWindow.ShouldClose() {
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

		updateStarField(dt)
		drawStarField(shader)

		glfw.PollEvents()
		glfwWindow.SwapBuffers()
	}
}

//create random stars
func initStarField() {

	perlinMap = genPerlin(cliConfig.Starfield_Width, cliConfig.Starfield_Height, noiseConfig)

	file1 := "./assets/starfield/stars/planets.png"
	file2 := file1
	// file1 := "cmd/starfield/stars_1.png"
	// file2 := "cmd/starfield/stars_2.png"
	star1SpriteSheetId := spriteloader.LoadSpriteSheet(file1)
	star2SpriteSheetId := spriteloader.LoadSpriteSheet(file2)

	//load all sprites from spritesheet
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			//for stars
			spriteloader.LoadSprite(star1SpriteSheetId,
				fmt.Sprintf("stars-1-%d", y*4+x),
				x, y,
			)
			spriteloader.LoadSprite(star2SpriteSheetId,
				fmt.Sprintf("stars-2-%d", y*4+x),
				x, y,
			)
		}
	}
	gaussianDepth := rand.Float32()
	for i := 0; i < cliConfig.StarAmount; i++ {
		spriteName := fmt.Sprintf("stars-%d-%d", rand.Intn(2)+1, rand.Intn(16))
		star := &Star{
			//bad generation position, TODO
			SpriteId:  spriteloader.GetSpriteIdByName(spriteName),
			Intensity: rand.Float32(),
			Depth:     rand.Float32(),
			Size:      float32(starConfig.Pixel_Size) * getRandomSize(),
		}

		if i < gaussianAmount {
			star.X, star.Y = getStarPosition(true, 1)
			star.GradientValue = 0.9
			star.Depth = gaussianDepth
			star.IsGaussian = true
		} else {
			star.X, star.Y = getStarPosition(false, 1)
			star.GradientValue = rand.Float32()
			star.IsGaussian = false
		}

		stars = append(stars, star)
	}

}

//regenerates positions of stars
func regenStarField() {
	gaussianAmount = cliConfig.StarAmount * starConfig.Gaussian_Percentage / 100
	for _, star := range stars {
		star.X, star.Y = getStarPosition(star.IsGaussian, 1)
		star.Size = float32(starConfig.Pixel_Size) * getRandomSize()
	}
}

//updates position of each star at each game iteration
func updateStarField(dt float32) {

	for _, star := range stars {
		// star.X = float32(math.Mod(float64(star.X+(speed*dt)), coords*2))
		if rightPressed {
			star.X += speed * dt * star.Depth
		}
		if leftPressed {
			star.X -= speed * dt * star.Depth
		}
		if star.X > float32(cliConfig.Starfield_Width) {
			star.X = 0
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
		shader.SetFloat("intensity", getIntensity(star.Intensity))
		spriteloader.DrawSpriteQuadCustom(star.X, star.Y, star.Size, star.Size, star.SpriteId, shader.ID)
	}

}

func getIntensity(intensity float32) float32 {

	// t := time.Now()

	// return myclamp(float32(math.Sin(float64(time.Now().UnixNano()))))
	ttime := float64(time.Now().UnixNano() / (10000000 * int64(starConfig.Intensity_Period)))

	rad := math.Pi / 180 * ttime

	sinusoid := math.Sin(rad * float64(intensity))

	result := float32(sinusoid)*0.55 + 0.45

	return float32(result)
}
func getRandomSize() float32 {
	return (1 + 3*rand.Float32()/2)
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
func getStarPosition(isGaussian bool, depth uint) (float32, float32) {

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
			x := rand.Intn(cliConfig.Starfield_Width)
			y := rand.Intn(cliConfig.Starfield_Height)
			perlinProb := perlinMap[x][y]
			deleteChance := (1 - perlinProb)
			if deleteChance > 0.5 {
				continue
			}
			// xPos, yPos = convertIntToSpriteCoords(x, y, 0, cliConfig.Width, 0, cliConfig.Height)
			// fmt.Println(x, "         ", y)
			xPos, yPos = float32(x), float32(y)
			break
		}
	}

	if starConfig.Star_Separation_Enabled > 0 {
		if starPositionInConflict(xPos, yPos) && depth < 5 {
			return getStarPosition(isGaussian, depth+1)
		}
	}
	return xPos, yPos
}

func starPositionInConflict(xPos, yPos float32) bool {

	mindistance := 10.0

	for _, star := range stars {
		if math.Abs(float64(star.X-xPos)) < mindistance ||
			math.Abs(float64(star.Y-yPos)) < mindistance {
			return true
		}
	}
	return false

}

func convertToWorldCoords(xPos, yPos, minX, maxX, minY, maxY float32) (float32, float32) {
	diffX := maxX - minX
	diffY := maxY - minY
	b := (xPos - minX) / diffX
	c := (yPos - minY) / diffY
	x := float32(cliConfig.Starfield_Width) * b
	y := float32(cliConfig.Starfield_Height) * c
	return x, y
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
			result = utility.ClampF(result-min/(max-min), 0.0, 1.0)
			grid[y] = append(grid[y], result)
		}

	}

	return grid
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
	if A < 0.6 {
		A = 0.6
	}
	theta = float64(utility.DegToRad(float32(starConfig.Gaussian_Angle)))
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
				saveResult()
			}
		case glfw.KeyP:
			pauseAutoMove = !pauseAutoMove
		case glfw.KeyA:
			leftPressed = true
		case glfw.KeyLeft:
			leftPressed = true
		case glfw.KeyD:
			rightPressed = true
		case glfw.KeyRight:
			rightPressed = true
		}

	} else if a == glfw.Release {
		switch k {
		case glfw.KeyA:
			leftPressed = false
		case glfw.KeyLeft:
			leftPressed = false
		case glfw.KeyD:
			rightPressed = false
		case glfw.KeyRight:
			rightPressed = false
		}
	}

}

func saveResult() {

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

//yaml reload function
func initReloadConfig() {
	starConfigReloaded := make(chan struct{})
	perlinConfigReloaded := make(chan struct{})
	done := make(chan struct{})
	go utility.CheckAndReload("./cmd/starfield/config/star.yaml", &starConfig, starConfigReloaded)
	go utility.CheckAndReload("./cmd/starfield/config/perlin.yaml", &noiseConfig, perlinConfigReloaded)
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
	<-done
	<-done
	close(done)
}
