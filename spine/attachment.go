package spine

type Attachment interface {
	GetName() string
}

type RegionAttachment struct {
	Name  string
	Path  string
	Size  Vector
	Local Transform
	Color Color
}

func (attach *RegionAttachment) GetName() string { return attach.Name }
func (attach *RegionAttachment) GetImage() string {
	if attach.Path != "" {
		return attach.Path
	}
	return attach.Name
}

type PointAttachment struct {
	Name  string
	Local Transform
	Color Color
}

func (attach *PointAttachment) GetName() string { return attach.Name }

type BoundingBoxAttachment struct{ VertexAttachment }

type MeshAttachment struct {
	VertexAttachment

	Hull      int
	UV        []Vector
	Triangles [][3]int
	Edges     [][2]int
}

type VertexAttachment struct {
	Name  string
	Path  string
	Local Transform
	Size  Vector
	Color Color

	Weighted     bool
	BindingCount int
	Vertices     []Vertex
}

func (attach *VertexAttachment) GetName() string { return attach.Name }
func (attach *VertexAttachment) GetImage() string {
	if attach.Path != "" {
		return attach.Path
	}
	return attach.Name
}

func (attach *VertexAttachment) CalculateWorldVertices(skel *Skeleton, slot *Slot) []Vector {
	vertices := make([]Vector, len(attach.Vertices))
	if !attach.Weighted {
		final := slot.Bone.World.Mul(attach.Local.Affine())
		if len(slot.Deform) == 0 {
			for i := range attach.Vertices {
				p := attach.Vertices[i].Position
				vertices[i] = final.Transform(p)
			}
		} else {
			if len(vertices) != len(slot.Deform) {
				panic("invalid deform")
			}
			for i := range attach.Vertices {
				p := attach.Vertices[i].Position
				vertices[i] = final.Transform(p.Add(slot.Deform[i]))
			}
		}
	} else {
		bones := skel.Bones
		if len(slot.Deform) == 0 {
			for i := range attach.Vertices {
				v := &attach.Vertices[i]
				var w Vector
				for _, binding := range v.Bindings {
					bone := bones[binding.Bone]
					p := binding.Position
					w = w.Add(bone.World.WeightedTransform(binding.Weight, p))
				}
				vertices[i] = w
			}
		} else {
			deform, di := slot.Deform, 0
			for i := range attach.Vertices {
				v := &attach.Vertices[i]
				var w Vector
				for _, binding := range v.Bindings {
					bone := bones[binding.Bone]
					p := binding.Position.Add(deform[di])
					di++
					w = w.Add(bone.World.WeightedTransform(binding.Weight, p))
				}
				vertices[i] = w
			}
		}
	}
	return vertices
}

type Vertex struct {
	Position Vector
	Bindings []VertexBinding
}

type VertexBinding struct {
	Bone     int
	Position Vector
	Weight   float32
}
