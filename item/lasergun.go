package item

import (
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/world"
)

func UseLaserGun(info ItemUseInfo) {
	worldCoords := info.WorldCoords()
	// adding 0.5 here because raytrace assumes top left coords,
	// but player and tile positions are stored using centered coords

	playerPos := info.PlayerCoords()
	positions := cxmath.Raytrace(
		float64(playerPos.X())+0.5, float64(playerPos.Y())+0.5,
		float64(worldCoords.X())+0.5, float64(worldCoords.Y())+0.5)
	for _, pos := range positions {
		if info.World.Planet.TileIsSolid(int(pos.X), int(pos.Y)) {
			direction := worldCoords.Sub(playerPos).Normalize()
			length := pos.Vec2().Sub(playerPos).Len() + 0.5
			targetPos := playerPos.Add(direction.Mul(length))

			closePlayerPos, closeTargetPos :=
				info.World.Planet.MinimizeDistance(playerPos, targetPos)
			particles.CreateLaser(closePlayerPos, closeTargetPos)

			tile, destroyed :=
				info.World.Planet.DamageTile(int(pos.X), int(pos.Y), world.TopLayer)
			_ = tile

			if destroyed {
				//itemTypeId := GetItemTypeIdForTile(tile)
				//CreateWorldItems(tile.TileTypeID, pos.Vec2())
				//CreateWorldItem(itemTypeId, pos.Vec2())
			}

			return
		}
	}

	// hit nothing - visual effect only
	//particles.CreateLaser(playerPos, worldCoords )
	closePlayerPos, closeTargetPos :=
		info.World.Planet.MinimizeDistance(playerPos, worldCoords)
	particles.CreateLaser(closePlayerPos, closeTargetPos)
}

func RegisterLaserGunItemType() ItemTypeID {
	laserGunItemType := NewItemType(render.GetSpriteIDByName("lasergun"))
	laserGunItemType.Use = UseLaserGun
	laserGunItemType.Name = "Laser Gun"
	return AddItemType(laserGunItemType)
}
