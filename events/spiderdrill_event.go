package events

var OnSpiderJump onSpiderJump

type SpiderEventData struct {
	// agent *agents.Agent
	WaitingFor float32
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
