package spine

import (
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	skeleton, err := ReadJSON(strings.NewReader(BoxSkeletonJSON))
	if err != nil {
		t.Error(err)
	}

	testFloat(t, "size.X", skeleton.Size.X, 138.56)
	testFloat(t, "size.Y", skeleton.Size.Y, 186.55)
	testInt(t, "len(bones)", len(skeleton.Bones), 13)

	for _, bone := range skeleton.Bones {
		if bone.Parent == nil && bone.Name != "root" {
			t.Errorf("unbound bone parent: %v", bone.Name)
		}
	}

	// { "name": "back-shin", "parent": "back-thigh", "length": 31.17, "rotation": -50.15, "x": 44.3, "y": 0.06, "color": "ff0008ff" },
	ankleLeft := skeleton.FindBone("ankle-left")
	if ankleLeft == nil {
		t.Error("left-left missing")
	} else {
		testFloat(t, "ankle-left.length", ankleLeft.Length, 70)
		testFloat(t, "ankle-left.rotation", ankleLeft.Local.Rotate, 1.9326031)
		testFloat(t, "ankle-left.x", ankleLeft.Local.Translate.X, 70)
		testFloat(t, "ankle-left.y", ankleLeft.Local.Translate.Y, 0)
		// TODO: testColor
	}

	testInt(t, "len(slots)", len(skeleton.Slots), 9)
	for _, slot := range skeleton.Slots {
		if slot.Name == "" || slot.Bone == nil {
			t.Errorf("invalid slot: %v", slot)
		}
	}

	testInt(t, "len(events)", len(skeleton.Events), 1)
	testInt(t, "len(skins)", len(skeleton.Skins), 3)
	testInt(t, "len(animations)", len(skeleton.Animations), 1)
	testInt(t, "len(animations[0].Timelines)", len(skeleton.Animations[0].Timelines), 15)
}

func testFloat(t *testing.T, fmt string, got, exp float32) {
	t.Helper()
	if got != exp {
		t.Errorf("invalid %v: got %v exp %v", fmt, got, exp)
	}
}

func testInt(t *testing.T, fmt string, got, exp int) {
	t.Helper()
	if got != exp {
		t.Errorf("invalid %v: got %v exp %v", fmt, got, exp)
	}
}

