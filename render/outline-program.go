package render

var outlineProgram Program

func InitOutlineProgram() {
	outlineProgram = CompileProgram(
		"./assets/shader/outline.vert", "./assets/shader/outline.frag" )
}
