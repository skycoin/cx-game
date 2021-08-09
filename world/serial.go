package world

import (
	"log"
	"encoding/json"
	"io/ioutil"
)

func (world *World) Save(fname string) {
	buf, err := json.MarshalIndent(world, "", "  ")
	if err != nil {
		log.Fatalf("saving world: %v",err)
	}
	err = ioutil.WriteFile(fname, buf, 0644)
	if err != nil {
		log.Fatalf("writing world file: %v",err)
	}
}

func (world *World) Load(fname string) {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("reading world file: %v",err)
	}
	err = json.Unmarshal(buf, world)
	if err != nil {
		log.Fatalf("loading world: %v",err)
	}
}
