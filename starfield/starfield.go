package starfield

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/models"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
)

func init() {
	// seed rand so stars will be random each program run
	rand.Seed(time.Now().UnixNano())
}

var (
	p *models.Player
)

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

type starSettings struct {
	Star_size               int
	Star_Separation_Enabled int
	Gaussian_Percentage     int
	Gaussian_Angle          int
	Gaussian_Offset_X       float32
	Gaussian_Offset_Y       float32
	Gaussian_Sigma_X        float32
	Gaussian_Sigma_Y        float32
	Gaussian_Constant       float32
	Intensity_Period        float32
}

type StarFieldSettings struct {
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

	perlinMap    [][]float32
	perlinSpread float32 = 0.8
	//cli options
	windowConfig *StarFieldSettings = &StarFieldSettings{
		StarAmount: 1500,
	}
	gaussianAmount int
	//star options (pixelsize)
	starConfig *starSettings = &starSettings{
		Star_size:           2,
		Gaussian_Percentage: 70,
		Gaussian_Angle:      45,
		Gaussian_Offset_X:   0.3,
		Gaussian_Offset_Y:   0.2,
		Gaussian_Sigma_X:    0.4,
		Gaussian_Sigma_Y:    0.1,
		Gaussian_Constant:   1,
		Intensity_Period:    0.5,
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

	//variable for star move speed
	speed  float32 = 4
	shader *utility.Shader
)

//create random stars
func InitStarField(window *render.Window, player *models.Player) {
	p = player
	shader = utility.NewShader("./assets/shader/starfield/shader.vert", "./assets/shader/starfield/shader.frag")
	shader.Use()

	for i := 1; i < 12; i++ {
		tex := getGradient(uint(i))
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_1D, tex)
	}

	windowConfig.Width = window.Width
	windowConfig.Height = window.Height
	windowConfig.Starfield_Width = int(float32(windowConfig.Width) * 1.3)
	windowConfig.Starfield_Height = int(float32(windowConfig.Width) * 1.3)
	spriteloader.InitSpriteloader(window)
	spriteloader.DEBUG = false

	perlinMap = genPerlin(windowConfig.Starfield_Width, windowConfig.Starfield_Height, noiseConfig)

	file1 := "./assets/starfield/stars/planets.png"
	file2 := file1
	// file1 := "cmd/starfield/stars_1.png"
	// file2 := "cmd/starfield/stars_2.png"
	star1SpriteSheetId := spriteloader.LoadSpriteSheet(file1)
	star2SpriteSheetId := spriteloader.LoadSpriteSheet(file2)

	ortho := mgl32.Ortho2D(0, float32(window.Width), 0, float32(window.Height))
	shader.SetMat4("projection", &ortho)

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
	for i := 0; i < windowConfig.StarAmount; i++ {
		spriteName := fmt.Sprintf("stars-%d-%d", rand.Intn(2)+1, rand.Intn(16))
		star := &Star{
			//bad generation position, TODO
			SpriteId:  spriteloader.GetSpriteIdByName(spriteName),
			Intensity: rand.Float32(),
			Depth:     rand.Float32(),
			Size:      float32(starConfig.Star_size) * getRandomSize(),
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

//updates position of each star at each game iteration
func UpdateStarField(dt float32) {
	for _, star := range stars {
		if p.IsMoving("left") {
			star.X += speed * dt * star.Depth
		}
		if p.IsMoving("right") {
			star.X -= speed * dt * star.Depth
		}
		if p.IsMoving("up") {
			star.Y -= speed * dt * star.Depth
		}
		if p.IsMoving("down") {
			star.Y += speed * dt * star.Depth
		}
		if star.X > float32(windowConfig.Starfield_Width) {
			star.X = 0
		}
		star.X += 7 * dt * star.Depth * (rand.Float32() - 0.5)
	}

}

//draws stars via drawquad function
func DrawStarField() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ALPHA, gl.ONE_MINUS_SRC_ALPHA)
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
	ttime := float64(float32(time.Now().UnixNano()) / (10000000 * starConfig.Intensity_Period))

	rad := math.Pi / 180 * ttime

	sinusoid := math.Sin(rad * float64(intensity))

	result := float32(sinusoid)*0.55 + 0.45

	return float32(result)
}

func getRandomSize() float32 {
	return (1 + 3*rand.Float32()/2)
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

var max_depth = 5

// Get Random Position for a star
//  isGaussian - boolean to specify is star should be in a gaussian band
func getStarPosition(isGaussian bool, depth int) (float32, float32) {

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
			x := rand.Intn(windowConfig.Starfield_Width)
			y := rand.Intn(windowConfig.Starfield_Height)
			prob := perlinMap[x][y]
			deleteChance := (1 - prob)
			if deleteChance > perlinSpread {
				continue
			}
			xPos, yPos = float32(x), float32(y)
			break
		}
	}

	if starConfig.Star_Separation_Enabled > 0 {
		if starPositionInConflict(xPos, yPos) {
			if depth < max_depth {
				return getStarPosition(isGaussian, depth+1)
			}
		}
	}
	return xPos, yPos
}

func starPositionInConflict(xPos, yPos float32) bool {
	mindistance := 00.3

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
	x := float32(windowConfig.Starfield_Width) * b
	y := float32(windowConfig.Starfield_Height) * c
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
