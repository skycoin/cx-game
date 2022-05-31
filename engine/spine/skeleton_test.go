package spine

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestSkeletonUpdateTransform(t *testing.T) {
	data, err := ReadJSON(strings.NewReader(`{
		"skeleton": {"width": 100, "height": 100},
		"skins": { "default": {} },
		"bones": [
			{ "name": "root", "y": 10},
			{ "name": "U", "parent": "root", "y": 10 },
			{ "name": "UU", "parent": "U", "y": 10 },
			{ "name": "UUU", "parent": "UU", "y": 10 },
			{ "name": "UUUC", "parent": "UUU", "rotation": 90 },
			{ "name": "UUUW", "parent": "UUU", "rotation": -90 },
			{ "name": "UUUCU", "parent": "UUUC", "y": 10 },
			{ "name": "UUUWU", "parent": "UUUW", "y": 10 },
			{ "name": "UUUS", "parent": "UUU", "scaleX": 0.5, "scaleY": 0.5 },
			{ "name": "UUUSU", "parent": "UUUS", "y": 10 }
		]
	}`))
	if err != nil {
		t.Fatal(err)
	}

	skeleton := NewSkeleton(data)
	skeleton.Update()

	testPoint(t, Vector{0, 10}, skeleton.FindBone("root").World.Transform(V0))
	testPoint(t, Vector{0, 20}, skeleton.FindBone("U").World.Transform(V0))
	testPoint(t, Vector{0, 30}, skeleton.FindBone("UU").World.Transform(V0))
	testPoint(t, Vector{0, 40}, skeleton.FindBone("UUU").World.Transform(V0))

	testPoint(t, Vector{-10, 40}, skeleton.FindBone("UUUCU").World.Transform(V0))
	testPoint(t, Vector{10, 40}, skeleton.FindBone("UUUWU").World.Transform(V0))
	testPoint(t, Vector{0, 45}, skeleton.FindBone("UUUSU").World.Transform(V0))
}

func TestAnimateTransform(t *testing.T) {
	data, err := ReadJSON(strings.NewReader(`{
		"skeleton": {"width": 100, "height": 100},
		"bones": [
			{ "name": "root", "y": 10},
			{ "name": "U", "parent": "root", "y": 10 },
			{ "name": "UU", "parent": "U", "y": 10 },
			{ "name": "UUU", "parent": "UU", "y": 10 },
			{ "name": "UUUC", "parent": "UUU", "rotation": 90 },
			{ "name": "UUUCU", "parent": "UUUC", "y": 10 }
		],
		"skins": { "default": {} },
		"animations": {
			"run": {
				"bones": {
					"U": {
						"translate": [
							{ "time": 0, "y": -10 },
							{ "time": 1, "y": 0 },
							{ "time": 2, "y": 10 }
						]
					},
					"UUUC": {
						"rotate": [
							{ "time": 0, "angle": 0  },
							{ "time": 1, "angle": 180 },
							{ "time": 2, "angle": 0  }
						]
					}
				}
			}
		}
	}`))

	if err != nil {
		t.Fatal(err)
	}

	skeleton := NewSkeleton(data)
	skeleton.Update()

	uuu := skeleton.FindBone("UUU")
	uuucu := skeleton.FindBone("UUUCU")

	testPoint(t, Vector{0, 40}, uuu.World.Transform(V0))
	testPoint(t, Vector{-10, 40}, uuucu.World.Transform(V0))

	animation := skeleton.Data.Animations[0]
	fmt.Println(animation)

	animation.Apply(skeleton, 0, true)
	skeleton.Update()

	testPoint(t, Vector{0, 30}, uuu.World.Transform(V0))
	testPoint(t, Vector{-10, 30}, uuucu.World.Transform(V0))

	animation.Apply(skeleton, 0.5, true)
	skeleton.Update()

	testPoint(t, Vector{0, 35}, uuu.World.Transform(V0))
	testPoint(t, Vector{0, 25}, uuucu.World.Transform(V0))

	animation.Apply(skeleton, 1.0, true)
	skeleton.Update()

	testPoint(t, Vector{0, 40}, uuu.World.Transform(V0))
	testPoint(t, Vector{10, 40}, uuucu.World.Transform(V0))

	animation.Apply(skeleton, 1.5, true)
	skeleton.Update()

	testPoint(t, Vector{0, 45}, uuu.World.Transform(V0))
	testPoint(t, Vector{0, 35}, uuucu.World.Transform(V0))

	animation.Apply(skeleton, 1.99999, true)
	skeleton.Update()

	testPoint(t, Vector{0, 50}, uuu.World.Transform(V0))
	testPoint(t, Vector{-10, 50}, uuucu.World.Transform(V0))

	animation.Apply(skeleton, 2.5, true)
	skeleton.Update()

	testPoint(t, Vector{0, 35}, uuu.World.Transform(V0))
	testPoint(t, Vector{0, 25}, uuucu.World.Transform(V0))
}

