package agent_draw

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

var moving bool

//var bodyMultiplier float32

const (
	degree = 180 / math.Pi
	radian = math.Pi / 180
)

func spineDrawPlayerSprite(
	agent *agents.Agent, ctx DrawHandlerContext,
	spriteID render.SpriteID, zOffset float32,
) {

	body := &agent.Transform

	bodyMultiplier := body.Size.Y / agent.Meta.Skeleton.Data.Size.Y

	// fmt.Println(body.Size.Y)
	// fmt.Println(agent.Skeleton.Data.Size.Y)
	// fmt.Println(bodyMultiplier)
	// agent.Skeleton.Data.Size.X = agent.Skeleton.Data.Size.X // * bodyMultiplier
	// agent.Skeleton.Data.Size.Y = agent.Skeleton.Data.Size.Y // * bodyMultiplier

	//drawn one frame behind
	alpha := timer.GetTimeBetweenTicks() / constants.MS_PER_TICK

	// var interpolatedPos cxmath.Vec2
	// if !body.Pos.Equal(body.PrevPos) {
	// 	interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))
	// } else {
	// 	interpolatedPos = body.Pos

	// }

	// translate := mgl32.Translate3D(
	// 	interpolatedPos.X,
	// 	interpolatedPos.Y,
	// 	zOffset,
	// )

	//	hitboxToRender := 1 / constants.PLAYER_RENDER_TO_HITBOX
	// scaleX := -body.Size.X * body.Direction * hitboxToRender
	//scaleX := -body.Size.X //* hitboxToRender
	// scale := mgl32.Scale3D(scaleX, agent.Transform.Size.Y, 1)

	// agent.Skeleton.Data.Size.X = scaleX
	// agent.Skeleton.Data.Size.Y = agent.Transform.Size.Y

	if body.Direction != 1 {
		agent.Meta.Skeleton.FlipX = false
	} else {
		agent.Meta.Skeleton.FlipX = true
	}

	if body.IsOnGround() && input.GetAxis(input.HORIZONTAL) != 0 {
		if moving == false {
			agent.Meta.SetAnimation(6)
			agent.Meta.Play = true
			moving = true
		}

	} else {
		if moving == true {
			agent.Meta.SetAnimation(2)
			moving = false
		}

	}

	//transform := translate.Mul4(scale)
	// agent.Meta.Skeleton.Local.Translate.Set(float32(body.Pos.X), float32(body.Pos.Y-(body.Size.Y/2)))
	//	agent.Skeleton.Local.Scale.Set(scaleX*0.025, agent.Transform.Size.Y*0.025)
	// agent.Meta.Skeleton.Local.Scale.Set(-1*bodyMultiplier, 1*bodyMultiplier)
	// agent.Skeleton.Local.Scale.Set(-1*0.077, 1*0.077)

	skeletonTraslate := cxmath.Vec2{body.Pos.X, float32(body.Pos.Y - (body.Size.Y / 2))}
	skeletonScale := cxmath.Vec2{-1 * bodyMultiplier, 1 * bodyMultiplier}
	// fmt.Println(agent.Skeleton.Data.Size)
	// fmt.Println(agent.Transform.Size)

	agent.Meta.Update(float64(alpha), skeletonTraslate, skeletonScale)
	// agent.Skeleton.Local.Scale.Set(scaleX, agent.Transform.Size.Y)
	//agent.Transform.Size = cxmath.Vec2{agent.Skeleton.Data.Size.X, agent.Skeleton.Data.Size.Y}
	for i, slot := range agent.Meta.Skeleton.Order {
		bone := slot.Bone
		switch attachment := slot.Attachment.(type) {
		case nil:
		case *spine.RegionAttachment:
			if attachment.Name == "head1" ||
				attachment.Name == "front_arm1" ||
				attachment.Name == "front_hand1" ||
				attachment.Name == "front_foot1" {
				break
			}
			spriteID := agent.Meta.GetImage(attachment.Name)

			local := attachment.Local.Affine()
			final := bone.World.Mul(local)

			// S_T := GetSpriteTransform(final)
			S_T := GetSpriteTransform2(final.Col64())

			spriteTransform := mgl32.Ident4().
				Mul4(mgl32.Translate3D(S_T.Translation.X(), S_T.Translation.Y(), zOffset+float32(i))).
				Mul4(mgl32.HomogRotate3DZ(S_T.Rotation)).

				// Mul4(mgl32.Scale3D(-spriteSize[0]*0.065*body.Direction*hitboxToRender, spriteSize[1]*0.065, 1))
				//Mul4(mgl32.Scale3D(S_T.Scale.X()*9, S_T.Scale.Y()*9, 1))
				Mul4(mgl32.Scale3D(1/S_T.Scale.X()*bodyMultiplier*0.5, 1/S_T.Scale.Y()*bodyMultiplier*0.5, 1)).
				Mul4(mgl32.ShearZ3D(S_T.Skew.X(), S_T.Skew.Y()))
			// set blending mode
			// BUG: incorrect, should use blending mode not compositing mode
			switch slot.Data.Blend {
			case spine.Normal:
				// MISSING
			case spine.Additive:
			//	draw.CompositeMode = ebiten.CompositeModeLighter
			case spine.Multiply:
				// MISSING
			case spine.Screen:
				// MISSING
				//	draw.CompositeMode = ebiten.CompositeModeLighter
			}

			render.DrawWorldSprite(spriteTransform, spriteID, render.NewSpriteDrawOptions())
			//_, _, _, _, tx, ty := draw.GeoM.GetElements()
			// poject := Project(mgl32.Vec2{final.Translation().X, final.Translation().Y}, final)

			// if attachment.Name == "goggles" {
			// 	fmt.Println("goggles sizex:", m.M_width)
			// 	fmt.Println("goggles sizey:", m.M_height)
			// 	fmt.Println("goggles transform:", final.Translation().Y)

			// }
			//if attachment.Name == "gun" || attachment.Name == "head" || attachment.Name == "goggles" {
			// q := CreateQuad(float32(m.M_height), float32(m.M_width), float32(i), xform, m)
			// pos = append(pos, q...)
			//}
			// fmt.Println(i)
			//offset += 50.0
			//target.DrawImage(m, &draw)
		default:
			panic(fmt.Sprintf("unknown attachment %v", attachment))
		}
	}

	//	render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions())
}

