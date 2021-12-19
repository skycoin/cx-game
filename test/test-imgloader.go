package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	sl "github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called
	// from the main thread.
	runtime.LockOSThread()
}

var il *sl.ImgLoader
var spriteId int = -1

//test loader callback function
func loaderCb(il *sl.ImgLoader) {
	fmt.Printf("loader cb \n")
	id := sl.AddSpriteSheet("assets/starfield/stars/planets.png", il)
	fmt.Printf("AddSpriteSheet id %v\n", id)
	if id < 0 {
		return
	}

	sl.LoadSprite(id, "star", 2, 1)
	spriteId = sl.GetSpriteIdByName("star")
	fmt.Printf("sprite id %v\n", spriteId)
}

func getPng(dir string) []string {
	var ret []string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".png") {
				return nil
			}
			ret = append(ret, path)
			return nil
		})

	if err != nil {
		log.Printf("error dir %s\n", dir)
	}

	return ret
}

func main() {
	imgList := getPng("./assets")

	il = sl.NewImgLoader(imgList, loaderCb)
	il.Run()

	log.Print("running test-imageloader")
	win := render.NewWindow(640, 480, true)
	window := win.Window
	defer glfw.Terminate()
	sl.InitSpriteloader(&win)

	for !window.ShouldClose() {
		if il.NeedUpdate() {
			il.Update()
		}

		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// if load finish ,the img will render
		if spriteId >= 0 {
			sl.DrawSpriteQuad(0, 0, 2, 2, spriteId)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
