package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
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
	StarAmount          int
	Gaussian_Percentage int

	gaussianAmount int
}

type RandType int

const (
	width  = 800
	height = 600

	NORMAL RandType = iota
	MEAN
)

var (
	noise = &noiseSettings{
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

	//breaks if more than 1
	perlinMap           [][]float32
	probability         = 0.3
	stars               []float32
	starAmount          int     = 1500
	gaussian_percentage int     = 25
	gaussian_amount             = starAmount * gaussian_percentage / 100
	maxStarSize         float32 = 3
	// fileHasChanged chan struct{}
)

func init() {
	runtime.LockOSThread()
}
func main() {
	// fileHasChanged = make(chan struct{})
	go utility.CheckAndReload("./test/gaussian/config/perlin.yaml", noise, nil)

	win := render.NewWindow(width, height, false)
	window := win.Window

	shader := utility.NewShader(
		"./test/gaussian/shaders/vertex.glsl",
		"./test/gaussian/shaders/fragment.glsl",
	)
	shader.Use()

	gl.Enable(gl.POINT)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)

	// tex := utility.Make1DTexture(noise.GradFile)
	//pick from 1-11 gradient
	tex := utility.Make1DTexture(noise.GradFile)
	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_1D, tex)
	shader.SetInt("texture_1d", 0)

	projection := mgl32.Ortho2D(0, width, 0, height)
	shader.SetMat4("projection", &projection)

	vao1 := genStarField(0.5)
	vao2 := genStarField(0.9)

	var current uint32 = vao1

	go func() {
		for {
			time.Sleep(1500 * time.Millisecond)
			if current == vao1 {
				current = vao2
			} else {
				current = vao1
			}
		}
	}()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindVertexArray(current)

		gl.DrawArrays(gl.POINTS, 0, int32(len(stars)))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func genStarField(n float32) uint32 {
	stars = []float32{}
	perlinMap = genPerlin()

	for i := 0; i < starAmount; i++ {
		var xPos, yPos, gradientValue float32
		gradientValue = 0.3
		if i < gaussian_amount {
			for {
				xPos, yPos = getRandPosition(MEAN)
				gradientValue = n
				z := gaussianTheta(xPos, yPos)
				if 1-z > 0.4*rand.Float32() {
					continue
				}
				xPos, yPos = xPos*width, yPos*height
				break
			}
		} else {
			xPos, yPos = genStar()
		}

		stars = append(stars, xPos, yPos, 0.0, rand.Float32()*maxStarSize, gradientValue)
	}

	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao)
	gl.BufferData(gl.ARRAY_BUFFER, len(stars)*4, gl.Ptr(stars), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 5*4, gl.PtrOffset(4*4))
	gl.EnableVertexAttribArray(2)

	return vao
}

func genPerlin() [][]float32 {
	grid := make([][]float32, 0)
	myPerlin := perlin.NewPerlin2D(
		noise.Seed,
		noise.X,
		noise.Xs,
		noise.Gradmax,
	)
	max := float32(math.Sqrt2 / (1.9 * noise.Contrast))
	min := float32(-math.Sqrt2 / (1.9 * noise.Contrast))
	for y := 0; y < height; y++ {
		grid = append(grid, []float32{})
		for x := 0; x < width; x++ {
			result := myPerlin.Noise(
				float32(x)*noise.Scale,
				float32(y)*noise.Scale,
				noise.Persistance,
				noise.Lacunarity,
				noise.Octaves,
			)
			result = clamp(result-min/(max-min), 0.0, 1.0)
			grid[y] = append(grid[y], result)
		}

	}

	return grid
}

func genStar() (float32, float32) {
	xPos := rand.Intn(width)
	yPos := rand.Intn(height)

	perlinProb := perlinMap[yPos][xPos]

	deleted := rand.Float32() * perlinProb
	if deleted > float32(probability) {
		return float32(xPos), float32(yPos)
	} else {
		return genStar()
	}
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

//gaussian
func gaussianTheta(x32, y32 float32) float32 {
	// x and y values range from 0 to 1
	x, y := float64(x32), float64(y32)
	var sigmaX, sigmaY, x0, y0 float64
	var A float64 = 1
	var theta float64
	sigmaX = 0.1
	sigmaY = 0.3
	x0 = 0.7
	y0 = 0.7
	theta = float64(DegToRad(45))

	a := math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaY, 2))
	b := -math.Sin(2*theta)/(4*math.Pow(sigmaX, 2)) + math.Sin(2*theta)/(4*math.Pow(sigmaY, 2))
	c := math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaY, 2))

	result := A * math.Exp(-(a*math.Pow(x-x0, 2) + 2*b*(x-x0)*(y-y0) + c*math.Pow(y-y0, 2)))

	return float32(result)
}

func DegToRad(angle float32) float32 {
	return math.Pi / 180 * angle
}

func getRandPosition(randType RandType) (float32, float32) {
	switch randType {
	case MEAN:
		var x, y float32
		n := 2
		for i := 0; i < n; i++ {
			x += rand.Float32()
			y += rand.Float32()
		}
		return x / float32(n), y / float32(n)
	default:
		return rand.Float32(), rand.Float32()
	}
}
