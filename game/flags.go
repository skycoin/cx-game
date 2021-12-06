package game

import (
	"flag"

	"github.com/skycoin/cx-game/constants"
)

var tmxPtr = flag.String("map", constants.DefaultTmxPath, "tmx path")

type StartupFlags struct {
	TmxPath string
}

func ParseStartupFlags() StartupFlags {
	return StartupFlags {
		TmxPath: *tmxPtr,
	}
}
