package ui

import (
	"fmt"
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/ui/glfont"
	"github.com/skycoin/cx-game/render"
)

const healthBarDividerWidth = float32(0.1)

type HealthBar struct {
	filledDivider, unfilledDivider spriteloader.SpriteID
	border, fill                   StretchingNineSlice
	Width, Height                  float32
	Scale                          float32
	HpPerTick                      int
	Font                           *glfont.Font
}

var DefaultFont *glfont.Font

func NewHealthBar() HealthBar {
	font, err := glfont.LoadFont("./assets/font/GravityBold8.ttf", 24, 640, 480)
	if err != nil {
		log.Fatal(err)
	}
	DefaultFont = font
	return HealthBar{
		filledDivider: spriteloader.LoadSingleSprite(
			"./assets/hud/hud_hp_bar_div1.png", "hp-bar-vertical-divider"),
		unfilledDivider: spriteloader.LoadSingleSprite(
			"./assets/hud/hud_hp_bar_div2.png", "hp-bar-horizontal-divider"),
		border: NewStretchingNineSlice(
			spriteloader.LoadSingleSprite(
				"./assets/hud/hud_hp_bar_border.png", "hp-bar-border"),
			NineSliceDims{1.0 / 6.0, 1.0 / 6.0, 1.0 / 8.0, 2.0 / 8.0},
		),
		fill: NewStretchingNineSlice(
			spriteloader.LoadSingleSprite(
				"./assets/hud/hud_hp_bar_fill.png", "hp-bar-fill"),
			NineSliceDims{1.0 / 5.0, 1.0 / 5.0, 1.0 / 5.0, 1.0 / 5.0},
		),
		Width: 8, Height: 1,
		Scale: 0.5, HpPerTick: 25,
		Font: font,
	}
}

func (bar HealthBar) Draw(ctx render.Context, hp, maxHP int) {
	topLeftCtx := ctx.
		PushLocal(mgl32.Translate3D(-0.5, 0.5, 0)).
		PushLocal(cxmath.Scale(bar.Scale))
	hpFrac := float32(hp) / float32(maxHP)
	bar.fill.Draw(topLeftCtx, mgl32.Vec2{hpFrac * bar.Width, bar.Height})
	for tick := 0; tick < maxHP; tick += bar.HpPerTick {
		var divider spriteloader.SpriteID
		if hp > tick {
			divider = bar.filledDivider
		} else {
			divider = bar.unfilledDivider
		}
		x := bar.Width * float32(tick) / float32(maxHP)
		spriteloader.DrawSpriteQuadContext(
			topLeftCtx.
				PushLocal(mgl32.Translate3D(x, -0.5, 0)).
				PushLocal(mgl32.Scale3D(healthBarDividerWidth, 1, 1)),
			divider,
			spriteloader.NewDrawOptions(),
		)
	}
	bar.border.Draw(topLeftCtx, mgl32.Vec2{bar.Width, bar.Height})
	text := fmt.Sprintf("%d/%d", hp, maxHP)
	bar.Font.SetColor(1, 1, 1, 1)
	fs := float32(0.3) // font scale
	x := 130 - bar.Font.Width(fs, text)
	bar.Font.Printf(x, 37, fs, text)
	//utility.DrawColorQuad(ctx, mgl32.Vec4{1,0,0,1})
	/*
		spriteloader.DrawSpriteQuadContext(
			ctx.PushLocal(mgl32.Scale3D(0.1,1,1)), bar.filledDivider)
	*/
	// TODO
}

type CircleIndicator struct {
	spriteID spriteloader.SpriteID
}

func NewCircleIndicator(spriteID spriteloader.SpriteID) CircleIndicator {
	return CircleIndicator{spriteID: spriteID}
}

// x describes how full circle is
func (indicator CircleIndicator) Draw(ctx render.Context, x float32) {
	DrawArc(ctx.MVP(), x)
	spriteloader.DrawSpriteQuadContext(
		ctx, indicator.spriteID, spriteloader.NewDrawOptions())
}

// all values are normalized to [1,1] range
type HUDState struct {
	Health    int
	MaxHealth int

	Fullness  float32 // opposite of hunger
	Hydration float32
	Oxygen    float32
	Fuel      float32
}

type HUD struct {
	Health HealthBar

	Fullness  CircleIndicator
	Hydration CircleIndicator
	Oxygen    CircleIndicator
	Fuel      CircleIndicator

	hpIconSpriteID spriteloader.SpriteID
}

var hud HUD

func initHUD() {
	hud = HUD{
		Health: NewHealthBar(),

		Fullness: NewCircleIndicator(spriteloader.LoadSingleSprite(
			"./assets/hud/hud_status_food.png", "status_food")),
		Hydration: NewCircleIndicator(spriteloader.LoadSingleSprite(
			"./assets/hud/hud_status_water.png", "status_water")),
		Oxygen: NewCircleIndicator(spriteloader.LoadSingleSprite(
			"./assets/hud/hud_status_oxygen.png", "status_oxygen")),
		Fuel: NewCircleIndicator(spriteloader.LoadSingleSprite(
			"./assets/hud/hud_status_fuel.png", "status_fuel")),

		hpIconSpriteID: spriteloader.LoadSingleSprite(
			"./assets/hud/hud_hp_icon.png", "hp_icon"),
	}
}

func DrawHUD(state HUDState) {
	hud.Draw(state)
}

const hudPadding = 1
const circleYOffset = float32(-1.2)
const circlePadding = 1.2

func (h HUD) Draw(state HUDState) {
	y := circleYOffset
	ctx := render.CenterToTopLeft(spriteloader.Window.DefaultRenderContext()).
		// padding
		PushLocal(mgl32.Translate3D(hudPadding, -hudPadding, 0))

	spriteloader.DrawSpriteQuadContext(
		ctx, h.hpIconSpriteID, spriteloader.NewDrawOptions())

	h.Health.Draw(
		ctx.PushLocal(
			mgl32.Translate3D(1, 0, 0).
				Mul4(cxmath.Scale(0.9)),
		), state.Health, state.MaxHealth)
	// TODO offset these
	h.Fullness.Draw(
		ctx.PushLocal(mgl32.Translate3D(0, y, 0)),
		state.Fullness,
	)
	h.Hydration.Draw(
		ctx.PushLocal(mgl32.Translate3D(1*circlePadding, y, 0)), state.Hydration)
	h.Oxygen.Draw(
		ctx.PushLocal(mgl32.Translate3D(2*circlePadding, y, 0)), state.Oxygen)
	h.Fuel.Draw(
		ctx.PushLocal(mgl32.Translate3D(3*circlePadding, y, 0)), state.Fuel)
}

func DrawAgentHUD(agent *agents.Agent) {
	DrawHUD(HUDState{
		Health:    agent.HealthComponent.Current,
		MaxHealth: agent.HealthComponent.Max,

		Fullness: 1, Hydration: 1, Oxygen: 1, Fuel: 1,
	})
}
