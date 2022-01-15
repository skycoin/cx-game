package renderer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type GL_ERROR_CODE int32

const (
	NONE GL_ERROR_CODE = 0
)

func GLClearError() {
	for gl.GetError() != gl.NO_ERROR {
		fmt.Println("Cleared an Error")
	}
}

func GLCheckError() {
	for err := gl.GetError(); err != 0; {
		fmt.Println("got a An Error: ", err)
	}
}
