package item

import (
	"github.com/skycoin/cx-game/world"
)

func CanPlaceTileTypeAt(World *world.World, tt *world.TileType, x,y int) bool {
	//midlayer and toplayer can't occupy same tile
	layersToCheck := layersToCheckForPlace(tt.Layer)
	occupyingTilesAreClear :=
		tilesAreClear(World,layersToCheck,x,y,x+int(tt.Width),y+int(tt.Height))

	wouldBeSupported := tileTypeWouldBeSupported(World,tt,x,y)

	return occupyingTilesAreClear && wouldBeSupported
}

func tileTypeWouldBeSupported(
	World *world.World,tt *world.TileType, x,y int,
) bool {
	if tt.NeedsGround {
		belowTilesAreSolid := tilesAreSolid(
			World, []world.LayerID { world.TopLayer },
			x,y-1,x+int(tt.Width),y,
		)
		if !belowTilesAreSolid { return false }
	}
	if tt.NeedsRoof {
		aboveTilesAreSolid := tilesAreSolid(
			World, []world.LayerID { world.TopLayer },
			x,y+int(tt.Height),x+int(tt.Width),y+int(tt.Height)+1,
		)
		if !aboveTilesAreSolid { return false }
	}
	return true
}
