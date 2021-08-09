package ui

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	//"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
)

const dialogueScale = 0.3
const dialogueFadeTime = 0.5

type DialogueBox struct {
	WorldTransform mgl32.Mat4
	Time           float32
	Text           string
	Alignment      TextAlignment
}

var dialogueBoxes = []DialogueBox{}

func PlaceDialogueBox(
	Text string, Alignment TextAlignment, Time float32,
	WorldTransform mgl32.Mat4,
) {
	dialogueBoxes = append(dialogueBoxes, DialogueBox{
		WorldTransform: WorldTransform,
		Time:           Time,
		Alignment:      Alignment,
		Text:           Text,
	})
}

func TickDialogueBoxes(dt float32) {

	oldDialogueBoxes := dialogueBoxes
	newDialogueBoxes := []DialogueBox{}

	for _, box := range oldDialogueBoxes {
		box.Time -= dt
		if box.Time > 0 {
			newDialogueBoxes = append(newDialogueBoxes, box)
		}
	}

	dialogueBoxes = newDialogueBoxes
}

func DrawDialogueBoxes(ctx render.Context) {
	for _, box := range dialogueBoxes {
		boxLocalTransform := box.WorldTransform.
			Mul4(cxmath.Scale(dialogueScale))

		boxCtx := ctx.PushLocal(boxLocalTransform)
		/*
			modelViewMatrix :=
				mgl32.Translate3D(-cam.X,-cam.Y,0).Mul4(box.WorldTransform).
				Mul4(cxmath.Scale(dialogueScale))
		*/
		// TODO fade out
		opacity := float32(math.Min(float64(box.Time/dialogueFadeTime), 1))
		color := mgl32.Vec4{1, 1, 1, opacity}
		DrawString(box.Text, color, box.Alignment, boxCtx)
	}
}
