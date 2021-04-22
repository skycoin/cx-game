package spriteloader

import (
	"log"
	"os"
	"image/png"
	"image"
	//"image/draw"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
)

var window render.Window
// call this before loading any spritesheets
func initSpriteloader(_window render.Window) {
	window = _window
	quadVao = makeQuadVao()
}

type Spritesheet struct {
	spriteWidth int
	spriteHeight int
	tex uint32
}

type Sprite struct {
	spriteSheetId int
	x int
	y int
}

var spritesheets = []Spritesheet{};
var sprites = []Sprite{};
var spriteIdsByName = make(map[string]int);

func LoadSpriteSheet(fname string) int {
	log.Print("loading sprite sheet from "+fname)

	img := LoadPng(fname)

	spritesheets = append(spritesheets, Spritesheet{
		spriteWidth: img.Bounds().Dx() / 32,
		spriteHeight: img.Bounds().Dx() / 32,
		tex: makeTexture(img),
	})

	return len(spritesheets)-1
}

func LoadSprite(spriteSheetId int, name string, x,y int) {
	sprites = append(sprites,Sprite{spriteSheetId,x,y})
	spriteIdsByName[name] = len(sprites)-1
}

func GetSpriteIdByName(name string) int {
	spriteId, ok := spriteIdsByName[name]
	if (!ok) {
		log.Fatalf("sprite with name [%v] does not exist",name)
	}
	return spriteId
}

func DrawSpriteQuad(xpos, ypos, xwidth, yheight, spriteId int) {
	// this method assumes:
	// - perspective is already set correctly
	gl.Uniform1i(
		gl.GetUniformLocation(window.Program, gl.Str("ourTexture\x00")), 0,
	)
	worldTranslate := mgl32.Translate3D(float32(xpos), float32(ypos), 0)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(window.Program,gl.Str("world\x00")),
		1,false,&worldTranslate[0],
	)
	gl.BindVertexArray(quadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func LoadPng(fname string) *image.RGBA {
	imgFile, err := os.Open(fname)

	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	im, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	return image.NewRGBA(im.Bounds())
}

func makeTexture(img *image.RGBA) uint32 {
	var tex uint32
	gl.GenTextures(1,&tex);

	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(img.Rect.Dx()), int32(img.Rect.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix),
	)

	return tex;
}

// x,y,z,u,v
var quadVertexAttributes = []float32{
	1, 1, 0, 1, 0,
	1, -1, 0, 1, 1,
	-1, 1, 0, 0, 0,

	1, -1, 0, 1, 1,
	-1, -1, 0, 0, 1,
	-1, 1, 0, 0, 0,
}

var quadVao uint32
func makeQuadVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1,&vbo);

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


