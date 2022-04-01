package main

// import (
// 	"fmt"
// 	_ "image/png"
// 	"log"
// 	"runtime"

// 	//"github.com/go-gl/gl/v2.1/gl"
// 	"github.com/go-gl/gl/v4.1-core/gl"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// 	"github.com/go-gl/mathgl/mgl32"
// 	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBufferDY"
// 	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/UI_Injector"
// 	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBufferDY"
// 	n "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/character"
// 	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/renderer"
// 	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/shader"
// 	vertexArray "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArrayDY"
// 	vertexbufferLayout "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayoutDY"
// 	"github.com/skycoin/cx-game/engine/spine"
// 	"github.com/skycoin/cx-game/render"
// 	"github.com/skycoin/cx-game/test/spine-animation/animation"
// 	"github.com/skycoin/cx-game/world"
// )

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

// var (
// 	DrawCollisionBoxes = false
// 	fps                *render.Fps
// 	FPS                int
// )

// var CurrentPlanet *world.Planet

// const (
// 	width  = 960
// 	height = 540
// )

// var (
// 	characters     []*n.Character
// 	character      *n.Character
// 	characterIndex int
// )

// var (
// 	positions = []float32{
// 		// // X    Y     U     V   texIndex
// 		// -50.0, -50.0, 0.0, 0.0, 0.0, //1
// 		// 50.0, -50.0, 1.0, 0.0, 0.0, //2
// 		// 50.0, 50.0, 1.0, 1.0, 0.0, // 3
// 		// -50.0, 50.0, 0.0, 1.0, 0.0, //4

// 		// 50.0, -50.0, 0.0, 0.0, 1.0, //1
// 		// 150.0, -50.0, 1.0, 0.0, 1.0, //2
// 		// 150.0, 50.0, 1.0, 1.0, 1.0, // 3
// 		// 50.0, 50.0, 0.0, 1.0, 1.0, //4

// 	}

// 	// indices = []uint32{
// 	// 	0, 1, 2,
// 	// 	2, 3, 0,

// 	// 	4, 5, 6,
// 	// 	6, 7, 4,
// 	// }
// 	samplers = [31]int32{
// 		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
// 	}
// )

// type Vertex struct {
// 	Position  mgl32.Vec2
// 	TexCoords mgl32.Vec2
// 	TexID     float32
// }

// func CreateQuad(x, y, textureID float32) []float32 {

// 	var size float32 = 100.0
// 	var pos []float32
// 	v0 := Vertex{}
// 	v0.Position = mgl32.Vec2{x, y}
// 	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
// 	v0.TexID = textureID
// 	pos = append(pos, v0.Position.X(), v0.Position.Y(), v0.TexCoords.X(), v0.TexCoords.Y(), v0.TexID)

// 	v1 := Vertex{}
// 	v1.Position = mgl32.Vec2{x + size, y}
// 	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
// 	v1.TexID = textureID
// 	pos = append(pos, v1.Position.X(), v1.Position.Y(), v1.TexCoords.X(), v1.TexCoords.Y(), v1.TexID)

// 	v2 := Vertex{}
// 	v2.Position = mgl32.Vec2{x + size, y + size}
// 	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
// 	v2.TexID = textureID
// 	pos = append(pos, v2.Position.X(), v2.Position.Y(), v2.TexCoords.X(), v2.TexCoords.Y(), v2.TexID)

// 	v3 := Vertex{}
// 	v3.Position = mgl32.Vec2{x, y + size}
// 	v3.TexCoords = mgl32.Vec2{0.0, 1.0}
// 	v3.TexID = textureID
// 	pos = append(pos, v3.Position.X(), v3.Position.Y(), v3.TexCoords.X(), v3.TexCoords.Y(), v3.TexID)

// 	//fmt.Printf("numbers=%v\n", pos)
// 	return pos

// }

// func initGlfw() *glfw.Window {
// 	if err := glfw.Init(); err != nil {
// 		panic(err)
// 	}
// 	glfw.WindowHint(glfw.Resizable, glfw.False)
// 	glfw.WindowHint(glfw.ContextVersionMajor, 4)
// 	glfw.WindowHint(glfw.ContextVersionMinor, 0)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
// 	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

// 	window, err := glfw.CreateWindow(width, height, "batch", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	window.MakeContextCurrent()
// 	glfw.SwapInterval(1)

// 	return window
// }

// // initOpenGL initializes OpenGL and returns an intiialized program.
// func initOpenGL() {
// 	// // Initialize Glow
// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	version := gl.GoStr(gl.GetString(gl.VERSION))
// 	fmt.Println("OpenGL version", version)
// }

// var ib *indexbuffer.IndexBuffer
// var vb *vertexbuffer.VertexBuffer
// var va *vertexArray.VertexArray
// var vbl *vertexbufferLayout.VertexbufferLayout

