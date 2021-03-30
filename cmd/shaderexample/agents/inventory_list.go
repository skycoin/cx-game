package agents

type Inventory struct {
	InventoryID    int
	InventorySize  int
	InventoryXSize int
	InventoryYSize int
	ItemSlot       []InventorySlot
}

type InventorySlot struct {
	ItemID       int
	ItemQuantity int
}
