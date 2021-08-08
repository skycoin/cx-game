package events

var OnSpiderJump onSpiderJump
var OnSpiderBeforeJump onSpiderBeforeJump

type SpiderEventData struct {
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