// func codeT() {
// 	{

// 		// var m_QuadVA uint32
// 		// var m_QuadVB uint32
// 		// var m_QuadIB uint32

// 		// gl.GenVertexArrays(1, &m_QuadVA)
// 		// gl.BindVertexArray(m_QuadVA)

// 		// gl.GenBuffers(1, &m_QuadVB)
// 		// gl.BindBuffer(gl.ARRAY_BUFFER, m_QuadVB)
// 		// gl.BufferData(gl.ARRAY_BUFFER, 4*len(positions), gl.Ptr(positions), gl.STATIC_DRAW)

// 		// gl.EnableVertexArrayAttrib(m_QuadVB, 0)
// 		// gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

// 		// gl.GenBuffers(1, &m_QuadIB)
// 		// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m_QuadIB)
// 		// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

// 		// tex := Texture.SetUpTexture("./cat.png")

// 		// tex.Bind(0)
// 		// renderer.GLClearError()
// 		// shader.SetUniForm1i("u_Texture", 0)

// 	}
// }
// func main() {

// 	runtime.LockOSThread()

// 	window := initGlfw()
// 	defer glfw.Terminate()
// 	initOpenGL()

// 	fps = render.NewFps(false)

// 	var UI *UI_Injector.UI_Injector

// 	UI = UI_Injector.SetUpUI()
// 	go UI.ListenForChanges()
// 	var objectAdjustment = UI
// 	objectAdjustment.Object[0].CharacterAnimButton = -1
// 	var proj mgl32.Mat4 = mgl32.Ortho(0.0, 960.0, 0.0, 540.0, -1.0, 1.0)
// 	var view mgl32.Mat4 = mgl32.Translate3D(0.0, 0.0, 0.0)
// 	gl.Enable(gl.BLEND)
// 	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
// 	shader := shader.SetupShader("./../../assets/shader/spine/basic.shader")
// 	shader.Bind()
// 	shader.SetUniForm4f("u_Color", 0.8, 0.3, 0.8, 1.0)

// 	for _, loc := range animation.LoadList("./../../test/spine-animation/animation") {
// 		character, err := n.LoadCharacter(loc)
// 		if err != nil {
// 			log.Println(loc.Name, err)
// 			break
// 		}

// 		for _, skin := range character.Skeleton.Data.Skins {
// 			for _, att := range skin.Attachments {
// 				if _, ismesh := att.(*spine.MeshAttachment); ismesh {
// 					log.Println(loc.Name, "Unsupported")
// 					break
// 				}
// 			}
// 		}

// 		characters = append(characters, character)
// 	}

// 	character = characters[0]

// 	const MaxQuadCount = 1000
// 	const MaxVertexCount = MaxQuadCount * 4
// 	const MaxIndexCount = MaxQuadCount * 6

// 	var indices = make([]uint32, MaxIndexCount)
// 	var offset uint32 = 0
// 	for i := 0; i < MaxIndexCount; i += 6 {
// 		indices[i+0] = 0 + offset
// 		indices[i+1] = 1 + offset
// 		indices[i+2] = 2 + offset

// 		indices[i+3] = 2 + offset
// 		indices[i+4] = 3 + offset
// 		indices[i+5] = 0 + offset

// 		offset += 4
// 	}

// 	//setup vertex array
// 	va = vertexArray.SetUpVertxArray()
// 	// setup and run vertex buffer
// 	vb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
// 	//setup vertex layout
// 	vbl = &vertexbufferLayout.VertexbufferLayout{}
// 	//add vertex buffer to vertex bufferlayout

// 	vbl.Push(gl.FLOAT, 2)
// 	vbl.Push(gl.FLOAT, 2)
// 	vbl.Push(gl.FLOAT, 1)
// 	va.AddBuffer(vb, vbl)

// 	// setup and run index buffer
// 	ib = indexbuffer.RunIndexBuffer(indices, len(indices))

// 	//	tex := Texture.SetUpTexture("./cat.png")
// 	//	tex1 := Texture.SetUpTexture("./sprite.png")

// 	//	tex.Bind(0)
// 	//	tex1.Bind(1)
// 	renderer.GLClearError()
// 	//	gl.Uniform1iv(location int32, count int32, value *int32)
// 	//loc := gl.GetUniformLocation(shader.M_renderID, gl.Str("u_Texture\x00"))
// 	//gl.Uniform1iv(loc, 2, &samplers[0])
// 	shader.SetUniform1iv("u_Texture", int32(len(samplers)), &samplers[0])
// 	//shader.SetUniForm1i("u_Texture", 0)
// 	renderer.GLCheckError()
// 	// va.Unbind()
// 	// vb.Unbind()
// 	// ib.Unbind()
// 	// shader.UnBind()

