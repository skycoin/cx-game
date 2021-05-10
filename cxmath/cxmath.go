package cxmath

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/spriteloader"
)

func ConvertScreenCoordsToWorld(x,y float32, projection mgl32.Mat4) mgl32.Vec4 {
	homogenousClipCoords := mgl32.Vec4 { x,y,-1.0,1.0}
	cameraCoords := projection.Inv().Mul4x1(homogenousClipCoords)
	rayEye := mgl32.Vec4 { cameraCoords.X(), cameraCoords.Y(), -1.0, 0 }
	worldCoords := rayEye.Mul(spriteloader.SpriteRenderDistance)
	worldCoords[3]=1
	return worldCoords
}

func Scale(factor float32) mgl32.Mat4 {
	return mgl32.Scale3D(factor,factor,factor)
}
