package enemies

import "github.com/skycoin/cx-game/spriteloader"

type Slime struct {
	SlimeAnimated spriteloader.SpriteAnimated
}

func (slime *Slime) NewSlime() {
	// slime.SlimeAnimated = spriteloader.NewSpriteAnimated("./assets/slime.json")
}
