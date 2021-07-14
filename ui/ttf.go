package ui

import (
	"log"


	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/ui/glfont"
)

const numFonts = 10
const fontSize = 24

type FontID uint32
var fonts [numFonts]*glfont.Font

const DefaultFontID = FontID(0)

func LoadFont(id FontID, fname string, scale int32) {
	font,err := glfont.LoadFont(fname,18,100,100)
	if err!=nil { log.Fatal(err) }

	fonts[id] = font
}

func initFonts() {
	LoadFont(DefaultFontID, "./assets/font/GravityBold8.ttf", 24)
}

func drawString(ctx render.Context, str string) {
	
}
