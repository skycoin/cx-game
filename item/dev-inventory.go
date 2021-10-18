package item

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/world"
)

func idFor(name string) world.TileTypeID {
	id, ok := world.IDFor(name)
	if !ok {
		log.Fatalf("cannot find tile type ID for \"%s\"", name)
	}
	return id
}

func NewDevInventory() types.InventoryID {
	inventoryId := NewInventory(10, 8)
	inventory := GetInventoryById(inventoryId)
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 0)] =
		InventorySlot{LaserGunItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(2, 0)] =
		InventorySlot{GunItemTypeID, 1, 0}

	pipeTileType, ok := world.GetTileTypeByID(idFor("pipe"))
	if !ok {
		log.Fatal("Cannot find pipe tile type")
	}
	pipeTile := pipeTileType.CreateTile(world.TileCreationOptions{})
	pipeItemTypeID := GetItemTypeIdForTile(pipeTile)
	inventory.Slots[inventory.ItemSlotIndexForPosition(3, 0)] =
		InventorySlot{pipeItemTypeID, 20, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(4, 0)] =
		InventorySlot{FurnitureToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(5, 0)] =
		InventorySlot{TileToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(6, 0)] =
		InventorySlot{EnemyToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(7, 0)] =
		InventorySlot{PipePlaceToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(8, 0)] =
		InventorySlot{PipeConnectToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(9, 0)] =
		InventorySlot{BgToolItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 0)] =
		InventorySlot{DevDestroyToolID, 1, 0}

	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 1)] = InventorySlot{
		GetItemTypeIdForTileTypeID(idFor("stone")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 2)] = InventorySlot{
		GetItemTypeIdForTileTypeID(idFor("regolith")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 3)] = InventorySlot{
		GetItemTypeIdForTileTypeID(idFor("regolith-wall")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 4)] = InventorySlot{
		GetItemTypeIdForTileTypeID(idFor("bedrock")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 5)] = InventorySlot{
		GetItemTypeIdForTileTypeID(idFor("water")),
		1, 0,
	}
	return inventoryId
}
