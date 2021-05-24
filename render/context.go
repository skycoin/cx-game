package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Context struct {
	World mgl32.Mat4
	Size mgl32.Vec2
	Projection mgl32.Mat4
}

func (ctx Context) PushLocal(local mgl32.Mat4) Context {
	newCtx := ctx
	newCtx.World = newCtx.World.Mul4(local)
	return newCtx
}

func (ctx Context) PushView(view mgl32.Mat4) Context {
	newCtx := ctx
	newCtx.Projection = newCtx.Projection.Mul4(view)
	return newCtx
}

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
	return Context {
		World: mgl32.Ident4(),
		Size: mgl32.Vec2 { w/PixelsPerTile,h/PixelsPerTile },
		Projection: projectTransform,
	}
}
