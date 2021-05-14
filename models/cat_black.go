package models

import (
	"fmt"
	"image"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type CatBlack struct {
	RGBA                 *image.RGBA
	Size                 image.Point
	Walk                 func()
	Sit                  func()
	SitStop              func()
	StartRunning         func()
	Running              func()
	width                float32
	height               float32
	X                    float32
	Y                    float32
	XVelocity, YVelocity float32
	movSpeed             float32
	jumpSpeed            float32
	SpriteSheetId        int
}

const (
	walkSprite         = 0
	sitSprite          = 1
	startRunningSprite = 2
	runningSprite      = 3
)

var stopPlay chan bool

// private method
func play(action int, fcount int, lwindow *glfw.Window, lspriteSheetId int) {
	stopPlay = make(chan bool)
	j := 0
	for {
		select {
		default:
			time.Sleep(100 * time.Millisecond)
			spriteloader.LoadSprite(lspriteSheetId, "blackcat", action, j)
			spriteId := spriteloader.GetSpriteIdByName("blackcat")
			fmt.Println("spriteId. ", spriteId, " j. ", j)
			// gl.ClearColor(1, 1, 1, 1)
			// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			// spriteloader.DrawSpriteQuad(0, 0, 2, 1, spriteId)
			lwindow.SwapBuffers()
			j++
			if j == fcount {
				j = 0
			}
		case <-stopPlay:
			return
		}
	}
}

func stop() {
	fmt.Println("Stop.")
	stopPlay <- true
	// sp := <-stopPlay
	// fmt.Println("channel. ", sp)
}

func NewCatBlack(lwin *render.Window, lwindow *glfw.Window) *CatBlack {
	spriteloader.InitSpriteloader(lwin)
	lspriteSheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/blackcat_sprite.png", 13, 4)
	catBlack := CatBlack{
		width:         2,
		height:        1,
		movSpeed:      0.05,
		jumpSpeed:     0.2,
		SpriteSheetId: lspriteSheetId,
		Walk: func() {
			play(walkSprite, 11, lwindow, lspriteSheetId)
		},
		Sit: func() {
			play(sitSprite, 5, lwindow, lspriteSheetId)
		},
		SitStop: func() {
			stop()
		},
		StartRunning: func() {
			play(startRunningSprite, 11, lwindow, lspriteSheetId)
		},
		Running: func() {
			play(runningSprite, 13, lwindow, lspriteSheetId)
		},
	}

	return &catBlack
}
