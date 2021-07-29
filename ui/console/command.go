package console

import (
	"log"
	"strings"

	"github.com/skycoin/cx-game/world/mapgen"
)

type Command func(string,CommandContext)

var commands = make(map[string]Command)

func LoadMap(line string, ctx CommandContext) {
	fname := strings.Split(line," ")[1]
	log.Printf("trying to load map from %s",fname)
	ctx.World.Load(fname)
}

func SaveMap(line string, ctx CommandContext) {
	fname := strings.Split(line," ")[1]
	log.Printf("trying to save map to %s",fname)
	ctx.World.Save(fname)
}

func NewPlanet(line string, ctx CommandContext) {
	ctx.World.Planet = *mapgen.GeneratePlanet()
}

func init() {
	commands["loadmap"] = LoadMap
	commands["savemap"] = SaveMap
	commands["newplanet"] = NewPlanet
}

func processCommand(line string, ctx CommandContext) {
	words := strings.Split(line," ")
	commandName := words[0]
	command,ok := commands[commandName]
	if !ok {
		log.Printf("unrecognized command [%s]",commandName)
		return
	}
	command(line, ctx)
}
