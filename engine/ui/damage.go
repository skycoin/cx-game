package ui

import (
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const (
	maxFontScale float32 = 0.5
	minFontScale float32 = 0.3
	fontShrinkFactor float32 = 10
	damageIndicatorLifetime float32 = 1
)

type DamageIndicator struct {
	Damage int
	TimeToLive float32
	Position mgl32.Vec2
}
var damageIndicators = []DamageIndicator{}

func CreateDamageIndicator( damage int, position mgl32.Vec2 ) {
	damageIndicators = append ( damageIndicators, DamageIndicator {
		damage, damageIndicatorLifetime, position } )
}

func DrawDamageIndicators(invCameraTransform mgl32.Mat4) {
	for _, damageIndicator := range damageIndicators {
		damageIndicator.Draw(invCameraTransform)
	}
}

func TickDamageIndicators(dt float32) {
	newDamageIndicators := make([]DamageIndicator, 0, len(damageIndicators))
	for _, damageIndicator := range damageIndicators {
		damageIndicator.TimeToLive -= dt
		if damageIndicator.TimeToLive > 0 {
			newDamageIndicators = append(newDamageIndicators, damageIndicator)
		}
	}
	damageIndicators = newDamageIndicators
}

func (d DamageIndicator) Alpha() float32 {
	return d.TimeToLive / damageIndicatorLifetime
}

func (d DamageIndicator) FontScale() float32 {
	x := 1 - d.TimeToLive / damageIndicatorLifetime // normalized time to live
	expFactor := maxFontScale - minFontScale
	expValue := expFactor * math32.Exp(-fontShrinkFactor*x)

	return expValue + minFontScale
}

func (d DamageIndicator) Draw(invCameraTransform mgl32.Mat4) {
	text := strconv.Itoa(d.Damage)

	homogPos := d.Position.Vec4(0,1)
	camPos := invCameraTransform.Mul4x1(homogPos)
	windowPos := render.Projection.Mul4x1(camPos)
	// this is hacky
	x := 640/2 + windowPos.X()*320
	y := 480/2 - windowPos.Y()*240
	// copy font such that we can safely mutate the color
	fontCopy := *DefaultFont
	fontCopy.SetColor(1,0,0,d.Alpha())
	fontCopy.Printf(x,y, d.FontScale(), text)
}
