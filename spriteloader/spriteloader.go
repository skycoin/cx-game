package spriteloader

import (
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Vec2i struct {
	X int32
	Y int32
}

type Vec2ui struct {
	X uint32
	Y uint32
}

type Texture struct {
	id   uint32
	size Vec2ui
}

type SpriteSheet struct {
	id      uint32
	texture Texture
	bounds  Vec2ui
}

type Sprite struct {
	id     uint32
	name   string
	ss     *SpriteSheet
	bounds image.Rectangle
}

const spriteWidth = 32
const spriteHeight = 32

var spriteNameMap map[string]Sprite = make(map[string]Sprite)
var sprites []Sprite = make([]Sprite, 0)
var spriteSheets []SpriteSheet = make([]SpriteSheet, 0)

// Helper functions

func loadTexture(file string) Texture {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	size := Vec2ui{
		X: uint32(rgba.Bounds().Size().X),
		Y: uint32(rgba.Bounds().Size().Y),
	}

	return Texture{id: texture, size: size}
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))

	return vao
}

func glDraw(vao uint32, texId uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

// SpriteSheet

// LoadSpriteSheet from a png file given it's path
func LoadSpriteSheet(file string) uint32 {
	texture := loadTexture(file)

	bounds := Vec2ui{
		X: texture.size.X / spriteWidth,
		Y: texture.size.Y / spriteHeight,
	}

	ss := SpriteSheet{
		id:      uint32(len(spriteSheets)),
		texture: texture,
		bounds:  bounds,
	}

	spriteSheets = append(spriteSheets, ss)

	return ss.id
}

// LoadSprite from a sprite sheet with it's id
func LoadSprite(ssId uint32, name string, x uint32, y uint32) uint32 {
	ss := spriteSheets[ssId]

	if _, ok := spriteNameMap[name]; ok {
		log.Fatal("name for sprite already exist")
	}

	if x >= ss.bounds.X || y >= ss.bounds.Y {
		log.Fatalln("sprite out of bounds")
	}

	x0 := int(x * spriteWidth)
	y0 := int(y * spriteHeight)
	x1 := int(x0 + spriteWidth)
	y1 := int(y0 + spriteHeight)

	rect := image.Rect(x0, y0, x1, y1)

	sprite := Sprite{
		id:     uint32(len(sprites)),
		name:   name,
		ss:     &ss,
		bounds: rect,
	}

	sprites = append(sprites, sprite)
	spriteNameMap[name] = sprite

	return sprite.id
}

// Sprite

// GetSpriteByName
func GetSpriteByName(name string) uint32 {
	return spriteNameMap[name].id
}

// DrawQuad draws a sprite providing it's position, size and it's id
func DrawQuad(xpos float32, ypos float32, width float32, height float32, spriteId uint32) {
	sprite := sprites[spriteId]

	hw := width / 2.0
	hh := height / 2.0

	x0 := xpos - hw
	y0 := ypos - hh
	x1 := xpos + hw
	y1 := ypos + hh

	tw := 1.0 / float32(sprite.ss.texture.size.X)
	th := 1.0 / float32(sprite.ss.texture.size.Y)

	tx0 := tw * float32(sprite.bounds.Min.X)
	tx1 := tw * float32(sprite.bounds.Max.X)
	ty0 := th * float32(sprite.bounds.Min.Y)
	ty1 := th * float32(sprite.bounds.Max.Y)

	vertices := []float32{
		x1, y1, 0, tx1, ty0,
		x1, y0, 0, tx1, ty1,
		x0, y1, 0, tx0, ty0,

		x1, y0, 0, tx1, ty1,
		x0, y0, 0, tx0, ty1,
		x0, y1, 0, tx0, tx0,
	}

	vao := makeVao(vertices)

	glDraw(vao, sprite.ss.texture.id)
}
