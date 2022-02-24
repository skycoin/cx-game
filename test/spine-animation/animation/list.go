package animation

import "path/filepath"

type Location struct {
	Name   string
	Dir    string
	JSON   string
	Atlas  string
	Images string
}

var list = []Location{
	// {"Alien", "alien", "alien-ess.json", "alien.atlas", ""},
	// {"Alien Pro", "alien", "alien-pro.json", "alien.atlas", ""},
	// {"Goblins", "goblins", "goblins-ess.json", "goblins.atlas", ""},
	// {"Goblins Pro", "goblins", "goblins-pro.json", "goblins.atlas", ""},
	{"Powerup", "powerup", "powerup-ess.json", "powerup.atlas", ""},
	// {"Powerup Pro", "powerup", "powerup-pro.json", "powerup.atlas", ""},
	// {"Raptor", "raptor", "raptor-pro.json", "raptor.atlas", ""},
	// {"Speedy", "speedy", "speedy-ess.json", "speedy.atlas", ""},
	// {"Spineboy", "spineboy", "spineboy-ess.json", "spineboy.atlas", ""},
	//	{"Spineboy Pro", "spineboy", "spineboy-pro.json", "spineboy.atlas", ""},
	// {"Spinosaurus", "spinosaurus", "spinosaurus-ess.json", "spinosaurus-ess.atlas", ""},
	// {"Stretchyman", "stretchyman", "stretchyman-pro.json", "stretchyman.atlas", ""},
	// {"Tank", "tank", "tank-pro.json", "tank.atlas", ""},
	// {"Vine", "vine", "vine-pro.json", "vine.atlas", ""},
	// {"Robot2", "skeleton", "skeleton.json", "Robot2.atlas", ""},
}

// TODO: read this from folder structure instead
func LoadList(root string) []Location {
	xs := make([]Location, 0)
	for _, loc := range list {
		dir := filepath.Join(root, loc.Dir)
		xs = append(xs, Location{
			Name:   loc.Name,
			Dir:    dir,
			JSON:   filepath.Join(dir, "export", loc.JSON),
			Atlas:  filepath.Join(dir, "export", loc.Atlas),
			Images: filepath.Join(dir, "images"),
		})
	}
	return xs
}
