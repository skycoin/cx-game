package render

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/cxmath/math32i"
)

type SpriteID int

// Usage:
// Prepare(...)
// DrawSprite(...) - many times
// Flush()

var worldWidth float32
var cameraTransform mgl32.Mat4
var minFilter, magFilter int32 = gl.NEAREST, gl.NEAREST

func SetWorldWidth(w float32) { worldWidth = w }
func SetCameraTransform(mat mgl32.Mat4) {
	cameraTransform = mat
}

type SpriteDrawOptions struct {
	Outline bool
}

func NewSpriteDrawOptions() SpriteDrawOptions {
	return SpriteDrawOptions{}
}

func (opts SpriteDrawOptions) Framebuffer() Framebuffer {
	if opts.Outline {
		return FRAMEBUFFER_PLANET
	} else {
		return FRAMEBUFFER_MAIN
	}
}

type SpriteDraw struct {
	Sprite      Sprite
	Model       mgl32.Mat4
	View        mgl32.Mat4
	UVTransform mgl32.Mat3
	LightColor  mgl32.Vec3
	Options     SpriteDrawOptions
}

var spriteDrawsPerAtlasPerFramebuffer = map[Framebuffer]map[Texture][]SpriteDraw{}
var spriteDrawsPerAtlasUI = map[Texture][]SpriteDraw{}

func drawSprite(model, view mgl32.Mat4, id SpriteID, opts SpriteDrawOptions) {
	sprite := sprites[id]
	atlas := sprite.Texture
	framebuffer := opts.Framebuffer()
	spriteDrawsPerAtlas, ok := spriteDrawsPerAtlasPerFramebuffer[framebuffer]
	if !ok {
		spriteDrawsPerAtlasPerFramebuffer[framebuffer] =
			map[Texture][]SpriteDraw{}
	}
	spriteDrawsPerAtlas = spriteDrawsPerAtlasPerFramebuffer[framebuffer]
	spriteDrawsPerAtlas[atlas] =
		append(spriteDrawsPerAtlasPerFramebuffer[framebuffer][atlas],
			SpriteDraw{
				Sprite:      sprite,
				Model:       model,
				View:        view,
				UVTransform: sprite.Transform,
			})
}

func drawUISprite(model, view mgl32.Mat4, id SpriteID, opts SpriteDrawOptions) {
	sprite := sprites[id]
	atlas := sprite.Texture

	spriteDrawsPerAtlasUI[atlas] =
		append(spriteDrawsPerAtlasUI[atlas],
			SpriteDraw{
				Sprite:      sprite,
				Model:       model,
				View:        view,
				UVTransform: sprite.Transform,
			})
}

// unaffected by camera movement
func DrawUISprite(transform mgl32.Mat4, id SpriteID, opts SpriteDrawOptions) {
	view := mgl32.Ident4()
	drawUISprite(transform, view, id, opts)
}

func wrapTransform(raw mgl32.Mat4) mgl32.Mat4 {
	rawX := raw.At(0, 3)
	x := math32.PositiveModulo(rawX, worldWidth)
	camX := cameraTransform.At(0, 3)
	if x-camX > worldWidth/2 {
		x -= worldWidth
	}
	if x-camX < -worldWidth/2 {
		x += worldWidth
	}

	translate := mgl32.Translate3D(x-rawX, 0, 0)
	return translate.Mul4(raw)
}

// affected by camera movement
// TODO implement wrap-around in here
func DrawWorldSprite(transform mgl32.Mat4, id SpriteID, opts SpriteDrawOptions) {
	wrappedTransform := wrapTransform(transform)
	model := wrappedTransform
	view := cameraTransform.Inv()
	drawSprite(model, view, id, opts)
}

//draw without flushing

func Flush(zoom float32) {
	flushSpriteDraws(zoom)
	// flushColorDraws(Projection)
	if drawBBoxLines {
		flushBBoxLineDraws()
	}

}

func FlushUI() {
	flushSpriteUIDraws()
	flushColorDraws(Projection)
}

func drawFramebufferSprites(framebuffer Framebuffer) {
	framebuffer.Bind(gl.FRAMEBUFFER)
	spriteDrawsPerAtlas := spriteDrawsPerAtlasPerFramebuffer[framebuffer]
	for atlas, spriteDraws := range spriteDrawsPerAtlas {
		drawAtlasSprites(atlas, spriteDraws)
	}
}

func flushSpriteUIDraws() {
	spriteProgram.Use()
	defer spriteProgram.StopUsing()

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.DEPTH_TEST)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(QuadVao)

	spriteProgram.SetMat4("projection", &Projection)

	for atlas, spriteDraws := range spriteDrawsPerAtlasUI {
		drawAtlasSprites(atlas, spriteDraws)
	}
	spriteDrawsPerAtlasUI = make(map[Texture][]SpriteDraw)

}

