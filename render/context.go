package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

// the render context contains information
// used to transform graphics to the screen
type Context struct {
	World      mgl32.Mat4 // 4x4 matrix for local => world
	Size       mgl32.Vec2 // window size for ratio / frustum calculations
	Projection mgl32.Mat4 // 4x4 projection matrix for world => screen
}

// push a local transformation onto the render context.
// can use to build hierarchies.
func (ctx Context) PushLocal(local mgl32.Mat4) Context {
	newCtx := ctx
	newCtx.World = newCtx.World.Mul4(local)
	return newCtx
}

// push a view matrix onto the render context.
// in the result, the projection matrix is really 
// a projection-view matrix
func (ctx Context) PushView(view mgl32.Mat4) Context {
	newCtx := ctx
	newCtx.Projection = newCtx.Projection.Mul4(view)
	return newCtx
}

// compute the model-view-projection matrix for this render context
func (ctx Context) MVP() mgl32.Mat4 {
	return ctx.Projection.Mul4(ctx.World)
}

const PixelsPerTile = 32

func (window Window) DefaultRenderContext() Context {
	w := float32(window.Width)
	h := float32(window.Height)
	projectTransform := mgl32.Ortho(
		-w/2/PixelsPerTile, w/2/PixelsPerTile,
		-h/2/PixelsPerTile, h/2/PixelsPerTile,
		-1, 1000,
	)
	return Context{
		World:      mgl32.Ident4(),
		Size:       mgl32.Vec2{w / PixelsPerTile, h / PixelsPerTile},
		Projection: projectTransform,
	}
}

// use for rendering things that need to wrap around world
type WorldContext interface {
	Projection() mgl32.Mat4
	ModelToModelView(model mgl32.Mat4) mgl32.Mat4
	ModelToModelViewProjection(model mgl32.Mat4) mgl32.Mat4
}
