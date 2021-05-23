package stats

import (
	"github.com/skycoin/cx-game/item"
	"testing"
)

const WORLD_HASH = "exampleworldhash"
const ITEM_TYPE_ID =  item.ItemTypeID(7)

func TestLogFrame(t *testing.T) {
	globalStats := NewGlobalStats()
	worldStats := globalStats.AddWorld(WORLD_HASH)
	got := worldStats.Frames
	if got!=0 {
		t.Errorf("globalStats.Frames initially = %d; want 0",got)
	}
	worldStats.LogFrame()
	got = worldStats.Frames
	if got!=1 {
		t.Errorf("worldStats.Frames = %d after logging frame; want 1",got)
	}
	got = globalStats.World[WORLD_HASH].Frames
	if got!=1 {
		t.Errorf(
			"globalStats.World[WORLD_HASH].Frames = %d "+
			"after logging frame; want 1",got)
	}
}

func TestItemProduce(t *testing.T) {
	globalStats := NewGlobalStats()
	worldStats := globalStats.AddWorld(WORLD_HASH)
	itemStats := worldStats.GetItemStats(ITEM_TYPE_ID)
	got := itemStats.Produced
	if got!=0 {
		t.Errorf("itemStat.Produced initially = %d; want 0",got)
	}
	itemStats.Produced++
	got = globalStats.World[WORLD_HASH].Item[ITEM_TYPE_ID].Produced
	if got!=1 {
		t.Errorf("itemStat.Produced = %d after producing item; want 1",got)
	}
}

func TestSerialize(t *testing.T) {
	globalStats := NewGlobalStats()
	worldStats := globalStats.AddWorld(WORLD_HASH)
	itemStats := worldStats.GetItemStats(ITEM_TYPE_ID)
	itemStats.Produced++

	globalStats.Save()
	deserializedGlobalStats := LoadGlobalStats()

	got := deserializedGlobalStats.World[WORLD_HASH].Item[ITEM_TYPE_ID].Produced
	if got!=1 {
		t.Errorf("itemStat.Produced = %d after producing item; want 1",got)
	}

}