func flushSpriteDraws(zoom float32) {
	spriteProgram.Use()
	defer spriteProgram.StopUsing()

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.DEPTH_TEST)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(QuadVao)

	spriteProgram.SetMat4("projection", &Projection)

	physicalViewport := currentViewport
	virtualViewport :=
		Viewport{
			0, 0,
			constants.VIRTUAL_VIEWPORT_WIDTH2,
			constants.VIRTUAL_VIEWPORT_HEIGHT2,
		}
	virtualViewport.Use()
	FRAMEBUFFER_PLANET.Bind(gl.FRAMEBUFFER)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	drawFramebufferSprites(FRAMEBUFFER_PLANET)

	physicalViewport.Use()

	drawFramebufferSprites(FRAMEBUFFER_MAIN)

	spriteDrawsPerAtlasPerFramebuffer = // clear sprite draws
		make(map[Framebuffer]map[Texture][]SpriteDraw)

	FRAMEBUFFER_MAIN.Bind(gl.FRAMEBUFFER)

	outlineProgram.Use()
	defer outlineProgram.StopUsing()

	texelSize := mgl32.Vec2{
		zoom * 1.0 / float32(constants.VIRTUAL_VIEWPORT_WIDTH2),
		zoom * 1.0 / float32(constants.VIRTUAL_VIEWPORT_HEIGHT2),
	}
	outlineProgram.SetVec2("texelSize", &texelSize)
	outlineProgram.SetVec4("borderColor", &constants.OUTLINE_BORDER_COLOR)

	gl.BindVertexArray(Quad2Vao)

	gl.BindTexture(gl.TEXTURE_2D, RENDERTEXTURE_PLANET)
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // draw quad

	gl.BindVertexArray(QuadVao)

}

func drawAtlasSprites(atlas Texture, spriteDraws []SpriteDraw) {
	atlas.Bind()
	atlas.SetTextureFiltering(minFilter, magFilter)
	defer atlas.Unbind()

	uniforms := extractUniforms(spriteDraws)
	for _, batch := range uniforms.Batch(constants.DRAW_SPRITE_BATCH_SIZE) {
		drawInstancedQuads(batch)
	}
}

func extractUniforms(spriteDraws []SpriteDraw) Uniforms {
	uniforms := NewUniforms(int32(len(spriteDraws)))
	for idx, spriteDraw := range spriteDraws {
		uniforms.Models[idx] = spriteDraw.Model
		uniforms.Views[idx] = spriteDraw.View
		uniforms.UVTransforms[idx] = spriteDraw.UVTransform
	}
	return uniforms
}

type Uniforms struct {
	Count        int32
	Models       []mgl32.Mat4
	Views        []mgl32.Mat4
	UVTransforms []mgl32.Mat3
}

func NewUniforms(count int32) Uniforms {
	return Uniforms{
		Count:        count,
		Models:       make([]mgl32.Mat4, count),
		Views:        make([]mgl32.Mat4, count),
		UVTransforms: make([]mgl32.Mat3, count),
	}
}

func (u Uniforms) Batch(batchSize int32) []Uniforms {
	numBatches := divideRoundUp(u.Count, batchSize)
	batches := make([]Uniforms, numBatches)

	for i := int32(0); i < numBatches; i++ {
		start := batchSize * i
		stop := math32i.Min(batchSize*(i+1), u.Count)
		batches[i] = u.Range(start, stop)
	}

	return batches
}

func (u Uniforms) Range(start, stop int32) Uniforms {
	return Uniforms{
		Count:        stop - start,
		Models:       u.Models[start:stop],
		Views:        u.Views[start:stop],
		UVTransforms: u.UVTransforms[start:stop],
	}
}

func divideRoundUp(a, b int32) int32 {
	if a%b == 0 {
		return a / b
	} else {
		return a/b + 1
	}
}

func drawInstancedQuads(batch Uniforms) {
	spriteProgram.SetMat4s("models", batch.Models)
	spriteProgram.SetMat4s("views", batch.Views)
	spriteProgram.SetMat3s("uvtransforms", batch.UVTransforms)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, batch.Count)
}

//toggles and sets texture filtering options
//minifying includes nearest and linear interpolation, as well as mipmap
//magnifying includes nearest and linear interpolation

func ToggleFiltering() {
	var message string
	if magFilter == gl.LINEAR {
		magFilter = gl.NEAREST
		message += "MAG: NEAREST  "
	} else {
		magFilter = gl.LINEAR
		message += "MAG: LINEAR  "
	}

	switch minFilter {
	case gl.LINEAR:
		minFilter = gl.NEAREST
		message += "MIN: NEAREST"
	case gl.NEAREST:
		minFilter = gl.NEAREST_MIPMAP_NEAREST
		message += "MIN: NEAREST_MIPMAP_NEAREST"
	case gl.NEAREST_MIPMAP_NEAREST:
		minFilter = gl.NEAREST_MIPMAP_LINEAR
		message += "MIN: NEAREST_MIPMAP_LINEAR"
	case gl.NEAREST_MIPMAP_LINEAR:
		minFilter = gl.LINEAR_MIPMAP_NEAREST
		message += "MIN: NEAREST_M_LINEAR"
	case gl.LINEAR_MIPMAP_NEAREST:
		minFilter = gl.LINEAR_MIPMAP_LINEAR
		message += "MIN: LINEAR_M_LINEAR"
	case gl.LINEAR_MIPMAP_LINEAR:
		minFilter = gl.LINEAR
		message += "MIN: LINEAR"
	default:
		//should not branch this way
		log.Panic("something is wrong")
	}

	fmt.Println(message)
}

func CyclePixelSnap() {
	switch spriteProgram {
	case &spriteProgram1:
		spriteProgram = &spriteProgram2
		fmt.Println("PIXEL SNAPPING: snapping to a pixel")
	case &spriteProgram2:
		spriteProgram = &spriteProgram3
		fmt.Println("PIXEL SNAPPING: snapping to half a pixel")
	case &spriteProgram3:
		spriteProgram = &spriteProgram1
		fmt.Println("PIXEL SNAPPING: no snapping")
	default:
		log.Fatalf("Pixel snapping error")
	}
}
