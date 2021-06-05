package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/particles"
)

func UseLaserGun(info ItemUseInfo) {
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

	particles.CreateLaser(
		mgl32.Vec2{ info.Player.Pos.X, info.Player.Pos.Y},
		worldCoords.Vec2(),
	)

	for _,pos := range positions {
		tile := info.Planet.GetTile(int(pos.X),int(pos.Y),world.TopLayer)
		if tile.TileType != world.TileTypeNone {
			info.Planet.DamageTile(int(pos.X), int(pos.Y), world.TopLayer)
			return
		}
	}
}

func RegisterLaserGunItemType() ItemTypeID {
	// TODO use proper asset
	laserGunSpriteId :=spriteloader.LoadSingleSprite(
			"./assets/item/lasergun-temp.png","lasergun")
	laserGunItemType := NewItemType(laserGunSpriteId)
	laserGunItemType.Use = UseLaserGun
	laserGunItemType.Name = "Laser Gun"
	return AddItemType(laserGunItemType)
}
