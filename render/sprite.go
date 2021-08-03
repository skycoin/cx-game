package render

import (
	"log"
	"strconv"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var spriteProgram Program
const numInstances int = 100

func initSprite() {
	config := NewShaderConfig(
		"./assets/shader/spritesheet.vert", "./assets/shader/spritesheet.frag",
	)
	config.Define("NUM_INSTANCES", strconv.Itoa(numInstances))

	spriteProgram = config.Compile()
}

func Init() { initSprite() }

type SpriteSheet struct {
	Texture Texture
	Sprites []Sprite
}

func (sheet SpriteSheet) Sprite(name string) (Sprite,bool) {
	for _,sprite := range sheet.Sprites {
		if sprite.Name == name { return sprite,true }
	}
	return Sprite{},false
}

type Sprite struct {
	Name string
	Transform mgl32.Mat3
	Model mgl32.Mat4
}

type SpriteRenderParams struct {
	Sprite Sprite
	MVP mgl32.Mat4
}

type SpritesRenderParams []SpriteRenderParams

func (params SpritesRenderParams) SpriteTransforms() []mgl32.Mat3 {
	transforms := make([]mgl32.Mat3, len(params))
	for idx := range transforms {
		transforms[idx] = params[idx].Sprite.Transform
	}
	return transforms
}

func (params SpritesRenderParams) SpriteModels() []mgl32.Mat4 {
	models := make([]mgl32.Mat4, len(params))
	for idx := range models {
		models[idx] = params[idx].Sprite.Model
	}
	return models
}

func (params SpritesRenderParams) MVPs() []mgl32.Mat4 {
	MVPs := make([]mgl32.Mat4, len(params))
	for idx := range MVPs {
		MVPs[idx] = params[idx].MVP
	}
	return MVPs
}

func (sheet *SpriteSheet) Draw(params SpritesRenderParams) {
	spriteProgram.Use()
	defer spriteProgram.StopUsing()

	sheet.Texture.Bind()
	defer sheet.Texture.Unbind()

	gl.BindVertexArray(QuadVao)

	spriteProgram.SetMat3s("texTransforms",params.SpriteTransforms())
	spriteProgram.SetMat4s("SpriteModels",params.SpriteModels())
	spriteProgram.SetMat4s("MVPs",params.MVPs())

	instances := int32(len(params))
	if int(instances) > numInstances { log.Fatal("exceeded instance limit") }
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, vertsPerQuad, instances)
}
