package cxmath

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"
)

func ConvertScreenCoordsToWorld(x,y float32, projection mgl32.Mat4) mgl32.Vec2 {
	homogenousClipCoords := mgl32.Vec4 { x,y,-1.0,1.0}
	cameraCoords := projection.Inv().Mul4x1(homogenousClipCoords)
	return cameraCoords.Vec2()
}

func Scale(factor float32) mgl32.Mat4 {
	return mgl32.Scale3D(factor,factor,factor)
}

func atan2f32(y,x float32) float32 {
	return float32(math.Atan2(float64(y),float64(x)))
}

func AngleTo(v1,v2 mgl32.Vec2) float32 {
	return atan2f32(v1.Y(),v1.X()) - atan2f32(v2.Y(),v2.X())
}
