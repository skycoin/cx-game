package inputhandler

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func BenchmarkEventType(b *testing.B) {
	//arrange
	button1 := glfw.KeySpace
	// button2 := glfw.MouseButton1

	//act
	for i := 0; i < b.N; i++ {
		_ 	= int(button1)
		// back := glfw.Key(converted)
		// back = back + 1
	}

	//assert
}
