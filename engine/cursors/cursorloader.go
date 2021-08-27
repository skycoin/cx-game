package input

type CursorLoader struct {
}

// ex.
// upLeft := image.Point{0, 0}
// lowRight := image.Point{16, 16}
// img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
// cyan := color.RGBA{100, 200, 200, 0xff}
// for x := 0; x < 16; x++ {
// 	for y := 0; y < 16; y++ {
// 		switch {
// 		case x < 16/2 && y < 16/2:
// 			img.Set(x, y, cyan)
// 		case x >= 16/2 && y >= 16/2:
// 			img.Set(x, y, color.White)
// 		default:
// 		}
// 	}
// }
// cursor := glfw.CreateCursor(img, 0, 0)
// window.SetCursor(cursor)
