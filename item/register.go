package item

var (
	LaserGunItemTypeID      ItemTypeID
	GunItemTypeID           ItemTypeID
	RockDustItemTypeID      ItemTypeID
	FurnitureToolItemTypeID ItemTypeID
	TileToolItemTypeID      ItemTypeID
	EnemyToolItemTypeID     ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
	RockDustItemTypeID = RegisterRockDustItemType()
	FurnitureToolItemTypeID = RegisterFurnitureToolItemType()
	TileToolItemTypeID = RegisterTileToolItemType()
	EnemyToolItemTypeID = RegisterEnemyToolItemType()

	AddDrops()
}
