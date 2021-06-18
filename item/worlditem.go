// a world item is an item that is "floating" in the world
package item

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
)

var bloopStreamer beep.StreamSeeker

const bloopVolume = 1

func InitWorldItem() {
	// setup bloop sound
	file, err := os.Open("./assets/sound/bloop.wav")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := wav.Decode(file)
	bloopStreamer = streamer
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

const pickupRadius = 0.2
const attractRadius = 2
const worldItemSize = 0.5

type WorldItem struct {
	physics.Body
	ItemTypeId ItemTypeID
}

func NewWorldItem(ItemTypeId ItemTypeID) *WorldItem {
	item := WorldItem{
		Body: physics.Body{
			Size: physics.Vec2{X: worldItemSize, Y: worldItemSize},
		},
		ItemTypeId: ItemTypeId,
	}
	physics.RegisterBody(&item.Body)
	return &item
}

func (item WorldItem) Draw(cam *camera.Camera) {
	if !cam.IsInBoundsF(item.Pos.X, item.Pos.Y) {
		return
	}
	spriteId := itemTypes[item.ItemTypeId].SpriteID
	z := -spriteloader.SpriteRenderDistance
	view := mgl32.Translate3D(-cam.X, -cam.Y, 0)
	world := mgl32.Translate3D(
		item.Pos.X,
		item.Pos.Y,
		z,
	).Mul4(cxmath.Scale(worldItemSize))
	modelView := view.Mul4(world)
	spriteloader.DrawSpriteQuadMatrix(modelView, spriteId)
}

func (item *WorldItem) Tick(
	planet *world.Planet, dt float32,
	playerPos physics.Vec2,
) bool {
	item.Vel.Y -= physics.Gravity * dt

	itemToPlayerDisplacement := playerPos.Sub(item.Pos)
	itemToPlayerDistSqr := itemToPlayerDisplacement.LengthSqr()
	if itemToPlayerDistSqr < attractRadius*attractRadius {
		attractForce := itemToPlayerDisplacement.Mult(1 / itemToPlayerDisplacement.LengthSqr())
		item.Vel = item.Vel.Add(attractForce)
	}

	//item.Move(planet, dt)
	didPickup := itemToPlayerDistSqr < pickupRadius*pickupRadius
	if didPickup {
		speaker.Play(bloopStreamer)
	}
	return didPickup
}
