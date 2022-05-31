package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
	c "github.com/skycoin/cx-game/test/spine-animation/character"
)

const (
	width  = 500
	height = 500
)

var (
	characters     []*c.Character
	character      *c.Character
	characterIndex int
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	test()
	last := time.Now()
	for !window.ShouldClose() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		// center := cx.Vec{X:100,Y:100}
		// center.Y = 100
		character.Update(dt, 250, 250)
		draw(window, program)
	}
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

func test() {

	for _, loc := range animation.LoadList("./animation") {
		character, err := c.LoadCharacter(loc)
		if err != nil {
			log.Println(loc.Name, err)
			continue
		}

		for _, skin := range character.Skeleton.Data.Skins {
			for _, att := range skin.Attachments {
				if _, ismesh := att.(*spine.MeshAttachment); ismesh {
					log.Println(loc.Name, "Unsupported")
					//	continue skip
				}
			}
		}

		characters = append(characters, character)
	}

	// var character *c.Character
	characterIndex := 0
	character = characters[characterIndex]

	log.Println("%v", character)

}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "spine Demo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	glfw.PollEvents()
	window.SwapBuffers()
}
