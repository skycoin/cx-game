package anim

import (
	"github.com/skycoin/cx-game/render"
)

var Program render.Program

func InitAnimatedSpriteLoader() {
	Program = render.CompileProgram(
		"./assets/shader/mvp.vert", "./assets/shader/anim.frag",
	)
}

func (p *Playback) PlayOnce(action string) {
	p.NextAction = p.Animation.Action(action)
	p.Repeat = false
}

func (p *Playback) PlayRepeating(action string) {
	p.NextAction = p.Animation.Action(action)
	p.Repeat = true
}

func (p *Playback) Frame() Frame {
	return p.Action.Frames[p.FrameIndex]
}

func (p *Playback) Update(dt float32) {
	p.SecondsIntoFrame += dt

	for p.SecondsIntoFrame > p.Frame().Duration {
		p.SecondsIntoFrame -= p.Frame().Duration
		p.FrameIndex++
		if p.FrameIndex == len(p.Action.Frames) {
			p.FrameIndex = 0
			oldAction := p.Action
			p.Action = p.NextAction
			// if non repeating, queue the current action up again.
			// usually used for defaulting back to "Idle"
			if !p.Repeat { p.NextAction = oldAction }
		}
	}
}

