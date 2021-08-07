package events

import "github.com/skycoin/cx-game/cxmath"

var OnSpiderJump onSpiderJump

// type SpiderFacingData struct {
// }

type SpiderJumpData struct {
	Vel cxmath.Vec2
}

// type OnSpiderFacingLeft struct {
// 	handlers []interface{ Handle(SpiderFacingData) }
// }

type onSpiderJump struct {
	handlers []interface{ Handle(SpiderJumpData) }
}

func (o *onSpiderJump) Register(handler interface{ Handle(SpiderJumpData) }) {
	o.handlers = append(o.handlers, handler)
}

func (o onSpiderJump) Trigger(data SpiderJumpData) {
	for _, handler := range o.handlers {
		go handler.Handle(data)
	}
}

// func (o *OnSpiderFacingLeft) Register(handler interface{ Handle(SpiderFacingData) }) {
// 	o.handlers = append(o.handlers, handler)
// }

// func (o OnSpiderFacingLeft) Trigger(data SpiderFacingData) {
// 	for _, handler := range o.handlers {
// 		go handler.Handle(data)
// 	}
// }
