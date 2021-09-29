package cxmath

import "github.com/go-gl/mathgl/mgl32"

func Vec3Mix(a, b mgl32.Vec3, f float32) mgl32.Vec3 {
	return Vec3Add(a.Mul(1.0-f), b.Mul(f))
}

func Vec3Add(a, b mgl32.Vec3) mgl32.Vec3 {
	return a.Add(b)
}

func Vec3ScalarMult(v mgl32.Vec3, scalar float32) mgl32.Vec3 {
	return v.Mul(scalar)
}
func Vec3ScalarAdd(v mgl32.Vec3, scalar float32) mgl32.Vec3 {
	return v.Add(mgl32.Vec3{scalar, scalar, scalar})
}