func TestAnimateColor(t *testing.T) {
	data, err := ReadJSON(strings.NewReader(`{
		"skeleton": {"width": 100, "height": 100},
		"bones": [
			{ "name": "root", "y": 10}
		],
		"slots": [
			{ "name": "star", "bone": "root", "attachment": "star", "color": "ff000000"}
		],
		"skins": {
			"default": {
				"star": {
					"star": {"width": 10, "height": 10}
				}
			}
		},
		"animations": {
			"run": {
				"slots": {
					"star": {
						"color": [
							{ "time": 0, "color": "ffffff" },
							{ "time": 1, "color": "ff00ff80" },
							{ "time": 2, "color": "0000ff00" }
						]
					}
				}
			}
		}
	}`))

	if err != nil {
		t.Fatal(err)
	}

	skeleton := NewSkeleton(data)
	skeleton.Update()

	slot := skeleton.FindSlot("star")
	if slot == nil {
		t.Fatal("unable to find slot")
	}
	testColor(t, RGBA(1, 0, 0, 0), slot.Color)

	animation := skeleton.Data.Animations[0]
	fmt.Println(animation.Timelines[0])

	animation.Apply(skeleton, 0, true)
	skeleton.Update()

	testColor(t, RGBA(1, 1, 1, 1), slot.Color)

	animation.Apply(skeleton, 0.5, true)
	skeleton.Update()

	testColor(t, RGBA(1, 0.5, 1, 0.75), slot.Color)

	animation.Apply(skeleton, 1.0, true)
	skeleton.Update()

	testColor(t, RGBA(1, 0, 1, 0.5), slot.Color)

	animation.Apply(skeleton, 1.5, true)
	skeleton.Update()

	testColor(t, RGBA(0.5, 0, 1, 0.25), slot.Color)

	animation.Apply(skeleton, 1.99999, true)
	skeleton.Update()

	testColor(t, RGBA(0, 0, 1, 0.0), slot.Color)

	animation.Apply(skeleton, 2.5, true)
	skeleton.Update()

	testColor(t, RGBA(1, 0.5, 1, 0.75), slot.Color)
}

func TestTransform(t *testing.T) {
	x := NewTransform()
	x.Translate.Y = 10
	x.Rotate = math.Pi / 2

	aff := x.Affine()
	testPoint(t, Vector{-1, 10}, aff.Transform(Vector{0, 1}))
}

func testPoint(t *testing.T, exp, got Vector) {
	t.Helper()
	if got.Sub(exp).Len() > 0.001 {
		t.Errorf("exp %v got %v", exp, got)
	}
}

func testColor(t *testing.T, exp, got Color) {
	t.Helper()

	delta := abs(exp.R-got.R) + abs(exp.G-got.G) + abs(exp.B-got.B) + abs(exp.A-got.A)
	if delta > 0.01 {
		t.Errorf("delta %v: exp %v got %v", delta, exp, got)
	}
}
