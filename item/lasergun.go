package item

import (
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/particles"
)

func UseLaserGun(info ItemUseInfo) {
	worldCoords := info.WorldCoords()
	// adding 0.5 here because raytrace assumes top left coords,
	// but player and tile positions are stored using centered coords
	positions := cxmath.Raytrace(
		float64(info.Player.Pos.X)+0.5,float64(info.Player.Pos.Y)+0.5,
		float64(worldCoords.X())+0.5,float64(worldCoords.Y()) + 0.5 )

	playerPos := info.PlayerCoords()
	for _,pos := range positions {
		if info.Planet.TileIsSolid(int(pos.X),int(pos.Y)) {
			direction := worldCoords.Sub(playerPos).Normalize()
			length := pos.Vec2().Sub(playerPos).Len() + 0.5
			targetPos := playerPos.Add(direction.Mul(length))

			closePlayerPos, closeTargetPos :=
				info.Planet.MinimizeDistance(playerPos, targetPos)
			particles.CreateLaser(closePlayerPos, closeTargetPos)

			tile,destroyed :=
				info.Planet.DamageTile(int(pos.X), int(pos.Y), world.TopLayer)

			if destroyed {
				itemTypeId := GetItemTypeIdForTile(tile)
				CreateWorldItem(itemTypeId, pos.Vec2())
			}

			return
		}
	}

	// hit nothing - visual effect only
	//particles.CreateLaser(playerPos, worldCoords )
	closePlayerPos, closeTargetPos :=
		info.Planet.MinimizeDistance(playerPos, worldCoords)
	particles.CreateLaser(closePlayerPos, closeTargetPos)
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
