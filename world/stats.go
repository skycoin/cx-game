package world

import (
	"log"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type WorldEventType int

type WorldEvent struct {
	Type WorldEventType
	ItemTypeID types.ItemTypeID
	AgentTypeID constants.AgentTypeID
	Tick int
}

// TODO pass in current world tick
func NewMobKilledEvent(id constants.AgentTypeID, tick int) WorldEvent {
	return WorldEvent {
		AgentTypeID: id, Type: EVENT_TYPE_MOB_KILLED,
	}
}

const (
	EVENT_TYPE_UNDEFINED WorldEventType = iota
	EVENT_TYPE_ITEM_CREATED
	EVENT_TYPE_ITEM_CONSUMED
	EVENT_TYPE_MOB_KILLED
)

type WorldStats struct {
	ItemsCreated map[types.ItemTypeID]int
	ItemsConsumed map[types.ItemTypeID]int
	MobsKilled map[constants.AgentTypeID]int
	EventLog []WorldEvent
}

func NewWorldStats() WorldStats {
	return WorldStats {
		ItemsCreated: make(map[types.ItemTypeID]int),
		ItemsConsumed: make(map[types.ItemTypeID]int),
		MobsKilled: make(map[constants.AgentTypeID]int),
	}
}

func (stats *WorldStats) Log(ev WorldEvent) {
	stats.EventLog = append(stats.EventLog, ev)
	if ev.Type == EVENT_TYPE_ITEM_CREATED {
		stats.ItemsCreated[ev.ItemTypeID]++
	}
	if ev.Type == EVENT_TYPE_ITEM_CONSUMED {
		stats.ItemsConsumed[ev.ItemTypeID]++
	}
	if ev.Type == EVENT_TYPE_MOB_KILLED {
		stats.MobsKilled[ev.AgentTypeID]++
	}
}

func (stats *WorldStats) Print() {
	log.Printf("%+v", stats)
}
