package mainmap

import (
	"fmt"
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

//checks if tile in the fullstrum
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

	spriteSheetId := spriteloader.LoadSpriteSheet("./assets/stars")
	spriteloader.LoadSprite(spriteSheetId, "blue", 0, 0)
	spriteloader.LoadSprite(spriteSheetId, "gray", 1, 0)
	spriteloader.LoadSprite(spriteSheetId, "sand", 2, 0)

	bspriteId := spriteloader.GetSpriteIdByName("blue")
	mspriteId := spriteloader.GetSpriteIdByName("gray")
	fSpriteId := spriteloader.GetSpriteIdByName("sand")

	var randomBackgroundId int
	randValue := rand.Float32() > 0.5
	switch randValue {
	case true:
		randomBackgroundId = bspriteId
	default:
		randomBackgroundId = fSpriteId
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			show := 1
			if rand.Float32() < 0.15 {
				show = 0
			}
			switch y*size + x {
			case 0:
				randomBackgroundId = bspriteId
			case 1:
				randomBackgroundId = mspriteId
			case 5:
				randomBackgroundId = bspriteId
			default:
				randomBackgroundId = fSpriteId
			}
			m.tiles[x][y] = &MapTile{
				spriteId:         y*size + x,
				tileIdBackground: randomBackgroundId,
				tileIdMid:        randomBackgroundId,
				tileIdFront:      fSpriteId,
				x:                x,
				y:                y,
				show:             show,
			}
		}
	}
}

//DrawMap draws map with data provided. It somehow draws from bottom right to top left need to figuree out why
func DrawMap() {

	tiles := make([]*MapTile, 0)
	for _, row := range m.tiles {
		for _, tile := range row {
			if !m.isInBounds(tile) {
				continue
			}

			// spriteloader.DrawSpriteQuad(tile.x, tile.y, 1, 1, tile.tileIdBackground)
			tiles = append(tiles, tile)
		}
	}
	firstTile := tiles[0]
	for i, tile := range tiles {
		if i == 0 {
			fmt.Printf("%v    %v\n", tile.x, tile.y)

		}

		xpos, ypos := convertCoordinates(tile.x, tile.y, firstTile)
		if tile.show == 0 {
			continue
		}

		spriteloader.DrawSpriteQuad(xpos, ypos, 1, 1, tile.tileIdMid)
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
