package item

import (
	"log"

	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/components/types"
)

func NewDevInventory() types.InventoryID {
	inventoryId := NewInventory(10, 8)
	inventory := GetInventoryById(inventoryId)
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 0)] =
		InventorySlot{LaserGunItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(2, 0)] =
		InventorySlot{GunItemTypeID, 1, 0}

	pipeTileType, ok := world.GetTileTypeByID(world.IDFor("pipe"))
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

	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 1)] = InventorySlot{
		GetItemTypeIdForTileTypeID(world.IDFor("stone")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 2)] = InventorySlot{
		GetItemTypeIdForTileTypeID(world.IDFor("regolith")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 3)] = InventorySlot{
		GetItemTypeIdForTileTypeID(world.IDFor("regolith-wall")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 4)] = InventorySlot{
		GetItemTypeIdForTileTypeID(world.IDFor("bedrock")),
		1, 0,
	}
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 5)] = InventorySlot{
		GetItemTypeIdForTileTypeID(world.IDFor("water")),
		1, 0,
	}
	return inventoryId
}