// 	var translationA = mgl32.Translate3D(objectAdjustment.Object[0].X, objectAdjustment.Object[0].Y, objectAdjustment.Object[0].Z)

// 	// var translationB = mgl32.Translate3D(400+objectAdjustment.Object[1].X, 200+objectAdjustment.Object[1].Y, objectAdjustment.Object[1].Z)

// 	render := renderer.SetupRender()

// 	//var test1 *TestClearColor.TestClearColor

// 	//test1 := TestClearColor.SetUpTestClearColor()

// 	var r float32 = 0.0
// 	var increment float32 = 0.5

// 	for !window.ShouldClose() {
// 		//render.Clear()
// 		// test1.OnUpdate(0.0)
// 		// test1.OnRender()
// 		// test1.M_ClearColor = [4]float32{objectAdjustment.Object[2].X, objectAdjustment.Object[2].Y, objectAdjustment.Object[2].Z}
// 		// fmt.Println(test1.M_ClearColor)
// 		{
// 			fps.Tick()
// 			if objectAdjustment.Object[4].CharacterAnimButton == 0 {
// 				characterIndex = (characterIndex + len(characters) - 1) % len(characters)
// 				character = characters[characterIndex]
// 			}
// 			if objectAdjustment.Object[4].CharacterAnimButton == 1 {
// 				characterIndex = (characterIndex + len(characters) + 1) % len(characters)
// 				character = characters[characterIndex]
// 			}

// 			if objectAdjustment.Object[4].CharacterAnimButton == 5 {
// 				character.NextAnimation(-1)
// 			}
// 			if objectAdjustment.Object[4].CharacterAnimButton == 4 {
// 				character.NextAnimation(1)
// 			}

// 			if objectAdjustment.Object[4].CharacterAnimButton == 6 {
// 				character.NextSkin(-1)
// 			}
// 			if objectAdjustment.Object[1].CharacterAnimButton == 7 {
// 				character.NextSkin(1)
// 			}

// 			objectAdjustment.Object[4].CharacterAnimButton = -1
// 			character.Update(1/float64(fps.GetCurFps()), width/2, height/2)

// 			// if ebiten.IsRunningSlowly() {
// 			// 	return nil
// 			// }

// 			//	screen.Clear()

// 			//	character.Draw(screen)

// 			//	ebitenutil.DebugPrint(screen, character.Description())
// 		}

// 		//************** Code Before Tests scanes ***********************//
// 		{

// 			render.Clear()
// 			{
// 				var positions []float32
// 				// q0 := CreateQuad(objectAdjustment.Object[1].X+(-50.0), objectAdjustment.Object[1].Y+(-50.0), 0.0)
// 				// q1 := CreateQuad(150.0, -50.0, 1.0)
// 				// q2 := CreateQuad(250.0, -50.0, 1.0)
// 				// q3 := CreateQuad(350.0, -50.0, 1.0)
// 				// positions = append(positions, q0...)
// 				// positions = append(positions, q1...)
// 				// positions = append(positions, q2...)
// 				// positions = append(positions, q3...)
// 				positions = character.Draw()
// 				//fmt.Println(q0)

// 				vb.Bind()
// 				//gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(positions)*5*4, gl.Ptr(positions))
// 				vb.BufferSubData(positions)
// 			}
// 			//tex.Bind(0)
// 			//tex1.Bind(1)
// 			//fmt.Println(translation)
// 			// shader.SetUniForm4f("u_Color", r, 0.3, 0.8, 1.0)

// 			{
// 				var model mgl32.Mat4 = translationA.Add(mgl32.Translate3D(objectAdjustment.Object[0].X*10, objectAdjustment.Object[0].Y*10, objectAdjustment.Object[0].Z*10))
// 				mvp := proj.Mul4(view).Mul4(model)
// 				shader.Bind()
// 				shader.SetUniFormMat4f("u_MVP", mvp)
// 				render.DrawDY(va, ib, shader)
// 			}

// 			// {
// 			// 	var model mgl32.Mat4 = translationB.Add(mgl32.Translate3D(objectAdjustment.Object[1].X*10, objectAdjustment.Object[1].Y*10, objectAdjustment.Object[1].Z*10))
// 			// 	mvp := proj.Mul4(view).Mul4(model)
// 			// 	shader.Bind()
// 			// 	shader.SetUniFormMat4f("u_MVP", mvp)
// 			// 	render.DrawDY(va, ib, shader)
// 			// }

// 			if r > 1.9 {
// 				increment = -0.05
// 			} else if r < 0.0 {
// 				increment = 0.05
// 			}
// 			r = r + increment

// 		}

// 		//****************************************************************//
// 		glfw.PollEvents()
// 		window.SwapBuffers()
// 	}
// 	shader.DeleteShader()

// }
