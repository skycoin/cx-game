package spine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
)

func ReadJSON(r io.Reader) (*SkeletonData, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	skeleton := &SkeletonData{}
	err = skeleton.load(&object{data})
	return skeleton, err
}

func (skel *SkeletonData) load(obj *object) error {
	info := obj.Child("skeleton")
	skel.Size.X = info.Float("width", 0)
	skel.Size.Y = info.Float("height", 0)
	skel.Images = info.String("images", "")

	// parse bones
	for _, info := range obj.List("bones") {
		bone := &BoneData{}
		bone.Index = len(skel.Bones)
		skel.Bones = append(skel.Bones, bone)

		bone.Name = info.String("name", "")
		bone.Parent = skel.FindBone(info.String("parent", ""))
		if bone.Parent != nil {
			bone.Parent.Children = append(bone.Parent.Children, bone)
		}
		bone.Length = info.Float("length", 1)
		bone.Local = info.Transform()

		bone.Color = info.Color("color", RGBA(0x98/0xFF, 0x98/0xFF, 0x98/0xFF, 0xFF/0xFF))

		bone.Inherit = parseInherit(info.String("transform", ""))
		if info.Has("inheritScale") {
			if info.Bool("inheritScale", true) {
				bone.Inherit |= InheritScale
			} else {
				bone.Inherit &^= InheritScale
			}
		}
		if info.Has("inheritRotation") {
			if info.Bool("inheritRotation", true) {
				bone.Inherit |= InheritRotation
			} else {
				bone.Inherit &^= InheritRotation
			}
		}
	}

	// parse slots
	for _, slotInfo := range obj.List("slots") {
		slot := &SlotData{}
		slot.Index = len(skel.Slots)
		skel.Slots = append(skel.Slots, slot)

		slot.Name = slotInfo.String("name", "")
		slot.Bone = skel.FindBone(slotInfo.String("bone", ""))
		slot.Attachment = slotInfo.String("attachment", "")
		slot.Color = slotInfo.Color("color", RGBA(1, 1, 1, 1))
		slot.Dark = slotInfo.Color("dark", RGBA(1, 1, 1, 1))
		slot.Blend = parseBlendMode(slotInfo.String("blend", ""))
	}

	for eventName, eventInfo := range obj.Map("events") {
		event := &EventData{}
		event.Index = len(skel.Events)
		skel.Events = append(skel.Events, event)

		event.Name = eventInfo.String("name", eventName)
		event.Int = eventInfo.Int("int", 0)
		event.Float = eventInfo.Float("float", 0)
		event.String = eventInfo.String("string", "")
	}

	for _, transformInfo := range obj.List("transform") {
		transform := &TransformConstraintData{}
		transform.Index = len(skel.TransformConstraints)
		skel.TransformConstraints = append(skel.TransformConstraints, transform)

		transform.Name = transformInfo.String("name", "")
		transform.Order = transformInfo.Int("order", 0)
		if transformInfo.Has("bones") {
			for _, boneName := range transformInfo.Strings("bones") {
				bone := skel.FindBone(boneName)
				if bone == nil {
					return errors.New("bone missing: " + boneName)
				}
				transform.Bones = append(transform.Bones, bone)
			}
		} else {
			boneName := transformInfo.String("bone", "")
			bone := skel.FindBone(boneName)
			if bone == nil {
				return errors.New("bone missing: " + boneName)
			}
			transform.Bones = []*BoneData{bone}
		}

		transform.Target = skel.FindBone(transformInfo.String("target", ""))
		transform.Offset = transformInfo.Transform()
		transform.Offset.Scale.X = transformInfo.Float("scaleX", 0)
		transform.Offset.Scale.Y = transformInfo.Float("scaleY", 0)

		transform.Mix.Rotate = transformInfo.Float("rotateMix", 1)
		transform.Mix.Translate = transformInfo.Float("translateMix", 1)
		transform.Mix.Scale = transformInfo.Float("scaleMix", 1)
		transform.Mix.Shear = transformInfo.Float("shearMix", 1)

		transform.Local = transformInfo.Bool("local", false)
		transform.Relative = transformInfo.Bool("relative", false)
	}

	// parse skins
	for skinName, skinInfo := range obj.Map("skins") {
		skin := NewSkin(skinName)
		skel.Skins = append(skel.Skins, skin)
		if skinName == "default" {
			skel.DefaultSkin = skin
		}

		for slotName, attachments := range skinInfo.Children() {
			slot := skel.FindSlot(slotName)
			if slot == nil {
				return errors.New("slot missing")
			}

			for attachmentName, attachInfo := range attachments.Children() {
				attachType := attachInfo.String("type", "region")
				attachName := attachInfo.String("name", attachmentName)
				switch attachType {
				case "region":
					attach := &RegionAttachment{}
					attach.Name = attachName
					attach.Path = attachInfo.String("path", "")
					attach.Size.X = attachInfo.Float("width", 0)
					attach.Size.Y = attachInfo.Float("height", 0)
					attach.Local = attachInfo.Transform()
					attach.Color = attachInfo.Color("color", RGBA(1, 1, 1, 1))

					skin.AddAttachment(slot.Index, attachmentName, attach)
				case "point":
					attach := &PointAttachment{}
					attach.Name = attachName
					attach.Local = attachInfo.Transform()
					attach.Color = attachInfo.Color("color", RGBA(1, 1, 1, 1))

					skin.AddAttachment(slot.Index, attachmentName, attach)
				case "mesh":
					attach := &MeshAttachment{}
					attach.Name = attachName
					attach.Path = attachInfo.String("path", "")
					attach.Size.X = attachInfo.Float("width", 0)
					attach.Size.Y = attachInfo.Float("height", 0)
					attach.Local = attachInfo.Transform()
					attach.Color = attachInfo.Color("color", RGBA(1, 1, 1, 1))

					triangles := attachInfo.Ints("triangles")
					attach.Triangles = make([][3]int, len(triangles)/3)
					for i := range attach.Triangles {
						attach.Triangles[i] = [3]int{triangles[i*3], triangles[i*3+1], triangles[i*3+2]}
					}

					edges := attachInfo.Ints("edges")
					attach.Edges = make([][2]int, len(edges)/2)
					for i := range attach.Edges {
						attach.Edges[i] = [2]int{edges[i*2], edges[i*2+1]}
					}

					uvs := attachInfo.Floats("uvs")
					vertices := attachInfo.Floats("vertices")
					attach.Weighted = len(vertices) > len(uvs)
					attach.UV = make([]Vector, len(uvs)/2)
					for i := range attach.UV {
						attach.UV[i] = Vector{uvs[i*2], uvs[i*2+1]}
					}
					//TODO: error checks for invalid vertices array

					attach.Hull = attachInfo.Int("hull", 0)

					if attach.Weighted {
						attach.Vertices = make([]Vertex, len(attach.UV))
						k := 0
						for i := range attach.Vertices {
							vertex := &attach.Vertices[i]
							boneCount := int(vertices[k])
							k++
							vertex.Bindings = make([]VertexBinding, boneCount)
							for l := range vertex.Bindings {
								bind := &vertex.Bindings[l]
								bind.Bone = int(vertices[k])
								bind.Position = Vector{vertices[k+1], vertices[k+2]}
								bind.Weight = vertices[k+3]
								k += 4
								attach.BindingCount++
							}
						}
					} else {
						attach.Vertices = make([]Vertex, len(vertices)/2)
						for i := range attach.Vertices {
							attach.Vertices[i].Position = Vector{vertices[i*2], vertices[i*2+1]}
							attach.BindingCount++
						}
					}

					if len(attach.Vertices) != len(attach.UV) {
						errors.New("vertices and uv len mismatch")
					}
					skin.AddAttachment(slot.Index, attachmentName, attach)
				case "boundingbox":
					attach := &BoundingBoxAttachment{}
					attach.Name = attachName
					attach.Color = attachInfo.Color("color", RGBA(0x60/0xFF, 0xF0/0xFF, 0, 1))

					vertexCount := attachInfo.Int("vertexCount", 0)
					vertices := attachInfo.Floats("vertices")
					if vertexCount <= 0 {
						vertexCount = len(vertices)
					}

					attach.Weighted = len(vertices) > vertexCount*2
					if attach.Weighted {
						attach.Vertices = make([]Vertex, vertexCount)
						k := 0
						for i := range attach.Vertices {
							vertex := &attach.Vertices[i]
							boneCount := int(vertices[k])
							k++
							vertex.Bindings = make([]VertexBinding, boneCount)
							for l := range vertex.Bindings {
								bind := &vertex.Bindings[l]
								bind.Bone = int(vertices[k])
								bind.Position = Vector{vertices[k+1], vertices[k+2]}
								bind.Weight = vertices[k+3]
								k += 4
								attach.BindingCount++
							}
						}
					} else {
						attach.Vertices = make([]Vertex, len(vertices)/2)
						for i := range attach.Vertices {
							attach.Vertices[i].Position = Vector{vertices[i*2], vertices[i*2+1]}
							attach.BindingCount++
						}
					}

					if len(attach.Vertices) != vertexCount {
						errors.New("vertices and vertexCount mismatch")
					}
					skin.AddAttachment(slot.Index, attachmentName, attach)
				default:
					return errors.New("unhandled attachment type: " + attachType)
				}
			}
		}
	}

	// parse animations
	for animName, animInfo := range obj.Map("animations") {
		anim := &Animation{}
		anim.Name = animName
		skel.Animations = append(skel.Animations, anim)

		for boneName, timelineInfos := range animInfo.Map("bones") {
			bone := skel.FindBone(boneName)

			if timelineInfos.Has("rotate") {
				timeline := &RotateTimeline{}
				anim.Timelines = append(anim.Timelines, timeline)
				timeline.Bone = bone.Index
				for _, info := range timelineInfos.List("rotate") {
					info.addTimeStep(&timeline.CurveTimeline)
					angle := info.Float("angle", 0) * math.Pi * 2 / 360
					timeline.Angle = append(timeline.Angle, angle)
				}
			}

			if timelineInfos.Has("translate") {
				timeline := &TranslateTimeline{}
				anim.Timelines = append(anim.Timelines, timeline)
				timeline.Bone = bone.Index
				for _, info := range timelineInfos.List("translate") {
					info.addTimeStep(&timeline.CurveTimeline)
					x, y := info.Float("x", 0), info.Float("y", 0)
					timeline.Translate = append(timeline.Translate, Vector{x, y})
				}
			}

			if timelineInfos.Has("scale") {
				timeline := &ScaleTimeline{}
				anim.Timelines = append(anim.Timelines, timeline)
				timeline.Bone = bone.Index
				for _, info := range timelineInfos.List("scale") {
					info.addTimeStep(&timeline.CurveTimeline)
					x, y := info.Float("x", 0), info.Float("y", 0)
					timeline.Scale = append(timeline.Scale, Vector{x, y})
				}
			}

			if timelineInfos.Has("shear") {
				timeline := &ShearTimeline{}
				anim.Timelines = append(anim.Timelines, timeline)
				timeline.Bone = bone.Index
				for _, info := range timelineInfos.List("shear") {
					info.addTimeStep(&timeline.CurveTimeline)
					x := info.Float("x", 0) * math.Pi * 2 / 360
					y := info.Float("y", 0) * math.Pi * 2 / 360
					timeline.Shear = append(timeline.Shear, Vector{x, y})
				}
			}
		}

		for slotName, slotTimelines := range animInfo.Map("slots") {
			slot := skel.FindSlot(slotName)

			if slotTimelines.Has("attachment") {
				timeline := &AttachmentTimeline{}
				anim.Timelines = append(anim.Timelines, timeline)
				timeline.Slot = slot.Index
				for _, info := range slotTimelines.List("attachment") {
					info.addTimeStep(&timeline.CurveTimeline)
					attachment := info.String("name", "")
					timeline.Attachment = append(timeline.Attachment, attachment)
				}
			}

			if slotTimelines.Has("color") {
				timeline := &ColorTimeline{}
				timeline.Slot = slot.Index
				anim.Timelines = append(anim.Timelines, timeline)
				for _, info := range slotTimelines.List("color") {
					info.addTimeStep(&timeline.CurveTimeline)
					color := info.Color("color", RGBA(1, 1, 1, 1))
					timeline.Color = append(timeline.Color, color)
				}
			}

			if slotTimelines.Has("twoColor") {
				timeline := &TwoColorTimeline{}
				timeline.Slot = slot.Index
				anim.Timelines = append(anim.Timelines, timeline)
				for _, info := range slotTimelines.List("twoColor") {
					info.addTimeStep(&timeline.CurveTimeline)
					light := info.Color("light", RGBA(1, 1, 1, 1))
					dark := info.Color("dark", RGBA(0, 0, 0, 1))
					timeline.Color = append(timeline.Color, [2]Color{light, dark})
				}
			}
		}

		for constraintName, timelineInfo := range animInfo.Child("transform").ChildrenArray() {
			constraint := skel.FindTransformConstraint(constraintName)
			timeline := &TransformConstraintTimeline{}
			timeline.Constraint = constraint.Index
			anim.Timelines = append(anim.Timelines, timeline)

			for _, info := range timelineInfo {
				info.addTimeStep(&timeline.CurveTimeline)
				mix := TransformMix{}
				mix.Rotate = info.Float("rotateMix", 1)
				mix.Translate = info.Float("translateMix", 1)
				mix.Scale = info.Float("scaleMix", 1)
				mix.Shear = info.Float("shearMix", 1)
				timeline.Transform = append(timeline.Transform, mix)
			}
		}

		for skinName, skinTimelines := range animInfo.Map("deform") {
			skin := skel.FindSkin(skinName)
			for slotName, slotTimelines := range skinTimelines.Children() {
				slot := skel.FindSlot(slotName)
				for meshName, meshTimelines := range slotTimelines.ChildrenArray() {
					attachment := skin.Attachment(slot.Index, meshName)
					if attachment == nil {
						attachment = skel.DefaultSkin.Attachment(slot.Index, meshName)
					}
					if attachment == nil {
						return errors.New("attachment missing")
					}
					mesh, ok := attachment.(*MeshAttachment)
					if !ok {
						return errors.New("attachment not mesh")
					}

					timeline := &DeformTimeline{}
					anim.Timelines = append(anim.Timelines, timeline)
					timeline.Skin = skin
					timeline.Slot = slot.Index
					timeline.Mesh = meshName
					for _, info := range meshTimelines {
						info.addTimeStep(&timeline.CurveTimeline)
						offset := info.Int("offset", 0)
						vertices := info.Floats("vertices")
						deform := make([]Vector, mesh.BindingCount)
						for i := 0; i < len(vertices)/2; i++ {
							deform[offset/2+i] = Vector{vertices[i*2], vertices[i*2+1]}
						}
						if len(deform) != mesh.BindingCount {
							return errors.New("invalid deform length")
						}
						timeline.Deform = append(timeline.Deform, deform)
					}
				}
			}
		}

		orderKey := "drawOrder"
		if !animInfo.Has(orderKey) {
			orderKey = "draworder"
		}
		if animInfo.Has(orderKey) {
			timeline := &OrderTimeline{}
			anim.Timelines = append(anim.Timelines, timeline)
			for _, eventInfo := range animInfo.List(orderKey) {
				info.addTimeStep(&timeline.CurveTimeline)

				var order []int
				if eventInfo.Has("offsets") {
					order = make([]int, len(skel.Slots))
					for i := range order {
						order[i] = -1
					}

					offsets := eventInfo.List("offsets")
					unchanged := make([]int, len(order)-len(offsets))
					originalIndex, unchangedIndex := 0, 0

					for _, offsetInfo := range offsets {
						slotName := offsetInfo.String("slot", "")
						slotIndex := skel.FindSlot(slotName).Index

						for originalIndex != slotIndex {
							unchanged[unchangedIndex] = originalIndex
							unchangedIndex++
							originalIndex++
						}

						index := originalIndex + offsetInfo.Int("offset", 0)
						order[index] = originalIndex
						originalIndex++
					}
					for originalIndex < len(order) {
						unchanged[unchangedIndex] = originalIndex
						unchangedIndex++
						originalIndex++
					}
					for i := len(order) - 1; i >= 0; i-- {
						if order[i] == -1 {
							unchangedIndex--
							order[i] = unchanged[unchangedIndex]
						}
					}
				}
				timeline.Order = append(timeline.Order, order)
			}
		}

		if animInfo.Has("events") {
			timeline := &EventTimeline{}
			anim.Timelines = append(anim.Timelines, timeline)
			for _, eventInfo := range animInfo.List("events") {
				info.addTimeStep(&timeline.CurveTimeline)
				eventName := eventInfo.String("string", "")
				timeline.Event = append(timeline.Event, eventName)
			}
		}

		for _, timeline := range anim.Timelines {
			if timeline.Duration() > anim.Duration {
				anim.Duration = timeline.Duration()
			}
		}
	}

	// TODO: make order dependent loading
	sort.Slice(skel.Skins, func(i, k int) bool {
		return skel.Skins[i].Name < skel.Skins[k].Name
	})
	sort.Slice(skel.Animations, func(i, k int) bool {
		return skel.Animations[i].Name < skel.Animations[k].Name
	})

	// TODO: sanity checks

	if len(skel.Skins) == 0 {
		return errors.New("no skins defined")
	}
	if skel.DefaultSkin == nil {
		skel.DefaultSkin = skel.Skins[0]
	}

	return nil
}

