package cxmath

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

const eps = 0.01

func TestAngleToPositive(t *testing.T) {
	v1 := mgl32.Vec2{0,0}
	v2 := mgl32.Vec2{1,1}

	expected := -math.Pi/4
	got := float64(AngleTo(v1,v2))

	if math.Abs(expected-got)>eps {
		t.Errorf("expected %v; got %v",expected,got)
	}
}

func TestAngleToNegative(t *testing.T) {
	v1 := mgl32.Vec2{0,0}
	v2 := mgl32.Vec2{1,-1}

	expected := math.Pi/4
	got := float64(AngleTo(v1,v2))

	if math.Abs(expected-got)>eps {
		t.Errorf("expected %v; got %v",expected,got)
	}
}