// type Vertex struct {
// 	Position  mgl32.Vec2
// 	TexCoords mgl32.Vec2
// }
type SpriteTransform struct {
	Translation mgl32.Vec2
	Rotation    float32
	Scale       mgl32.Vec2
	Skew        mgl32.Vec2
}

func GetSpriteTransform(mat spine.Affine) SpriteTransform {
	var a = mat.M00
	var b = mat.M01
	var c = mat.M10
	var d = mat.M11
	var e = mat.M02
	var f = mat.M12

	var delta = a*d - b*c

	result := SpriteTransform{
		Translation: mgl32.Vec2{e, f},
		Rotation:    0,
		Scale:       mgl32.Vec2{0, 0},
		Skew:        mgl32.Vec2{0, 0},
	}

	// Apply the QR-like decomposition.
	if a != 0 || b != 0 {

		var r = float32(math.Sqrt(float64(a*a + b*b)))

		// result.Rotation = b > 0 ? math.acos(a / r) : -math.acos(a / r);
		if b > 0 {
			result.Rotation = float32(math.Acos(float64(a / r)))
		} else {
			result.Rotation = float32(-math.Acos(float64(a / r)))
		}
		result.Scale = mgl32.Vec2{r, delta / r}
		result.Skew = mgl32.Vec2{float32(math.Atan(float64((a*c + b*d) / (r * r)))), 0}
	} else if c != 0 || d != 0 {
		var s = float32(math.Sqrt(float64(c*c + d*d)))
		//   result.Rotation =	math.Pi / 2 - (d > 0 ? math.acos(-c / s) : -math.acos(c / s));
		if d > 0 {
			result.Rotation = math.Pi/2 - float32(math.Acos(float64(-c/s)))
		} else {
			result.Rotation = math.Pi/2 - float32(-math.Acos(float64(c/s)))
		}
		result.Scale = mgl32.Vec2{delta / s, s}
		result.Skew = mgl32.Vec2{0, float32(math.Atan(float64((a*c + b*d) / (s * s))))}
	} else {
		// a = b = c = d = 0
	}

	return result
}
func GetSpriteTransform2(mat [6]float64) SpriteTransform {
	var a = float32(mat[0])
	var b = float32(mat[1])
	var c = float32(mat[2])
	var d = float32(mat[3])
	var e = float32(mat[4])
	var f = float32(mat[5])

	var delta = a*d - b*c

	result := SpriteTransform{
		Translation: mgl32.Vec2{e, f},
		Rotation:    0,
		Scale:       mgl32.Vec2{1, 1},
		Skew:        mgl32.Vec2{0, 0},
	}

	// Apply the QR-like decomposition.
	if a != 0 || b != 0 {

		var r = float32(math.Sqrt(float64(a*a + b*b)))

		// result.Rotation = b > 0 ? math.acos(a / r) : -math.acos(a / r);
		if b > 0 {
			result.Rotation = float32(math.Acos(float64(a / r)))
		} else {
			result.Rotation = float32(-math.Acos(float64(a / r)))
		}
		result.Scale = mgl32.Vec2{r, delta / r}
		result.Skew = mgl32.Vec2{float32(math.Atan(float64((a*c + b*d) / (r * r)))), 0}
	} else if c != 0 || d != 0 {
		var s = float32(math.Sqrt(float64(c*c + d*d)))
		//   result.Rotation =	math.Pi / 2 - (d > 0 ? math.acos(-c / s) : -math.acos(c / s));
		if d > 0 {
			result.Rotation = math.Pi/2 - float32(math.Acos(float64(-c/s)))
		} else {
			result.Rotation = math.Pi/2 - float32(-math.Acos(float64(-c/s)))
		}
		result.Scale = mgl32.Vec2{delta / s, s}
		result.Skew = mgl32.Vec2{0, float32(math.Atan(float64((a*c + b*d) / (s * s))))}
	} else {
		// a = b = c = d = 0
	}

	return result
}
func GetSpriteTransform3(mat [6]float64) SpriteTransform {
	var a = mat[0]
	var b = mat[1]
	var c = mat[2]
	var d = mat[3]
	var e = mat[4]
	var f = mat[5]

	// var delta = a*d - b*c

	result := SpriteTransform{
		Translation: mgl32.Vec2{float32(e), float32(f)},
		Rotation:    0,
		Scale:       mgl32.Vec2{0, 0},
		Skew:        mgl32.Vec2{0, 0},
	}

	angle := math.Atan2(b, a)
	denom := math.Pow(a, 2) + math.Pow(b, 2)
	scaleX := math.Sqrt(denom)
	scaleY := (a*d - c*b) / scaleX
	skewX := math.Atan2(a*c+b*d, denom)

	// scaleX := math.Sqrt((a * a) + (c * c))
	// scaleY := math.Sqrt((b * b) + (d * d))

	// sign := math.Atan(-c / a)
	// rad := math.Acos(a / scaleX)
	// deg := rad * degree

	// var rotation float64
	// if deg > 90 && sign > 0 {
	// 	rotation = (360 - deg) * radian
	// } else if deg < 90 && sign < 0 {
	// 	rotation = (360 - deg) * radian
	// } else {
	// 	rotation = rad
	// }

	result.Rotation = float32(angle / (math.Pi / 180))
	result.Scale = mgl32.Vec2{float32(scaleX), float32(scaleY)}
	result.Skew = mgl32.Vec2{float32(skewX / (math.Pi / 180)), 0}

	// // Apply the QR-like decomposition.
	// if a != 0 || b != 0 {

	// 	var r = float32(math.Sqrt(float64(a*a + b*b)))

	// 	// result.Rotation = b > 0 ? math.acos(a / r) : -math.acos(a / r);
	// 	if b > 0 {
	// 		result.Rotation = float32(math.Acos(float64(a / r)))
	// 	} else {
	// 		result.Rotation = float32(-math.Acos(float64(a / r)))
	// 	}
	// 	result.Scale = mgl32.Vec2{r, delta / r}
	// 	result.Skew = mgl32.Vec2{float32(math.Atan(float64((a*c + b*d) / (r * r)))), 0}
	// } else if c != 0 || d != 0 {
	// 	var s = float32(math.Sqrt(float64(c*c + d*d)))
	// 	//   result.Rotation =	math.Pi / 2 - (d > 0 ? math.acos(-c / s) : -math.acos(c / s));
	// 	if d > 0 {
	// 		result.Rotation = math.Pi/2 - float32(math.Acos(float64(-c/s)))
	// 	} else {
	// 		result.Rotation = math.Pi/2 - float32(-math.Acos(float64(-c/s)))
	// 	}
	// 	result.Scale = mgl32.Vec2{delta / s, s}
	// 	result.Skew = mgl32.Vec2{0, float32(math.Atan(float64((a*c + b*d) / (s * s))))}
	// } else {
	// 	// a = b = c = d = 0
	// }

	return result
}

