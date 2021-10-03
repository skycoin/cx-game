package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/engine/input"
)

var (
	shockwave_center    mgl32.Vec2 = mgl32.Vec2{0.5, 0.5}
	shockwave_force     float32    = 0.2
	shockwave_size      float32    = 0.3
	shockwave_thickness float32    = 0
	shockwave_enabled   bool       = true
	inverted            bool       = false

	slider_size          float32 = 0
	slider_thickness     float32 = 0
	slider_thickness_gap float32 = 0.1
	slider_force         float32 = 0.2
	mask_flag            bool    = false
)

//size from 0 to 0.23
//

func InitShockwave() {
	input.RegisterSlider(
		&slider_size, 0.01, 0, 0.5, glfw.Key1, glfw.Key2, glfw.ModShift, "slider_size",
	)
	input.RegisterSlider(
		&slider_thickness, 0.01, 0, 1, glfw.Key3, glfw.Key4, glfw.ModShift, "slider_thickness",
	)
	input.RegisterSlider(
		&slider_force, 0.01, -1, 1, glfw.Key1, glfw.Key2, glfw.ModControl, "slider_force",
	)

	input.RegisterSlider(
		&slider_thickness_gap, 0.01, 0, 1, glfw.Key3, glfw.Key4, glfw.ModControl, "slider_thickness_gap",
	)
	input.RegisterFlag(func() {
		mask_flag = true
	}, func() {
		mask_flag = false
	}, glfw.KeyN, 0, "mask")

	input.RegisterFlag(func() { inverted = true }, func() { inverted = false }, glfw.KeyZ, 0, "inverted colors")
}

func SetupShockwave() {
	// fmt.Println("SHOCKWAVE ENABLED: ", shockwave_enabled)
	// EnableShockwave()
	ScreenShader.SetBool("u_shockwave_enabled", shockwave_enabled)
	ScreenShader.SetVec2("data.center", &shockwave_center)
	ScreenShader.SetFloat("data.force", slider_force)
	ScreenShader.SetFloat("data.size", slider_size)
	ScreenShader.SetFloat("data.thickness", slider_thickness)
	ScreenShader.SetFloat("data.thickness_gap", slider_thickness_gap)

	ScreenShader.SetBool("inverted", inverted)
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

func EnableInverted() {
	inverted = true
}

func DisableInverted() {
	inverted = false
}