func (obj *object) addTimeStep(timeline *CurveTimeline) {
	time := obj.Float("time", 0)
	timeline.Time = append(timeline.Time, time)

	if !obj.Has("curve") {
		if len(timeline.Curve) > 0 {
			curve := Curve{}
			curve.SetLinear()
			timeline.Curve = append(timeline.Curve, curve)
		}
		return
	}

	if len(timeline.Curve) == 0 {
		timeline.Curve = make([]Curve, len(timeline.Time)-1)
		for i := range timeline.Curve {
			timeline.Curve[i].SetLinear()
		}
	}

	curve := Curve{}
	curve.SetLinear()

	standard := obj.String("curve", "")
	if standard == "linear" || standard == "step" || standard == "stepped" {
		if standard == "step" || standard == "stepped" {
			curve.SetStep()
		}
		timeline.Curve = append(timeline.Curve, curve)
		return
	}
	if standard != "" {
		panic("unimplemented curve " + standard)
	}

	if any, ok := obj.Any("curve"); ok {
		if items, ok := any.([]interface{}); ok && len(items) == 4 {
			x0, ok0 := items[0].(float64)
			y0, ok1 := items[1].(float64)
			x1, ok2 := items[2].(float64)
			y1, ok3 := items[3].(float64)
			if ok0 && ok1 && ok2 && ok3 {
				curve.SetBezier(
					Vector{float32(x0), float32(y0)},
					Vector{float32(x1), float32(y1)})
				timeline.Curve = append(timeline.Curve, curve)
				return
			}
		}
	}

	timeline.Curve = append(timeline.Curve, curve)
}