// func CreateQuad(x, y, h, w float32, matrix geometry.Matrix) []float32 {

// 	var (
// 		horizontal = geometry.V(float64(w/2), 0)
// 		vertical   = geometry.V(0, float64(h/2))
// 	)
// 	var pos []float32
// 	v0 := Vertex{}
// 	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
// 	t0 := geometry.Vec{}.Sub(horizontal).Sub(vertical)
// 	xy0 := matrix.Project(t0)
// 	v0.Position = mgl32.Vec2{float32(xy0.X), float32(xy0.Y)}
// 	pos = append(pos, v0.Position.X(), v0.Position.Y(), 0, v0.TexCoords.X(), v0.TexCoords.Y())

// 	v1 := Vertex{}
// 	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
// 	t1 := geometry.Vec{}.Add(horizontal).Sub(vertical)
// 	xy1 := matrix.Project(t1)
// 	v1.Position = mgl32.Vec2{float32(xy1.X), float32(xy1.Y)}
// 	pos = append(pos, v1.Position.X(), v1.Position.Y(), 0, v1.TexCoords.X(), v1.TexCoords.Y())

// 	v2 := Vertex{}
// 	v2.Position = mgl32.Vec2{x + w, y + h}
// 	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
// 	t2 := geometry.Vec{}.Add(horizontal).Add(vertical)
// 	xy2 := matrix.Project(t2)
// 	v2.Position = mgl32.Vec2{float32(xy2.X), float32(xy2.Y)}
// 	pos = append(pos, v2.Position.X(), v2.Position.Y(), 0, v2.TexCoords.X(), v2.TexCoords.Y())

// 	v3 := Vertex{}
// 	v3.Position = mgl32.Vec2{x, y + h}
// 	v3.TexCoords = mgl32.Vec2{0.0, 1.0}
// 	t3 := geometry.Vec{}.Sub(horizontal).Add(vertical)
// 	xy3 := matrix.Project(t3)
// 	v3.Position = mgl32.Vec2{float32(xy3.X), float32(xy3.Y)}
// 	pos = append(pos, v3.Position.X(), v3.Position.Y(), 0, v3.TexCoords.X(), v3.TexCoords.Y())

// 	fmt.Printf("numbers=%v\n", pos)
// 	return pos

// }

func SpinePlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	// anim.Program.Use()
	// defer anim.Program.StopUsing()

	// gl.Enable(gl.DEPTH_TEST)

	for _, agent := range agents {
		spineDrawPlayerSprite(agent, ctx, agent.Meta.SuitSpriteID, constants.PLAYER_Z)
		//	drawPlayerSprite(agent, ctx, agent.Meta.HelmetSpriteID, constants.PLAYER_Z+1)
	}
}
