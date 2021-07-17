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

var spriteShader *render.Shader

// provide getter for other components with similar rendering needs
func SpriteShader() *render.Shader { return spriteShader }

// call this before loading any spritesheets
func InitSpriteloader(_window *render.Window) {
	Window = _window
	spriteLoaderIsInitialized = true
	spriteShader = render.NewShader(
		"./assets/shader/sprite.vert", "./assets/shader/sprite.frag")
}

type Spritesheet struct {
	tex            uint32
	xScale, yScale float32
}
type SpritesheetID uint32

type Sprite struct {
	spriteSheetId SpritesheetID
	x, y          int
}

type SpriteID uint32

// fetch internal sprite data for using custom OpenGL rendering
type SpriteMetadata struct {
	GpuTex         uint32
	PosX, PosY     int
	ScaleX, ScaleY float32
}

func GetSpriteMetadata(spriteID SpriteID) SpriteMetadata {
	sprite := sprites[spriteID]
	spritesheet := spritesheets[sprite.spriteSheetId]
	return SpriteMetadata{
		GpuTex: spritesheet.tex,
		PosX:   sprite.x,
		PosY:   sprite.y,
		ScaleX: spritesheet.xScale,
		ScaleY: spritesheet.yScale,
	}
}

var spritesheets = []Spritesheet{}
var sprites = []Sprite{}
var spriteIdsByName = make(map[string]SpriteID)

func AddSpriteSheet(path string, il *ImgLoader) SpritesheetID {
	img := il.GetImg(path)
	if img == nil {
		log.Fatalf("Cannot find image at path %v", path)
	}

	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(32) / float32(img.Bounds().Dx()),
		yScale: float32(32) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return SpritesheetID(len(spritesheets) - 1)
}

func LoadSpriteSheet(fname string) SpritesheetID {
	_, img, _ := LoadPng(fname)

	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(32) / float32(img.Bounds().Dx()),
		yScale: float32(32) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return SpritesheetID(len(spritesheets) - 1)
}

//Load spritesheet with rows and columns specified
func LoadSpriteSheetByColRow(fname string, row int, col int) SpritesheetID {
	_, img, _ := LoadPng(fname)

	if DEBUG {
		fmt.Println("img.Bounds(): ", img.Bounds())
		fmt.Println("xScale: ", float32(img.Bounds().Dx()/col)/float32(img.Bounds().Dx()))
		fmt.Println("yScale: ", float32(img.Bounds().Dy()/row)/float32(img.Bounds().Dy()))
	}
	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(img.Bounds().Dx()/col) / float32(img.Bounds().Dx()),
		yScale: float32(img.Bounds().Dy()/row) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})

	return SpritesheetID(len(spritesheets) - 1)
}

func LoadSpriteSheetByFrames(fname string, frames []Frames) SpritesheetID {
	_, img, _ := LoadPng(fname)

	if DEBUG {
		fmt.Println("img.Bounds().Dx: ", img.Bounds().Dx())
		fmt.Println("img.Bounds().Dy: ", img.Bounds().Dy())
		fmt.Println("xScale: ", float32(frames[0].Frame.W)/float32(img.Bounds().Dx()))
		fmt.Println("yScale: ", float32(frames[0].Frame.H)/float32(img.Bounds().Dy()))
	}
	spritesheets = append(spritesheets, Spritesheet{
		xScale: float32(frames[0].Frame.W) / float32(img.Bounds().Dx()),
		yScale: float32(frames[0].Frame.H) / float32(img.Bounds().Dy()),
		tex:    MakeTexture(img),
	})
	return SpritesheetID(len(spritesheets) - 1)
}

func LoadSingleSprite(fname string, name string) SpriteID {
	spritesheetId := LoadSpriteSheetByColRow(fname, 1, 1)
	LoadSprite(spritesheetId, name, 0, 0)
	return GetSpriteIdByName(name)
}

//Load sprite into internal sheet
func LoadSprite(spriteSheetId SpritesheetID, name string, x, y int) SpriteID {
	sprites = append(sprites, Sprite{spriteSheetId, x, y})
	spriteId := SpriteID(len(sprites) - 1)
	spriteIdsByName[name] = spriteId
	return spriteId
}

