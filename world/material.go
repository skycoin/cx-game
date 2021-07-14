package world

type MaterialID uint32
type Material struct {
	Name string
}
var materials = []Material{}

func (id MaterialID) Get() Material {
	return materials[id]
}

func RegisterMaterial(material Material) MaterialID {
	id := MaterialID(len(materials))
	materials = append(materials,material)
	return id
}

func init() {
	RegisterMaterial(Material{Name:"Null Material"})
}
