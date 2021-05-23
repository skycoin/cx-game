package stats

import (
	"time"
	"os"
	"path"
	"log"
	"encoding/json"
	"github.com/skycoin/cx-game/item"
)

type WorldID string

type ItemStats struct {
	Crafted uint
	Produced uint
	Consumed uint
}
func NewItemStats() *ItemStats {
	return &ItemStats {}
}

type WorldStats struct {
	Item map[item.ItemTypeID]*ItemStats
	Frames uint64
	StartTime time.Time
}
func NewWorldStats() *WorldStats {
	return &WorldStats {
		Item: make(map[item.ItemTypeID]*ItemStats),
	}
}

type GlobalStats struct {
	World map[WorldID]*WorldStats
}
func NewGlobalStats() *GlobalStats {
	return &GlobalStats {
		World: make(map[WorldID]*WorldStats),
	}
}

func (stats *GlobalStats) AddWorld(id WorldID) *WorldStats {
	stats.World[id] = NewWorldStats()
	stats.World[id].StartTime = time.Now()
	return stats.World[id]
}

func (stats *WorldStats) LogFrame() {
	stats.Frames += 1
}

func (stats *WorldStats) GetItemStats(itemTypeId item.ItemTypeID) *ItemStats {
	_,ok := stats.Item[itemTypeId]
	if !ok {
		stats.Item[itemTypeId] = NewItemStats()
	}
	return stats.Item[itemTypeId]
}

var statsFilename string
func init() {
	home, err := os.UserHomeDir()
	if err !=nil {
		log.Fatal(err)
	}
	appPath := path.Join(home,".cx-game")
	err = os.MkdirAll(appPath,0700)
	if err!=nil {
		log.Fatal(err)
	}
	statsFilename = path.Join(appPath, "stats.dat")
}

func (stats *GlobalStats) Save() {
	buf,err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(statsFilename, buf, 0600)
}

func LoadGlobalStats() GlobalStats {
	var stats GlobalStats

	buf,err := os.ReadFile(statsFilename)
	if err!=nil {
		log.Fatal(err)
	}
	json.Unmarshal(buf,&stats)
	return stats
}
