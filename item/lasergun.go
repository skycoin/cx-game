package item

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

func UseLaserGun(info ItemUseInfo) {
	log.Print("shooting laser gun")
	// click relative to camera
	camCoords :=
		mgl32.Vec4{
			info.ScreenX / render.PixelsPerTile,
			info.ScreenY / render.PixelsPerTile, 0, 1 }
	// click relative to world
	worldCoords := info.Camera.GetTransform().Mul4x1(camCoords)
	positions := cxmath.Raytrace(
		float64(info.Player.Pos.X),float64(info.Player.Pos.Y),
		float64(worldCoords.X()),float64(worldCoords.Y()))

	for _,pos := range positions {
		tile := info.Planet.GetTile(int(pos.X),int(pos.Y),world.TopLayer)
		if tile.TileType != world.TileTypeNone {
			info.Planet.DamageTile(int(pos.X), int(pos.Y), world.TopLayer)
			return
		}
	}
}

func RegisterLaserGunItemType() ItemTypeID {
	// TODO use "lasergun" instead of "redblip" once we have asset
	laserGunItemType := NewItemType(spriteloader.GetSpriteIdByName("RedBlip"))
	laserGunItemType.Use = UseLaserGun
	laserGunItemType.Name = "Laser Gun"
	return AddItemType(laserGunItemType)
}
