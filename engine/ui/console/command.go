package console

import (
	"log"
	"strings"

	"github.com/skycoin/cx-game/world/mapgen"
)

type Command func(string,CommandContext) string

var commands = make(map[string]Command)

func LoadMap(line string, ctx CommandContext) string {
	fname := strings.Split(line," ")[1]
	log.Printf("trying to load map from %s",fname)
	ctx.World.Load(fname)
	return ""
}

func SaveMap(line string, ctx CommandContext) string {
	fname := strings.Split(line," ")[1]
	log.Printf("trying to save map to %s",fname)
	ctx.World.Save(fname)
	return ""
}

func NewPlanet(line string, ctx CommandContext) string {
	ctx.World.Planet = *mapgen.GeneratePlanet()
	return ""
}

func Help(line string, ctx CommandContext) string {
	names := make([]string, 0, len(commands))
	for name,_ := range commands {
		names = append(names, name)
	}
	return strings.Join(names,", ")
}

func init() {
	commands["loadmap"] = LoadMap
	commands["savemap"] = SaveMap
	commands["newplanet"] = NewPlanet
	commands["help"] = Help
}

func processCommand(line string, ctx CommandContext) string {
	words := strings.Split(line," ")
	commandName := words[0]
	command,ok := commands[commandName]
	if !ok {
		log.Printf("unrecognized command [%s]",commandName)
		return ""
	}
	return command(line, ctx)
}
