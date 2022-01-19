package testClearColor

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type TestClearColor struct {
	M_ClearColor [4]float32
}

func SetUpTestClearColor() *TestClearColor {
	test := &TestClearColor{
		M_ClearColor: [4]float32{0.2, 0.3, 0.8, 1.0},
	}
	return test
}

func (test *TestClearColor) deleteTest() {

}

func (test *TestClearColor) OnUpdate(deltaTime float32) {

}
func (test *TestClearColor) OnRender() {
	gl.ClearColor(test.M_ClearColor[0], test.M_ClearColor[1], test.M_ClearColor[2], test.M_ClearColor[3])
	gl.CLEAR(gl.COLOR_BUFFER_BIT)
}
func (test *TestClearColor) OnUIRender() {

}
