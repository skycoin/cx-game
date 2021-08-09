package anim

import (
	"log"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader"
	animjson "github.com/skycoin/cx-game/engine/spriteloader/anim/json"
)

func parseKeyframeName(name string) (action string, index int) {
	start := strings.Index(name, "#") + 1
	end := strings.LastIndex(name, ".")
	middle := name[start:end]

	words := strings.Split(middle, " ")
	action = words[0]
	haveIndex := len(words) == 2
	if haveIndex {
		var err error
		index, err = strconv.Atoi(words[1])
		if err != nil {
			log.Fatalf("could not parse keyframe name [%s]", name)
		}
	} else {
		index = 0
	}
	return
}

func createActionsFromFrameTags(
	frameTags []animjson.FrameTagItem,
) map[string]Action {
	actions := make(map[string]Action)
	for _, frameTag := range frameTags {
		numFrames := frameTag.To - frameTag.From + 1
		actions[frameTag.Name] = Action{
			Name:   frameTag.Name,
			Frames: make([]Frame, numFrames),
		}
	}
	return actions
}

func populateActions(
	actions map[string]Action, keyframes map[string]animjson.KeyFrame,
	sheetDims cxmath.Vec2i,
) {
	for name, keyframe := range keyframes {
		action, index := parseKeyframeName(name)
		frame := keyframe.Frame

		offsetX := float32(frame.X) / float32(sheetDims.X)
		offsetY := float32(frame.Y) / float32(sheetDims.Y)

		scaleX := float32(frame.W) / float32(sheetDims.X)
		scaleY := float32(frame.H) / float32(sheetDims.Y)

		actions[action].Frames[index] = Frame{
			Duration: float32(keyframe.Duration) / 1000,
			Transform: mgl32.Translate2D(offsetX, offsetY).
				Mul3(mgl32.Scale2D(scaleX, scaleY)),
		}
	}
}

func LoadAnimationFromJSON(fname string) Animation {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("could not find animation spritesheet at %s", fname)
	}
	sheet := animjson.UnmarshalAnimatedSpriteSheet(buf)

	imgPath := "./assets/" + sheet.Meta.Image
	gpuTex := spriteloader.LoadTextureFromFileToGPU(imgPath)

	actions := createActionsFromFrameTags(sheet.Meta.FrameTags)
	populateActions(actions, sheet.Frames, gpuTex.Dims())

	return Animation{
		Texture: gpuTex.Gl,
		Actions: actions,
	}
}

func (anim *Animation) Action(name string) Action {
	for _, action := range anim.Actions {
		if action.Name == name {
			return action
		}
	}
	log.Fatalf("cannot find animation action with name [%s]", name)
	return Action{}
}

func (anim *Animation) NewPlayback(idleActionName string) Playback {
	action := anim.Action(idleActionName)
	// set up repeating
	return Playback{Animation: *anim, Action: action, NextAction: action}
}
