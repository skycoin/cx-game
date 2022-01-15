package agents

import (
	"fmt"
	"log"

	"os"

	//"github.com/skycoin/cx-game/components/agents/agent_draw"
	//	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
)

type Agent struct {
	PlayerControlled  bool
	AgentId           types.AgentID
	Meta              AgentMeta
	Handlers          AgentHandlers
	Timers            AgentTimers
	Transform         physics.Body
	Health            HealthComponent
	AnimationPlayback anim.Playback
	InventoryID       types.InventoryID

	Time  float64
	Play  bool
	Speed float64

	Skeleton       *spine.Skeleton
	Animation      *spine.Animation
	SkinIndex      int
	AnimationIndex int

	DebugCenter bool
	DebugBones  bool
}

func (a *Agent) LoadCharacter(loc animation.Location) (*Agent, error) {
	fmt.Println("test here: ", loc.JSON)
	rd, err := os.Open(loc.JSON)
	if err != nil {
		fmt.Println("hit Error: ")
		return nil, err
	}
	fmt.Println("test data: ")
	data, err := spine.ReadJSON(rd)
	if err != nil {
		return nil, err
	}

	data.Name = loc.Name

	//char.ImagesPath = loc.Images

	//char.Images = make(map[string]*pixel.PictureData)

	//fmt.Println("____________Debug point_________")

	//fmt.Printf("%v", char.Images)

	//	fmt.Println("________________________________")

	a.Play = true
	a.DebugBones = true
	a.DebugCenter = false

	a.Speed = 1
	a.Skeleton = spine.NewSkeleton(data)
	a.Skeleton.Skin = a.Skeleton.Data.DefaultSkin
	a.Animation = a.Skeleton.Data.Animations[0]

	a.AnimationIndex = 0
	a.SkinIndex = 0

	a.Skeleton.UpdateAttachments()
	a.Skeleton.Update()
	fmt.Printf("%v", a)
	return a, nil
}

func (char *Agent) NextAnimation(offset int) {
	char.AnimationIndex += offset
	for char.AnimationIndex < 0 {
		char.AnimationIndex += len(char.Skeleton.Data.Animations)
	}
	char.AnimationIndex = char.AnimationIndex % len(char.Skeleton.Data.Animations)
	char.Animation = char.Skeleton.Data.Animations[char.AnimationIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
}

func (char *Agent) NextSkin(offset int) {
	char.SkinIndex += offset
	for char.SkinIndex < 0 {
		char.SkinIndex += len(char.Skeleton.Data.Skins)
	}
	char.SkinIndex = char.SkinIndex % len(char.Skeleton.Data.Skins)
	char.Skeleton.Skin = char.Skeleton.Data.Skins[char.SkinIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
	char.Skeleton.UpdateAttachments()
}

func (char *Agent) Update(dt float64, center cxmath.Vec2) {
	if char.Play {
		char.Time += dt * char.Speed
	}

	char.Skeleton.Local.Translate.Set(float32(center.X), float32(center.Y))
	char.Skeleton.Local.Scale.Set(10, 10)
	char.Animation.Apply(char.Skeleton, float32(char.Time), true)
	char.Skeleton.Update()
}

type AgentHandlers struct {
	Draw types.AgentDrawHandlerID
	AI   types.AgentAiHandlerID
}

type AgentTimers struct {
	TimeSinceDeath float32
	WaitingFor     float32
}

// func newAgent(id int) *Agent {
// 	// agent := Agent{
// 	// 	AgentCategory:     constants.AGENT_CATEGORY_UNDEFINED,
// 	// 	AiHandlerID:       constants.AI_HANDLER_NULL,
// 	// 	DrawHandlerID:     constants.DRAW_HANDLER_NULL,
// 	// 	PhysicsState:      physics.Body{},
// 	// 	PhysicsParameters: physics.PhysicsParameters{Radius: 5},
// 	// }

// 	return &agent
// }

func (a *Agent) FillDefaults() {
	// if have no direction, default to right-facing / no flip
	if a.Transform.Direction == 0 {
		a.Transform.Direction = 1
	}
}

//prefabs

func (a *Agent) SetPosition(x, y float32) {
	a.Transform.Pos.X = x
	a.Transform.Pos.Y = y
}

func (a *Agent) SetSize(x, y float32) {
	a.Transform.Size.X = x
	a.Transform.Size.Y = y
}

func (a *Agent) SetVelocity(x, y float32) {
	a.Transform.Vel.X = x
	a.Transform.Vel.Y = y
}

func (a *Agent) TakeDamage(amount int) {
	a.Health.Current -= amount
	if a.Health.Current <= 0 {
		a.Health.Died = true
	}
}

func (a *Agent) Died() bool {
	return a.Health.Died
}

func (a *Agent) IsWaiting() bool {
	return a.Timers.WaitingFor > 0
}

func (a *Agent) WaitFor(seconds float32) {
	a.Timers.WaitingFor = seconds
}

func (a *Agent) Validate() {
	if a.Meta.Category == constants.AGENT_CATEGORY_UNDEFINED {
		log.Fatalf("Cannot create agent with undefined category: %+v", a)
	}
}