func (obj *object) Transform() Transform {
	var transform Transform
	transform.Translate.X = obj.Float("x", 0)
	transform.Translate.Y = obj.Float("y", 0)
	transform.Rotate = obj.Float("rotation", 0) * math.Pi * 2 / 360
	transform.Scale.X = obj.Float("scaleX", 1)
	transform.Scale.Y = obj.Float("scaleY", 1)
	transform.Shear.X = obj.Float("shearX", 0) * math.Pi * 2 / 360
	transform.Shear.Y = obj.Float("shearY", 0) * math.Pi * 2 / 360
	return transform
}

type object struct {
	data map[string]interface{}
}

func (obj *object) Any(key string) (interface{}, bool) {
	if obj == nil {
		return nil, false
	}
	value, ok := obj.data[key]
	return value, ok
}

func (obj *object) Has(key string) bool {
	_, ok := obj.data[key]
	return ok
}

func (obj *object) Child(key string) *object {
	if obj == nil {
		return nil
	}
	if any, ok := obj.Any(key); ok {
		if data, ok := any.(map[string]interface{}); ok {
			return &object{data}
		}
	}
	return nil
}

func (obj *object) List(key string) []*object {
	if obj == nil {
		return nil
	}
	if any, ok := obj.Any(key); ok {
		if anyitems, ok := any.([]interface{}); ok {
			xs := []*object{}
			for _, item := range anyitems {
				if data, ok := item.(map[string]interface{}); ok {
					xs = append(xs, &object{data})
				}
			}
			return xs
		}
	}
	return nil
}

