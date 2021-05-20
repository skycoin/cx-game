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

const (
	width  = 800
	height = 600
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
	probability = 0.3
	stars       []float32
	starAmount  int     = 500
	maxStarSize float32 = 3
	// fileHasChanged chan struct{}
)

func init() {
	runtime.LockOSThread()
}
func main() {
	// fileHasChanged = make(chan struct{})
	// go utility.CheckAndReload("./test/perlin_map/config/perlin.yaml", noise, nil)

	win := render.NewWindow(width, height, false)
	window := win.Window

	shader := utility.NewShader(
		"./test/perlin_map/shaders/perlin_vertex.glsl",
		"./test/perlin_map/shaders/perlin_fragment.glsl",
	)
	shader.Use()

	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)

	// tex := utility.Make1DTexture(noise.GradFile)
	//pick from 1-11 gradient
	tex := getGradient(1)
	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_1D, tex)
	shader.SetInt("texture_1d", 0)

	projection := mgl32.Ortho2D(0, width, 0, height)
	shader.SetMat4("projection", &projection)

	var vao1, vao2 uint32
	vao1 = genStarField()
	vao2 = vao1
	go func() {
		for {
			vao2 = genStarField()
			time.Sleep(1500 * time.Millisecond)
			vao1 = vao2
		}
	}()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindVertexArray(vao1)

		gl.DrawArrays(gl.POINTS, 0, int32(len(stars)))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func genStarField() uint32 {
	stars = []float32{}
	perlinMap := genPerlin()

	for i := 0; i < starAmount; i++ {
		var xPos, yPos float32
		xPos, yPos = genStar(perlinMap)

		stars = append(stars, xPos, yPos, 0.0, rand.Float32()*maxStarSize, rand.Float32())
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

func genStar(perlinMap [][]float32) (float32, float32) {
	xPos := rand.Intn(width)
	yPos := rand.Intn(height)

	perlinProb := perlinMap[yPos][xPos]

	deleted := rand.Float32() * perlinProb
	if deleted > float32(probability) {
		return float32(xPos), float32(yPos)
	} else {
		return genStar(perlinMap)
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
