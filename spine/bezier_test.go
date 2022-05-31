package spine

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

type bezierControl struct {
	P1 Vector
	P2 Vector
}

func (_ bezierControl) Generate(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(bezierControl{
		P1: Vector{
			X: r.Float32(),
			Y: r.Float32(),
		},
		P2: Vector{
			X: r.Float32(),
			Y: r.Float32(),
		},
	})
}

func TestBezier01D(t *testing.T) {
	const MaxDeviation = 0.1

	err := quick.Check(func(bezier bezierControl) bool {
		p1, p2 := bezier.P1, bezier.P2
		bez := NewBezier01D(p1, p2)
		for time := float32(0.0); time <= 1.0; time += 0.1 {
			p := Bezier01(p1, p2, time)
			r := bez.Y(p.X)
			if abs(r-p.Y) > MaxDeviation {
				t.Logf("D: %v", bez)
				t.Logf("delta %v", r-p.Y)
				t.Logf("B(%v, %v, %f) = %v; bez=%f", p1, p2, time, p, r)
				return false
			}
		}
		return true
	}, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestBezier01T(t *testing.T) {
	const MaxDeviation = 0.3

	err := quick.Check(func(bezier bezierControl) bool {
		p1, p2 := bezier.P1, bezier.P2
		if p1.Y < 0 || p1.Y > 1 {
			panic("invalid input: " + fmt.Sprintf("%f", p1.Y))
		}
		bez := NewBezier01T(p1, p2)
		for time := float32(0.0); time <= 1.0; time += 0.1 {
			p := Bezier01(p1, p2, time)
			r := bez.Y(p.X)
			if abs(r-p.Y) > MaxDeviation {
				t.Logf("D: %v", bez)
				t.Logf("delta %v", r-p.Y)
				t.Logf("B(%v, %v, %f) = %v; bez=%f", p1, p2, time, p, r)
				return false
			}
		}
		return true
	}, nil)

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkBezier01D(b *testing.B) {
	bez := NewBezier01D(Vector{0.25, 0.75}, Vector{0.75, 0.25})
	for i := 0; i < b.N; i++ {
		for time := float32(0.0); time < 1.0; time += 0.1 {
			_ = bez.Y(time)
		}
	}
}

func BenchmarkBezier01T(b *testing.B) {
	bez := NewBezier01T(Vector{0.25, 0.75}, Vector{0.75, 0.25})
	for i := 0; i < b.N; i++ {
		for time := float32(0.0); time < 1.0; time += 0.1 {
			_ = bez.Y(time)
		}
	}
}