func (obj *object) Floats(key string) []float32 {
	if obj == nil {
		return nil
	}
	if any, ok := obj.Any(key); ok {
		if anyitems, ok := any.([]interface{}); ok {
			xs := []float32{}
			for _, item := range anyitems {
				if data, ok := item.(float32); ok {
					xs = append(xs, data)
				} else if data, ok := item.(float64); ok {
					xs = append(xs, float32(data))
				} else if data, ok := item.(int32); ok {
					xs = append(xs, float32(data))
				} else if data, ok := item.(int64); ok {
					xs = append(xs, float32(data))
				} else {
					panic(fmt.Sprintf("unexpected %v", item))
				}
			}
			return xs
		}
	}
	return nil
}

func (obj *object) Ints(key string) []int {
	if obj == nil {
		return nil
	}
	if any, ok := obj.Any(key); ok {
		if anyitems, ok := any.([]interface{}); ok {
			xs := []int{}
			for _, item := range anyitems {
				if data, ok := item.(int); ok {
					xs = append(xs, data)
				} else if data, ok := item.(float64); ok {
					xs = append(xs, int(data))
				} else if data, ok := item.(float32); ok {
					xs = append(xs, int(data))
				} else if data, ok := item.(int32); ok {
					xs = append(xs, int(data))
				} else if data, ok := item.(int64); ok {
					xs = append(xs, int(data))
				} else {
					panic(fmt.Sprintf("unexpected %v", item))
				}
			}
			return xs
		}
	}
	return nil
}

