package render

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	drawBBoxLines     bool = false
	AllLines          []float32
	AllCollidingLines []float32
)

func DrawBBoxLines(lines, collidingLines []float32) {
	AllLines = append(AllLines, lines...)
	AllCollidingLines = append(AllCollidingLines, collidingLines...)
}

func flushBBoxLineDraws(projection mgl32.Mat4) {
	mvp := projection.Mul4(cameraTransform.Inv())
	DrawLines(
		AllLines,
		mgl32.Vec3{0, 0, 1},
		mvp,
	)
	//without this check crashes
	if len(AllCollidingLines) >= 4 {
		DrawLines(AllCollidingLines, mgl32.Vec3{1, 0, 0}, mvp)
	}
	AllLines = AllLines[:0]
	AllCollidingLines = AllCollidingLines[:0]
}

func ToggleBBox() {
	drawBBoxLines = !drawBBoxLines

	toggleStatus := "not active"
	if drawBBoxLines {
		toggleStatus = "active"
	}

	fmt.Printf("Toggled Bounding Box: %v\n", toggleStatus)
}

func IsBBoxActive() bool {
	return drawBBoxLines
}