const BoxSkeletonJSON = `
{
"skeleton": { "hash": "se5KLTMjFALRa3sLdCKOmLx83Tk", "spine": "3.6.38", "width": 138.56, "height": 186.55, "images": "./images/" },
"bones": [
	{ "name": "root", "x": 0.8 },
	{ "name": "body", "parent": "root", "length": 200, "rotation": 90, "x": -1.96, "y": 34.25, "scaleX": 0.4, "scaleY": 0.4 },
	{ "name": "leg-left", "parent": "body", "length": 70, "rotation": 81.69 },
	{ "name": "ankle-left", "parent": "leg-left", "length": 70, "rotation": 110.73, "x": 70 },
	{ "name": "leg-right", "parent": "body", "length": 70, "rotation": -76.25 },
	{ "name": "ankle-right", "parent": "leg-right", "length": 70, "rotation": -108.52, "x": 70 },
	{ "name": "elbow-left", "parent": "body", "length": 90, "rotation": -105.25, "x": 70 },
	{ "name": "elbow-right", "parent": "body", "length": 90, "rotation": -90, "x": 70 },
	{ "name": "eye", "parent": "body", "length": 48.53, "rotation": -90, "x": 165, "y": -75 },
	{ "name": "hand-left", "parent": "elbow-left", "length": 70.42, "x": 90 },
	{ "name": "hand-right", "parent": "elbow-right", "length": 70, "x": 90 },
	{
		"name": "hat",
		"parent": "body",
		"length": 66.53,
		"rotation": 78.5,
		"x": 235.09,
		"y": -4.09,
		"transform": "noRotationOrReflection"
	},
	{ "name": "pupil", "parent": "eye", "length": 23.3, "x": 50 }
],
"slots": [
	{ "name": "hand-left", "bone": "hand-left", "color": "aaaaaaff", "attachment": "limb" },
	{ "name": "ankle-left", "bone": "ankle-left", "color": "aaaaaaff", "attachment": "limb" },
	{ "name": "body", "bone": "body", "attachment": "body" },
	{ "name": "ankle-right", "bone": "ankle-right", "attachment": "limb" },
	{ "name": "eye", "bone": "eye", "attachment": "eye-socket" },
	{ "name": "pupil", "bone": "pupil", "attachment": "eye-pupil" },
	{ "name": "hat", "bone": "hat", "attachment": "hat" },
	{ "name": "hand-grab", "bone": "hand-right" },
	{ "name": "hand-right", "bone": "hand-right", "attachment": "limb" }
],
"skins": {
	"default": {
		"ankle-left": {
			"limb": { "x": 70, "rotation": 135, "width": 75, "height": 75 }
		},
		"ankle-right": {
			"limb": { "x": 70, "rotation": 45, "width": 75, "height": 75 }
		},
		"body": {
			"body": { "x": 105.83, "rotation": -90, "width": 292, "height": 292 }
		},
		"eye": {
			"eye-socket": { "x": 48.74, "y": 0.03, "width": 102, "height": 102 }
		},
		"hand-grab": {
			"mace": { "x": 67.89, "y": 97.45, "width": 87, "height": 213 },
			"sword": { "x": 70.21, "y": 104.28, "width": 75, "height": 177 }
		},
		"hand-left": {
			"limb": { "x": 70, "width": 75, "height": 75 }
		},
		"hand-right": {
			"limb": { "x": 70, "width": 75, "height": 75 }
		},
		"pupil": {
			"eye-pupil": { "x": 25, "width": 38, "height": 38 }
		}
	},
	"plain": {
		"hat": {
			"hat": { "name": "plain/hat", "x": 51.3, "rotation": -90, "width": 142, "height": 116 }
		}
	},
	"purple": {
		"hat": {
			"hat": { "name": "purple/hat", "x": 51.3, "rotation": -90, "width": 142, "height": 116 }
		}
	}
},
"events": {
	"step": {}
},
"animations": {
	"run": {
		"bones": {
			"ankle-right": {
				"rotate": [
					{ "time": 0, "angle": 30.12 },
					{ "time": 0.0667, "angle": 49.87 },
					{ "time": 0.3333, "angle": 63.25 },
					{ "time": 0.4667, "angle": 20.4 },
					{ "time": 0.8, "angle": 6.23 },
					{ "time": 1.1333, "angle": 30.12 }
				]
			},
			"leg-right": {
				"rotate": [
					{ "time": 0, "angle": -10.97 },
					{ "time": 0.0667, "angle": -47.87 },
					{ "time": 0.3667, "angle": 230.51 },
					{ "time": 0.5, "angle": 225.19 },
					{ "time": 0.7667, "angle": -52.55 },
					{ "time": 0.9, "angle": -20.29 },
					{ "time": 1.1333, "angle": -10.97 }
				]
			},
			"leg-left": {
				"rotate": [
					{ "time": 0, "angle": 59.76 },
					{ "time": 0.1333, "angle": -206.66 },
					{ "time": 0.2333, "angle": -197.25 },
					{ "time": 0.4, "angle": -186.35 },
					{ "time": 0.9333, "angle": 78.07 },
					{ "time": 1.1333, "angle": 59.76 }
				]
			},
			"ankle-left": {
				"rotate": [
					{ "time": 0, "angle": -207.81 },
					{ "time": 0.1667, "angle": -204.28 },
					{ "time": 0.2667, "angle": -153.01 },
					{ "time": 0.3667, "angle": -127.07 },
					{ "time": 0.9333, "angle": -179.8 },
					{ "time": 1.1333, "angle": -173.3 }
				]
			},
			"body": {
				"rotate": [
					{ "time": 0, "angle": -4.03 },
					{ "time": 0.2333, "angle": -0.74 },
					{ "time": 0.6667, "angle": -14.61 },
					{ "time": 0.9, "angle": -1.77 },
					{ "time": 1.1333, "angle": -4.03 }
				],
				"translate": [
					{ "time": 0, "x": 0, "y": 12.43 },
					{ "time": 0.0667, "x": 0, "y": 26.26 },
					{ "time": 0.1667, "x": 0, "y": 27.64 },
					{ "time": 0.2667, "x": 0, "y": 26.52 },
					{ "time": 0.3333, "x": 0, "y": 16.01 },
					{ "time": 0.4333, "x": 0, "y": 17.67 },
					{ "time": 0.5, "x": 0, "y": 25.42 },
					{ "time": 0.6333, "x": 0, "y": 29.57 },
					{ "time": 0.7667, "x": 0, "y": 30.13 },
					{ "time": 0.8667, "x": 0, "y": 20.72 },
					{ "time": 0.9667, "x": 0, "y": 20.44 },
					{ "time": 1.1333, "x": 0, "y": 6.62 }
				]
			},
			"hand-right": {
				"rotate": [
					{ "time": 0, "angle": 1.81 },
					{ "time": 0.1667, "angle": 95.69 },
					{ "time": 0.7333, "angle": 134.54 },
					{ "time": 0.9333, "angle": 90.21 },
					{ "time": 1.1333, "angle": 1.81 }
				]
			},
			"elbow-right": {
				"rotate": [
					{ "time": 0, "angle": 0.64 },
					{ "time": 0.1333, "angle": -75.48 },
					{ "time": 0.4333, "angle": 206.51 },
					{ "time": 0.7333, "angle": 199.04 },
					{ "time": 0.8667, "angle": -65.13 },
					{ "time": 1.1333, "angle": 0.64 }
				]
			},
			"elbow-left": {
				"rotate": [
					{ "time": 0, "angle": 219.42 },
					{ "time": 0.4, "angle": -12.04 },
					{ "time": 0.9667, "angle": 218.13 },
					{ "time": 1.1333, "angle": 219.42 }
				]
			},
			"hand-left": {
				"rotate": [
					{ "time": 0, "angle": 147.57 },
					{ "time": 0.3333, "angle": 48.99 },
					{ "time": 0.8667, "angle": 151.73 },
					{ "time": 1.1333, "angle": 147.57 }
				]
			},
			"eye": {
				"rotate": [
					{ "time": 0, "angle": -6.26 },
					{ "time": 0.2333, "angle": -13.11 },
					{ "time": 0.5667, "angle": 11.4 },
					{ "time": 0.7667, "angle": 22.75 },
					{ "time": 0.9, "angle": -3.7 },
					{ "time": 1.0667, "angle": 12.91 },
					{ "time": 1.1333, "angle": -6.26 }
				]
			},
			"hat": {
				"rotate": [
					{ "time": 0, "angle": 13.4 },
					{ "time": 0.1, "angle": 12.29 },
					{ "time": 0.2, "angle": 17.77 },
					{ "time": 0.3, "angle": 5.29 },
					{ "time": 0.4, "angle": -3.16 },
					{ "time": 0.5, "angle": 7.46 },
					{ "time": 0.7, "angle": -15.64 },
					{ "time": 0.9, "angle": -15.03 },
					{ "time": 0.9667, "angle": 1.94 },
					{ "time": 1.0667, "angle": 13.4 }
				],
				"translate": [
					{
						"time": 0,
						"x": 0,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.1667,
						"x": -0.75,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.2667,
						"x": 29.66,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.3667,
						"x": 7.33,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.5,
						"x": -2.95,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.6,
						"x": 12.08,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.7,
						"x": 33.47,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 0.8667,
						"x": -3.9,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{
						"time": 1,
						"x": -4.75,
						"y": 0,
						"curve": [ 0.541, 0, 0.543, 1 ]
					},
					{ "time": 1.1333, "x": 0, "y": 0 }
				],
				"scale": [
					{ "time": 0.1667, "x": 0.949, "y": 1 },
					{ "time": 0.2667, "x": 0.745, "y": 1 },
					{ "time": 0.3667, "x": 1.117, "y": 1 },
					{ "time": 0.5, "x": 0.838, "y": 1 },
					{ "time": 0.6, "x": 0.926, "y": 1 },
					{ "time": 0.7, "x": 0.858, "y": 1 },
					{ "time": 0.8667, "x": 1.226, "y": 1 },
					{ "time": 1, "x": 1, "y": 1 },
					{ "time": 1.1333, "x": 0.965, "y": 1 }
				]
			}
		},
		"events": [
			{ "time": 0.0333, "name": "step" },
			{ "time": 0.4667, "name": "step" }
		]
	}
}
}`
