package agents

import "github.com/skycoin/cx-game/render"

//todo structure better
type SpineData struct {
	Back_armSpriteID       render.SpriteID
	Back_footSpriteID      render.SpriteID
	Back_handSpriteID      render.SpriteID
	Back_legSpriteID       render.SpriteID
	BodySpriteID           render.SpriteID
	Front_armSpriteID      render.SpriteID
	Front_footSpriteID     render.SpriteID
	Front_handSpriteID     render.SpriteID
	Front_legtSpriteID     render.SpriteID
	HeadSpriteID           render.SpriteID
	IgnoringPlatformsFor_2 float32
}
