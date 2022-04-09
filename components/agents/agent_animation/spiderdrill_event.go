package events

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
)

var OnSpiderJump onSpiderJump
var OnSpiderBeforeJump onSpiderBeforeJump
var OnSpiderCollisionHorizontal onSpiderCollisionHorizontal

type SpiderEventData struct {
	Agent             *agents.Agent
	AnimationPlayback anim.Playback
}

type onSpiderCollisionHorizontal struct {
	handlers []interface{ OnSpiderCollisionHorizontal(SpiderEventData) }
}

func (o *onSpiderCollisionHorizontal) Register(handler interface{ OnSpiderCollisionHorizontal(SpiderEventData) }) {
	o.handlers = append(o.handlers, handler)
}

func (o onSpiderCollisionHorizontal) Trigger(data SpiderEventData) {
	for _, handler := range o.handlers {
		go handler.OnSpiderCollisionHorizontal(data)
	}
}

type onSpiderJump struct {
	handlers []interface{ OnSpiderDrillJump(SpiderEventData) }
}

func (o *onSpiderJump) Register(handler interface{ OnSpiderDrillJump(SpiderEventData) }) {
	o.handlers = append(o.handlers, handler)
}

func (o onSpiderJump) Trigger(data SpiderEventData) {
	for _, handler := range o.handlers {
		go handler.OnSpiderDrillJump(data)
	}
}

type onSpiderBeforeJump struct {
	handlers []interface{ OnSpiderDrillBeforeJump(SpiderEventData) }
}

func (o *onSpiderBeforeJump) Register(handler interface{ OnSpiderDrillBeforeJump(SpiderEventData) }) {
	o.handlers = append(o.handlers, handler)
}

func (o onSpiderBeforeJump) Trigger(data SpiderEventData) {
	for _, handler := range o.handlers {
		go handler.OnSpiderDrillBeforeJump(data)
	}
}