func (obj *object) Strings(key string) []string {
	if obj == nil {
		return nil
	}
	if any, ok := obj.Any(key); ok {
		if anyitems, ok := any.([]interface{}); ok {
			xs := []string{}
			for _, item := range anyitems {
				if data, ok := item.(string); ok {
					xs = append(xs, data)
				} else {
					panic(fmt.Sprintf("unexpected %v", item))
				}
			}
			return xs
		}
	}
	return nil
}

func (obj *object) Map(key string) map[string]*object {
	return obj.Child(key).Children()
}

func (obj *object) Children() map[string]*object {
	if obj == nil {
		return nil
	}
	xs := map[string]*object{}
	for key, item := range obj.data {
		if data, ok := item.(map[string]interface{}); ok {
			xs[key] = &object{data}
		}
	}
	return xs
}

func (obj *object) ChildrenArray() map[string][]*object {
	if obj == nil {
		return nil
	}
	xs := map[string][]*object{}
	for key, item := range obj.data {
		if data, ok := item.([]interface{}); ok {
			xss := []*object{}
			for _, childitem := range data {
				if child, ok := childitem.(map[string]interface{}); ok {
					xss = append(xss, &object{child})
				}
			}
			xs[key] = xss
		}
	}
	return xs
}

func (obj *object) String(key string, def string) string {
	if any, ok := obj.Any(key); ok {
		if str, ok := any.(string); ok {
			return str
		}
	}
	return def
}

