package mainmap

import (
	"math/rand"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

func NewMap() *Map {

	newMap := Map{
		bounds: Fullstrum{
			Left:   0,
			Right:  shownSize,
			Top:    shownSize,
			Bottom: 0,
		},
	}

	return &newMap
}

func GoTop() {
	m.bounds.Top += 1
	m.bounds.Bottom += 1
}
func GoBottom() {
	m.bounds.Top -= 1
	m.bounds.Bottom -= 1
}
func GoLeft() {
	m.bounds.Left += 1
	m.bounds.Right += 1
}
func GoRight() {
	m.bounds.Left -= 1
	m.bounds.Right -= 1
}

func (m *Map) isInBounds(tile *MapTile) bool {
	if tile.x > m.bounds.Left &&
		tile.x < m.bounds.Right &&
		tile.y < m.bounds.Top &&
		tile.y > m.bounds.Bottom {
		return true
	}
	return false
}

func InitMap(window *render.Window) {

	m := NewMap()

	spriteloader.InitSpriteloader(window)
	spriteSheetId := spriteloader.LoadSpriteSheet("./second_assets/32x32-test-tiles-01.png")
	spriteloader.LoadSprite(spriteSheetId, "blue_star", 0, 1)
	spriteloader.LoadSprite(spriteSheetId, "brown_square", 1, 1)
	spriteloader.LoadSprite(spriteSheetId, "green_star", 2, 1)

	bspriteId := spriteloader.GetSpriteIdByName("blue_star")
	// mspriteId := spriteloader.GetSpriteIdByName("brown_square")
	fSpriteId := spriteloader.GetSpriteIdByName("green_star")

	var randomBackgroundId int
	randValue := rand.Float32() > 0.5
	switch randValue {
	case true:
		randomBackgroundId = bspriteId
	default:
		randomBackgroundId = fSpriteId
	}
	for i := 0; i < size*size; i++ {
		show := 1
		if rand.Float32() < 0.15 {
			show = 0
		}
		m.tiles = append(m.tiles, &MapTile{
			spriteId:         i,
			tileIdBackground: bspriteId,
			tileIdMid:        randomBackgroundId,
			tileIdFront:      fSpriteId,
			x:                i % size,
			y:                i / size,
			show:             show,
		})
	}
}

func DrawMap() {

	tiles := make([]*MapTile, 0)
	for _, tile := range m.tiles {
		if !m.isInBounds(tile) {
			continue
		}

		tiles = append(tiles, tile)
	}

	for i, tile := range tiles {
		if tile.show == 0 {
			continue
		}
		spriteloader.DrawSpriteQuad(i%shownSize, i/size, 1, 1, tile.tileIdMid)
	}

}
