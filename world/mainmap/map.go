package mainmap

import (
	"fmt"
	"math/rand"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

func NewMap() *Map {

	newMap := Map{
		bounds: Fulstrum{
			Left:   0,
			Right:  shownSize,
			Top:    shownSize,
			Bottom: 0,
		},
	}

	return &newMap
}

//isInBounds checks if tile in the fullstrum
func (m *Map) isInBounds(tile *MapTile) bool {
	if tile.x > m.bounds.Left &&
		tile.x < m.bounds.Right &&
		tile.y < m.bounds.Top &&
		tile.y > m.bounds.Bottom {
		return true
	}
	return false
}

//InitMap inits map with some data
func InitMap(window *render.Window) {

	m = NewMap()

	spriteloader.InitSpriteloader(window)

	spriteSheetId := spriteloader.LoadSpriteSheet("./assets/8x8/test-tile-stone-02.png")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			spriteloader.LoadSprite(spriteSheetId, fmt.Sprintf("tile%d", i*8+j), j, i)
		}
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			show := 1
			if rand.Float32() < 0.15 {
				show = 0
			}
			spriteId := spriteloader.GetSpriteIdByName(fmt.Sprintf("tile%d", (size*y)%8*8+x%8))

			m.tiles[x][y] = &MapTile{
				tileIdBackground: spriteId,
				tileIdMid:        spriteId,
				tileIdFront:      spriteId,
				x:                x,
				y:                y,
				show:             show,
			}
		}
	}
}

//CreateMapTile creates map tile

//DrawMap draws map with data provided. It somehow draws from bottom right to top left need to figuree out why
func DrawMap() {

	var firstTile *MapTile
	for _, row := range m.tiles {
		for _, tile := range row {
			if !m.isInBounds(tile) {
				continue
			}
			if firstTile == nil {
				firstTile = tile
			}

			xpos, ypos := convertCoordinates(tile.x, tile.y, firstTile)
			if tile.show == 0 {
				continue
			}

			spriteloader.DrawSpriteQuad(xpos, ypos, 1, 1, tile.tileIdMid)
		}
	}

}

func convertCoordinates(x, y int, tile *MapTile) (int, int) {
	return x - tile.x, y - tile.y
}

//GoTop moves map one row up
func GoTop() {
	if m.bounds.Top == size {
		return
	}
	m.bounds.Top += 1
	m.bounds.Bottom += 1
}

//GoBottom moves map one row down
func GoBottom() {
	if m.bounds.Bottom == -1 {
		return
	}
	m.bounds.Top -= 1
	m.bounds.Bottom -= 1
}

//GoLeft moves map one column left
func GoLeft() {
	if m.bounds.Left == size-shownSize {
		return
	}
	m.bounds.Left += 1
	m.bounds.Right += 1
}

//GoRight moves map one column right
func GoRight() {
	if m.bounds.Right == shownSize-1 {
		return
	}
	m.bounds.Left -= 1
	m.bounds.Right -= 1
}
