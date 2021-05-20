package spriteloader

import (
	"fmt"
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
)

var spriteLoaderIsInitialized = false
var Window *render.Window

// call this before loading any spritesheets
func InitSpriteloader(_window *render.Window) {
	Window = _window
	QuadVao = MakeQuadVao()
	spriteLoaderIsInitialized = true
}

type Spritesheet struct {
	tex            uint32
	xScale, yScale float32
}

type Sprite struct {
	spriteSheetId int
	x, y          int
}

var spritesheets = []Spritesheet{}
var sprites = []Sprite{}
var spriteIdsByName = make(map[string]int)

func AddSpriteSheet(path string, il *ImgLoader) int {
	img := il.GetImg(path)
	if img == nil {
		return -1
	}

	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(32) / float32(img.Bounds().Dx()),
		yScale: float32(32) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return len(spritesheets) - 1
}

func LoadSpriteSheet(fname string) int {
	_, img, _ := LoadPng(fname)

	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(32) / float32(img.Bounds().Dx()),
		yScale: float32(32) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return len(spritesheets) - 1
}

//Load spritesheet with rows and columns specified
func LoadSpriteSheetByColRow(fname string, row int, col int) int {
	_, img, _ := LoadPng(fname)

	fmt.Println("xScale: ", float32(img.Bounds().Dx()/col)/float32(img.Bounds().Dx()))
	fmt.Println("yScale: ", float32(img.Bounds().Dy()/row)/float32(img.Bounds().Dy()))
	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(img.Bounds().Dx()/col) / float32(img.Bounds().Dx()),
		yScale: float32(img.Bounds().Dy()/row) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return len(spritesheets) - 1
}

func LoadSingleSprite(fname string, name string) int {
	spritesheetId := LoadSpriteSheetByColRow(fname, 1, 1)
	LoadSprite(spritesheetId, name, 0, 0)
	return GetSpriteIdByName(name)
}

//Load sprite into internal sheet
func LoadSprite(spriteSheetId int, name string, x, y int) {
	sprites = append(sprites, Sprite{spriteSheetId, x, y})
	spriteIdsByName[name] = len(sprites) - 1
}

//Get the id of loaded sprite by its registered name
func GetSpriteIdByName(name string) int {
	spriteId, ok := spriteIdsByName[name]
	if !ok {
		log.Fatalf("sprite with name [%v] does not exist", name)
	}
	return spriteId
}

var SpriteRenderDistance float32 = 10

//Draw sprite specified with spriteId at x,y position
func DrawSpriteQuad(xpos, ypos, xwidth, yheight float32, spriteId int) {
	worldTransform := mgl32.Mat4.Mul4(
		mgl32.Translate3D(float32(xpos), float32(ypos), -SpriteRenderDistance),
		mgl32.Scale3D(float32(xwidth), float32(yheight), 1),
	)
	DrawSpriteQuadMatrix(worldTransform, spriteId)
}

func DrawSpriteQuadMatrix(worldTransform mgl32.Mat4, spriteId int) {
	// TODO this method probably shouldn't be responsible
	// for setting up the projection matrix.
	// clarify responsibilities later
	sprite := sprites[spriteId]
	spritesheet := spritesheets[sprite.spriteSheetId]

	// bind texture
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// NOTE depth test is disabled for now,
	// as we assume that objects are drawn in the correct order
	gl.Disable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, spritesheet.tex)

	gl.UseProgram(Window.Program)
	gl.Uniform1ui(
		gl.GetUniformLocation(Window.Program, gl.Str("ourTexture\x00")),
		// spritesheet.tex,
		0,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(Window.Program, gl.Str("texScale\x00")),
		spritesheet.xScale, spritesheet.yScale,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(Window.Program, gl.Str("texOffset\x00")),
		float32(sprite.x), float32(sprite.y),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(Window.Program, gl.Str("world\x00")),
		1, false, &worldTransform[0],
	)

	aspect := float32(Window.Width) / float32(Window.Height)
	projectTransform := mgl32.Perspective(
		mgl32.DegToRad(45), aspect, 0.1, 100.0,
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(Window.Program, gl.Str("projection\x00")),
		1, false, &projectTransform[0],
	)

	gl.BindVertexArray(QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	// restore texScale and texOffset to defaults
	// TODO separate GPU programs such that this becomes unecessary
	gl.Uniform2f(
		gl.GetUniformLocation(Window.Program, gl.Str("texScale\x00")),
		1, 1,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(Window.Program, gl.Str("texOffset\x00")),
		0, 0,
	)
}

// upload an in-memory RGBA image to the GPU
func MakeTexture(img *image.RGBA) uint32 {
	if !spriteLoaderIsInitialized {
		log.Fatalln("Sprite loader is not initialized")
	}
	var tex uint32
	gl.GenTextures(1, &tex)

	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// sample nearest pixel such that we see nice pixel art
	// and not a blurry mess
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(img.Rect.Dx()), int32(img.Rect.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix),
	)

	return tex
}

// x,y,z,u,v
var quadVertexAttributes = []float32{
	0.5, 0.5, 0, 1, 0,
	0.5, -0.5, 0, 1, 1,
	-0.5, 0.5, 0, 0, 0,

	0.5, -0.5, 0, 1, 1,
	-0.5, -0.5, 0, 0, 1,
	-0.5, 0.5, 0, 0, 0,
}

var QuadVao uint32

func MakeQuadVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(quadVertexAttributes),
		gl.Ptr(quadVertexAttributes),
		gl.STATIC_DRAW,
	)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}

//temporary overloaded function
func DrawSpriteQuadCustom(xpos, ypos, xwidth, yheight float32, spriteId int, program uint32) {
	// TODO this method probably shouldn't be responsible
	// for setting up the projection matrix.
	// clarify responsibilities later
	sprite := sprites[spriteId]
	spritesheet := spritesheets[sprite.spriteSheetId]

	// bind texture
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// NOTE depth test is disabled for now,
	// as we assume that objects are drawn in the correct order
	gl.Disable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, spritesheet.tex)

	gl.UseProgram(program)
	gl.Uniform1ui(
		gl.GetUniformLocation(program, gl.Str("ourTexture\x00")),
		spritesheet.tex,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texScale\x00")),
		spritesheet.xScale, spritesheet.yScale,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texOffset\x00")),
		float32(sprite.x), float32(sprite.y),
	)

	worldTransform := mgl32.Mat4.Mul4(
		mgl32.Translate3D(float32(xpos), float32(ypos), -10),
		mgl32.Scale3D(float32(xwidth), float32(yheight), 1),
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("world\x00")),
		1, false, &worldTransform[0],
	)

	aspect := float32(Window.Width) / float32(Window.Height)
	projectTransform := mgl32.Perspective(
		mgl32.DegToRad(45), aspect, 0.1, 100.0,
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("projection\x00")),
		1, false, &projectTransform[0],
	)

	gl.BindVertexArray(QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	// restore texScale and texOffset to defaults
	// TODO separate GPU programs such that this becomes unecessary
	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texScale\x00")),
		1, 1,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texOffset\x00")),
		0, 0,
	)
}
