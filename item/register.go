package item

var (
	LaserGunItemTypeID        ItemTypeID
	GunItemTypeID             ItemTypeID
	RockDustItemTypeID        ItemTypeID
	FurnitureToolItemTypeID   ItemTypeID
	TileToolItemTypeID        ItemTypeID
	BgToolItemTypeID          ItemTypeID
	EnemyToolItemTypeID       ItemTypeID
	PipePlaceToolItemTypeID   ItemTypeID
	PipeConnectToolItemTypeID ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
	RockDustItemTypeID = RegisterRockDustItemType()
	FurnitureToolItemTypeID = RegisterFurnitureToolItemType()
	TileToolItemTypeID = RegisterTileToolItemType()
	BgToolItemTypeID = RegisterBgToolItemType()
	PipePlaceToolItemTypeID = RegisterPipeToolItemType()
	PipeConnectToolItemTypeID = RegisterPipeConnectToolItemType()
	EnemyToolItemTypeID = RegisterEnemyToolItemType()

	AddDrops()
}
