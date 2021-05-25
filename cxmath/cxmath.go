package cxmath

import (
	"github.com/go-gl/mathgl/mgl32"
	//"log"
)

func ConvertScreenCoordsToWorld(x,y float32, projection mgl32.Mat4) mgl32.Vec2 {
	homogenousClipCoords := mgl32.Vec4 { x,y,-1.0,1.0}
	cameraCoords := projection.Inv().Mul4x1(homogenousClipCoords)
	return cameraCoords.Vec2()
}

func Scale(factor float32) mgl32.Mat4 {
	return mgl32.Scale3D(factor,factor,factor)
}

func IntMax(x,y int) int {
	if x>y {
		return x
	} else {
		return y
	}
}

func IntMin(x,y int) int {
	if x<y {
		return x
	} else {
		return y
	}
}
