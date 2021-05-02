package world

const TileTypeNone = 0
const TileTypeNormal = 1

type Tile struct {
	SpriteID uint32
	TileType uint32
    Name string
}
