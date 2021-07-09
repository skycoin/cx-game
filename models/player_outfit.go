package models

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/spriteloader"
)

var (
	helmSpriteSheetId spriteloader.SpritesheetID
	suitSpriteSheetId spriteloader.SpritesheetID
)

const (

	//1 based
	DEFAULT_HELM = 24

	DEFAULT_SUIT = 5

	HELMETS_COUNT = 30
	SUIT_COUNT    = 5
)

//
func (player *Player) DrawOutfit(position mgl32.Vec2) {

	//draw suit
	var dirX int = 0
	if player.XDirection < 0 {
		dirX = 1
	}
	for i, suitSprite := range player.suitSpriteIds {

		spriteloader.DrawSpriteQuad(
			position.X()+float32((i+dirX)%2)-(0.5), position.Y()-float32(i/2)+0.25,
			// player sprite actually faces left so throw an extra (-) here
			-player.Size.X*player.XDirection/2, player.Size.Y/2, suitSprite,
		)
	}

	//draw helmet
	for i, helmSprite := range player.helmSpriteIds {
		spriteloader.DrawSpriteQuad(
			position.X()+float32((i+dirX)%2)-0.5, position.Y()-float32(i/2)+1.25,
			// player sprite actually faces left so throw an extra (-) here
			-player.Size.X*player.XDirection/2, player.Size.Y/2, helmSprite,
		)
	}
}

func loadOutfits() {
	//load all outfits into internal spritesheet
	helmSpriteSheetId = spriteloader.LoadSpriteSheet("./assets/character/character-helmets.png")
	suitSpriteSheetId = spriteloader.LoadSpriteSheet("./assets/character/character-suits.png")

	//load helmets
	for i := 0; i < HELMETS_COUNT; i++ {
		x := i % 10 * 2
		y := i/10*2 + i/10
		spriteloader.LoadSprite(helmSpriteSheetId, fmt.Sprintf("helmet-%d-1", i+1), x+1, y)
		spriteloader.LoadSprite(helmSpriteSheetId, fmt.Sprintf("helmet-%d-2", i+1), x, y)
		spriteloader.LoadSprite(helmSpriteSheetId, fmt.Sprintf("helmet-%d-3", i+1), x+1, y+1)
		spriteloader.LoadSprite(helmSpriteSheetId, fmt.Sprintf("helmet-%d-4", i+1), x, y+1)
	}

	//load suits
	for i := 0; i < SUIT_COUNT; i++ {
		x := i * 2
		spriteloader.LoadSprite(suitSpriteSheetId, fmt.Sprintf("suit-%d-1", i+1), x+1, 1)
		spriteloader.LoadSprite(suitSpriteSheetId, fmt.Sprintf("suit-%d-2", i+1), x, 1)
		spriteloader.LoadSprite(suitSpriteSheetId, fmt.Sprintf("suit-%d-3", i+1), x+1, 2)
		spriteloader.LoadSprite(suitSpriteSheetId, fmt.Sprintf("suit-%d-4", i+1), x, 2)
	}
}

func (player *Player) SetHelm(id int) {

	if id <= 0 || id > HELMETS_COUNT {
		return
	}
	player.helmId = id
	for i := 0; i < 4; i++ {
		player.helmSpriteIds[i] = spriteloader.GetSpriteIdByName(fmt.Sprintf("helmet-%d-%d", id, i+1))
	}
}

func (player *Player) SetHelmPrev() {
	newHelmId := player.helmId%HELMETS_COUNT - 1
	if newHelmId < 0 {
		newHelmId += HELMETS_COUNT
	}
	player.SetHelm(newHelmId)
}

func (player *Player) SetHelmNext() {
	player.SetHelm(player.helmId%HELMETS_COUNT + 1)
}

func (player *Player) SetSuit(id int) {

	if id <= 0 || id > SUIT_COUNT {
		return
	}
	player.suitId = id
	for i := 0; i < 4; i++ {
		player.suitSpriteIds[i] = spriteloader.GetSpriteIdByName(fmt.Sprintf("suit-%d-%d", id, i+1))
	}
}

func (player *Player) SetSuitPrev() {
	newSuitId := player.helmId%SUIT_COUNT - 1
	if newSuitId < 0 {
		newSuitId += SUIT_COUNT
	}
	player.SetHelm(newSuitId)
}
func (player *Player) SetSuitNext() {
	player.SetSuit(player.suitId%5 + 1)
}