func (obj *object) Int(key string, def int) int {
	if any, ok := obj.Any(key); ok {
		if data, ok := any.(int); ok {
			return data
		} else if data, ok := any.(int64); ok {
			return int(data)
		} else if data, ok := any.(int32); ok {
			return int(data)
		} else if data, ok := any.(float64); ok {
			return int(data)
		} else if data, ok := any.(float32); ok {
			return int(data)
		} else if str, ok := any.(string); ok {
			if v, err := strconv.ParseInt(str, 10, 64); err == nil {
				return int(v)
			}
		}
	}
	return def
}

func (obj *object) Float(key string, def float32) float32 {
	if any, ok := obj.Any(key); ok {
		if data, ok := any.(float32); ok {
			return data
		} else if data, ok := any.(float64); ok {
			return float32(data)
		} else if str, ok := any.(string); ok {
			if v, err := strconv.ParseFloat(str, 32); err == nil {
				return float32(v)
			}
		}
	}
	return def
}

func (obj *object) Bool(key string, def bool) bool {
	if any, ok := obj.Any(key); ok {
		if data, ok := any.(bool); ok {
			return data
		} else if str, ok := any.(string); ok {
			if str == "true" {
				return true
			} else if str == "false" {
				return false
			}
		}
	}
	return def
}

func (obj *object) Color(key string, def Color) Color {
	//TODO: error handling
	if any, ok := obj.Any(key); ok {
		if str, ok := any.(string); ok {
			color := def
			if len(str) == 3 || len(str) == 4 {
				color.R = parseHex1(str[0:1])
				color.G = parseHex1(str[1:2])
				color.B = parseHex1(str[2:3])
				if len(str) == 4 {
					color = color.WithAlpha(parseHex1(str[3:4]))
				} else {
					color.A = 1
				}
			} else if len(str) == 6 || len(str) == 8 {
				color.R = parseHex2(str[0:2])
				color.G = parseHex2(str[2:4])
				color.B = parseHex2(str[4:6])
				if len(str) == 8 {
					color = color.WithAlpha(parseHex2(str[6:8]))
				} else {
					color.A = 1
				}
			}
			return color
		}
	}
	return def
}

func parseHex1(s string) float32 {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic(err)
		return 1
	}
	return float32(v) / 0xF
}

func parseHex2(s string) float32 {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic(err)
		return 1
	}
	return float32(v) / 0xFF
}
