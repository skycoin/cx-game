package main

import (
	"flag"
	"os"
	"log"
	"runtime/pprof"

	"github.com/skycoin/cx-game/z_game/game"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	game.Run()
}
