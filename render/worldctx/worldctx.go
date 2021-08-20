package worldctx

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

type worldContextImpl struct {
	camera *camera.Camera
	planet *world.Planet
}

func (wc worldContextImpl) Projection() mgl32.Mat4 {
	return wc.camera.GetProjectionMatrix()
}

func (wc worldContextImpl) ModelToModelView(model mgl32.Mat4) mgl32.Mat4 {
	// TODO wrap around
	// disp := wc.planet.ShortestDisplacement(
	// 	wc.camera.GetTransform().Col(3).Vec2(),
	// 	model.Col(3).Vec2(),
	// )

	// // zero out position since we are handling that manually
	// model.SetCol(3, mgl32.Vec4{0, 0, 0, 1})

	// return mgl32.Translate3D(disp.X(), disp.Y(), 0).Mul4(model)
	return wc.camera.GetViewMatrix().Mul4(model)
	//return wc.camera.GetTransform().Inv().Mul4(model)
}

func (wc worldContextImpl) ModelToModelViewProjection(
	model mgl32.Mat4,
) mgl32.Mat4 {
	return wc.Projection().Mul4(wc.ModelToModelView(model))
}

func NewWorldRenderContext(
	camera *camera.Camera, planet *world.Planet,
) render.WorldContext {
	return worldContextImpl{camera: camera, planet: planet}
}
