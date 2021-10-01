package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	shockwave_center    mgl32.Vec2 = mgl32.Vec2{0.5, 0.5}
	shockwave_force     float32    = 0.2
	shockwave_size      float32    = 0.3
	shockwave_thickness float32    = 0
	shockwave_enabled   bool       = false
)

func SetupShockwave() {
	// fmt.Println("SHOCKWAVE ENABLED: ", shockwave_enabled)
	// EnableShockwave()
	ScreenShader.SetBool("u_shockwave_enabled", shockwave_enabled)
	ScreenShader.SetVec2("data.center", &shockwave_center)
	ScreenShader.SetFloat("data.force", shockwave_force)
	ScreenShader.SetFloat("data.size", shockwave_size)
	ScreenShader.SetFloat("data.thickness", shockwave_thickness)
}

func SetShockwaveSize(size float32) {
	shockwave_size = size
}

func SetShockwaveForce(force float32) {
	shockwave_force = force
}

func SetShockwaveThickness(thickness float32) {
	shockwave_thickness = thickness
}

func SetShockwaveCenter(x, y float32) {
	shockwave_center = mgl32.Vec2{x, y}
}

func EnableShockwave() {
	shockwave_enabled = true
}

func DisableShockwave() {
	shockwave_enabled = false
}
