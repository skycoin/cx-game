package common

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func GlAssert(glEnum uint32, truth bool) {
	var value bool
	gl.GetBooleanv(glEnum, &value)
	if value != truth {
		_, fn, line, _ := runtime.Caller(1)
		log.Fatalf("ASSERT FAILED on line: %s:%d\n", fn, line)
	}
}
