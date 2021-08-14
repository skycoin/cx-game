package console

import (
	"log"
	"strings"
	"strconv"

	"github.com/skycoin/cx-game/world/mapgen"
)

type Command func(string, CommandContext) string

var commands = make(map[string]Command)

func LoadMap(line string, ctx CommandContext) string {
	splitLine := strings.Split(line, " ")
	if len(splitLine) < 2 {
		log.Println("[LOADMAP ERROR] provide filename")
		return ""
	}
	fname := splitLine[1]
	log.Printf("trying to load map from %s", fname)
	ctx.World.Load(fname)
	return ""
}

func SaveMap(line string, ctx CommandContext) string {
	splitLine := strings.Split(line, " ")
	if len(splitLine) < 2 {
		log.Println("[SAVEMAP ERROR] provide filename")
		return ""
	}
	fname := splitLine[1]
	log.Printf("trying to save map to %s", fname)
	ctx.World.Save(fname)
	return ""
}

func NewPlanet(line string, ctx CommandContext) string {
	ctx.World.Planet = *mapgen.GeneratePlanet()
	return ""
}

func Help(line string, ctx CommandContext) string {
	names := make([]string, 0, len(commands))
	for name, _ := range commands {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

func Teleport(line string, ctx CommandContext) string {
	words := strings.Split(line, " ")

	x,err := strconv.ParseFloat(words[1],32)
	if err!=nil { log.Fatalf("Teleport() [x]: %v",err) }
	y,err := strconv.ParseFloat(words[2], 32)
	if err!=nil { log.Fatalf("Teleport() [y]: %v",err) }

	ctx.Player.PhysicsState.Pos.X = float32(x)
	ctx.Player.PhysicsState.Pos.Y = float32(y)
	return ""
}

func init() {
	commands["loadmap"] = LoadMap
	commands["savemap"] = SaveMap
	commands["newplanet"] = NewPlanet
	commands["tp"] = Teleport
	commands["help"] = Help
}

func processCommand(line string, ctx CommandContext) string {
	words := strings.Split(line, " ")
	commandName := words[0]
	command, ok := commands[commandName]
	if !ok {
		log.Printf("unrecognized command [%s]", commandName)
		return ""
	}
	return command(line, ctx)
}
