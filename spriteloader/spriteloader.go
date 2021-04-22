package spriteloader

import (
	"os"
	"image/png"
	"image"
	"github.com/go-gl/gl/v4.1-core/gl"
	/*
	"image/draw"
	"image/png"
	*/
	"log"
)

// call this before loading any spritesheets
func initSpriteloader() {
	quadVao = makeQuadVao();
}

type Spritesheet struct {
	img image.Image
	spriteWidth int
	spriteHeight int
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
		img: img,
		spriteWidth: img.Bounds().Dx() / 32,
		spriteHeight: img.Bounds().Dx() / 32,
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
	// TODO
}

func LoadPng(fname string) image.Image {
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


