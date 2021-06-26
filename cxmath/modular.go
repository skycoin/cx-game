package cxmath

// 1D modular wrap-around arithmetic
// Useful for planet wrapping calculations.

import (
	"math"
)


type Modular struct {
	modulo float32
}

func NewModular(modulo float32) Modular {
	return Modular { modulo: modulo }
}

// positive floating point modulo
func pfmod(x float32, m float32) float32 {
	result := float32(math.Mod(float64(x),float64(m)))
	if result < 0 {
		result += m
	}
	return result
}

func (m Modular) Mod(x float32) float32 {
	return pfmod(x,m.modulo)
}

// compute the SHORTEST displacement between two points
func (m Modular) Disp(x1,x2 float32) float32 {
	result := pfmod(x2-x1,m.modulo)
	if result > m.modulo/2 {
		result -= m.modulo
	}
	return result
}

// is x to the left of y on a number line?
func (m Modular) IsLeft(x,y float32) bool {
	return m.Disp(y,x) < 0
}

// is x to the right of y on a number line?
func (m Modular) IsRight(x,y float32) bool {
	return m.Disp(y,x) > 0
}