// convenient for loading multi-tiles,
// loads a rectangle of sprites from a spritesheet
func LoadSprites(
	spritesheetId SpritesheetID, name string,
	left, top, right, bottom int,
) []SpriteID {
	spriteIds := make([]SpriteID, (right-left+1)*(bottom-top+1))
	spriteIdIdx := 0
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			localX := x - left
			localY := y - bottom
			name := fmt.Sprintf("%s_%d_%d", name, localX, localY)
			spriteIds[spriteIdIdx] = LoadSprite(spritesheetId, name, x, y)
			spriteIdIdx++
		}
	}
	return spriteIds
}

//Get the id of loaded sprite by its registered name
func GetSpriteIdByName(name string) SpriteID {
	spriteId, ok := spriteIdsByName[name]
	if !ok {
		log.Fatalf("sprite with name [%v] does not exist", name)
	}
	return spriteId
}

var SpriteRenderDistance float32 = 10

//Draw sprite specified with spriteId at x,y position
//this function is for testing, will not be used later on
func DrawSpriteQuadOptions(
	xpos, ypos, xwidth, yheight float32, spriteId SpriteID,
	opts DrawOptions,
) {
	worldTransform := mgl32.Mat4.Mul4(
		mgl32.Translate3D(float32(xpos), float32(ypos), -SpriteRenderDistance),
		mgl32.Scale3D(float32(xwidth), float32(yheight), 1),
	)
	DrawSpriteQuadMatrix(worldTransform, spriteId, opts)
}

func DrawSpriteQuad(
	xpos, ypos, xwidth, yheight float32, spriteId SpriteID,
) {
	DrawSpriteQuadOptions(xpos, ypos, xwidth, yheight, spriteId, NewDrawOptions())
}

func DrawSpriteQuadMatrix(
	worldTransform mgl32.Mat4, spriteId SpriteID, opts DrawOptions,
) {
	DrawSpriteQuadContext(render.Context{
		World:      worldTransform,
		Projection: Window.GetProjectionMatrix(),
	}, spriteId, opts)
}

func DrawSpriteQuadContext(
	ctx render.Context, spriteId SpriteID, opts DrawOptions,
) {
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

	spriteShader.Use()
	spriteShader.SetUint("outTexture", 0)
	spriteShader.SetVec2F("texScale", spritesheet.xScale, spritesheet.yScale)
	spriteShader.SetVec2F("texOffset", float32(sprite.x), float32(sprite.y))
	spriteShader.SetMat4("world", &ctx.World)
	spriteShader.SetMat4("projection", &ctx.Projection)

	color := opts.Color()
	spriteShader.SetVec4("color", &color)

	gl.BindVertexArray(render.QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	// restore texScale and texOffset to defaults
	// TODO separate GPU programs such that this becomes unecessary
	spriteShader.SetVec2F("texScale", 1, 1)
	spriteShader.SetVec2F("texOffset", 0, 0)
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
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)
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
	//unbind
	gl.BindVertexArray(0)

	return vao
}

//temporary overloaded function
func DrawSpriteQuadCustom(
	xpos, ypos, xwidth, yheight float32,
	spriteId SpriteID, program uint32,
) {
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
		mgl32.Translate3D(float32(xpos), float32(ypos), 0),
		mgl32.Scale3D(float32(xwidth), float32(yheight), 1),
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("world\x00")),
		1, false, &worldTransform[0],
	)

	w := float32(Window.Width)
	h := float32(Window.Height)
	projectTransform := mgl32.Ortho(
		0, w,
		0, h,
		-1, 1,
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("projection\x00")),
		1, false, &projectTransform[0],
	)

	gl.BindVertexArray(render.QuadVao)
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

type DrawOptions struct {
	Alpha float32
}

func NewDrawOptions() DrawOptions {
	return DrawOptions{Alpha: 1}
}

func (opts DrawOptions) Color() mgl32.Vec4 {
	return mgl32.Vec4{1, 1, 1, opts.Alpha}
}
