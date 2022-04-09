package events

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
)

var OnSlimeJump onSlimeJump

type SlimeEventData struct {
	Agent             *agents.Agent
	AnimationPlayback anim.Playback
}

type onSlimeJump struct {
	handlers []interface{ OnSpiderDrillJump(SlimeEventData) }
}
