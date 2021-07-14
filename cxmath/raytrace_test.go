package cxmath

import (
	"testing"
)

func pointsEqual(got, expected []Vec2i) bool {
	if len(got) != len(expected) {
		return false
	}
	for i := 0; i < len(got); i++ {
		if got[i] != expected[i] {
			return false
		}
	}
	return true
}

func TestRaytrace(t *testing.T) {
	x0 := 0.5
	y0 := 0.5
	x1 := 2.5
	y1 := 3.5

	expected := []Vec2i{
		{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}, {2, 3},
	}
	intersects := Raytrace(x0, y0, x1, y1)

	if !pointsEqual(intersects, expected) {
		t.Errorf("expected %v; got %v", expected, intersects)
	}
}
