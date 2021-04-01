package cv

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"

	"log"
	"math"
	"os"
)

type Border struct {
	seq_num     int
	border_type int
}

type Point struct {
	y int
	x int
}

type Rect struct {
	pos    Point
	width  int
	height int
}

func NewPoint(xpos, ypos int) Point {
	return Point{
		y: ypos,
		x: xpos,
	}
}

func (p *Point) SetPoint(x, y int) {
	p.y = y
	p.x = x
}

func (p *Point) SamePoint(pnt Point) bool {
	return p.x == pnt.x && p.y == pnt.y
}

type Node struct {
	parent       int
	first_child  int
	next_sibling int
	border       Border
}

func NewNode(p, fc, ns int) Node {
	return Node{
		parent:       p,
		first_child:  fc,
		next_sibling: ns,
	}
}
func (n *Node) reset() {
	n.parent = -1
	n.first_child = -1
	n.next_sibling = -1
}

type Contour struct {
	points       []Point
	boundingRect Rect
	centerOfMass Point
}

func stepCCW(current *Point, pivot Point) {
	if current.x > pivot.x {
		current.SetPoint(pivot.x, pivot.y-1)
	} else if current.x < pivot.x {
		current.SetPoint(pivot.x, pivot.y+1)
	} else if current.y > pivot.y {
		current.SetPoint(pivot.x+1, pivot.y)
	} else if current.y < pivot.y {
		current.SetPoint(pivot.x-1, pivot.y)
	}
}

func stepCW(current *Point, pivot Point) {
	if current.x > pivot.x {
		current.SetPoint(pivot.x, pivot.y+1)
	} else if current.x < pivot.x {
		current.SetPoint(pivot.x, pivot.y-1)
	} else if current.y > pivot.y {
		current.SetPoint(pivot.x-1, pivot.y)
	} else if current.y < pivot.y {
		current.SetPoint(pivot.x+1, pivot.y)
	}
}

func pixelOutOfBounds(p Point, numys, numxs int) bool {
	return (p.x >= numxs || p.y >= numys || p.x < 0 || p.y < 0)
}

func markExamined(mark, center Point, checked *[4]bool) {
	loc := -1
	if mark.x > center.x {
		loc = 0
	} else if mark.x < center.x {
		loc = 2
	} else if mark.y > center.y {
		loc = 1
	} else if mark.y < center.y {
		loc = 3
	}
	if loc == -1 {
		panic("Error: markExamined Failed")
	}

	checked[loc] = true
}

func binarizeImage(src image.Image, threshold int, inverted bool) [][]int {
	numys := src.Bounds().Max.Y
	numxs := src.Bounds().Max.X
	ret := make([][]int, numys)
	for i := range ret {
		ret[i] = make([]int, numxs)
	}
	for r := 0; r < numys; r++ {
		for c := 0; c < numxs; c++ {
			rd, gr, bl, al := src.At(c, r).RGBA()
			gray := 0.3*float32(rd/0xff) + 0.59*float32(gr/0xff) + 0.11*float32(bl/0xff) + 0*float32(al)
			if (int(gray) < threshold && !inverted) || (int(gray) > threshold && inverted) {
				ret[r][c] = 1
			} else {
				ret[r][c] = 0
			}
		}
	}
	return ret
}

func imToGray(src image.Image) [][]int {
	numys := src.Bounds().Max.Y
	numxs := src.Bounds().Max.X
	ret := make([][]int, numys)
	for i := range ret {
		ret[i] = make([]int, numxs)
	}
	for r := 0; r < numys; r++ {
		for c := 0; c < numxs; c++ {
			rd, gr, bl, al := src.At(c, r).RGBA()
			ret[r][c] = int(0.3*float32(rd/0xff) + 0.59*float32(gr/0xff) + 0.11*float32(bl/0xff) + 0*float32(al))
		}
	}
	return ret
}

func followBorder(binaryImg *[][]int, x, y int, p2 *Point, NBD Border, contours *[][]Point) {
	numys := len(*binaryImg)
	numxs := len((*binaryImg)[0])
	current := NewPoint(p2.x, p2.y)
	start := NewPoint(x, y)
	var point_storage []Point

	for {
		if pixelOutOfBounds(current, numys, numxs) {
			stepCW(&current, start)
			if current.SamePoint(*p2) {
				(*binaryImg)[start.y][start.x] = -NBD.seq_num
				point_storage = append(point_storage, start)
				*contours = append(*contours, point_storage)
				return
			}
		} else if (*binaryImg)[current.y][current.x] == 0 {
			stepCW(&current, start)
			if current.SamePoint(*p2) {
				(*binaryImg)[start.y][start.x] = -NBD.seq_num
				point_storage = append(point_storage, start)
				*contours = append(*contours, point_storage)
				return
			}
		} else {
			break
		}
	}
	p1 := current
	p3 := start
	var p4 Point
	*p2 = p1
	var checked [4]bool
	for {
		current = *p2

		for i := 0; i < 4; i++ {
			checked[i] = false
		}

		markExamined(current, p3, &checked)
		stepCCW(&current, p3)

		for {
			if pixelOutOfBounds(current, numys, numxs) {
				markExamined(current, p3, &checked)
				stepCCW(&current, p3)
			} else if (*binaryImg)[current.y][current.x] == 0 {
				markExamined(current, p3, &checked)
				stepCCW(&current, p3)
			} else {
				break
			}
		}
		p4 = current

		if (p3.x+1 >= numxs || (*binaryImg)[p3.y][p3.x+1] == 0) && checked[0] {
			(*binaryImg)[p3.y][p3.x] = -NBD.seq_num
		} else if p3.x+1 < numxs && (*binaryImg)[p3.y][p3.x] == 1 {
			(*binaryImg)[p3.y][p3.x] = NBD.seq_num
		}
		point_storage = append(point_storage, p3)
		if p4.SamePoint(start) && p3.SamePoint(p1) {
			*contours = append(*contours, point_storage)
			return
		}

		*p2 = p3
		p3 = p4
	}
}

func BoundingRect(contour []Point) Rect {
	max := NewPoint(0, 0)
	min := NewPoint(9999, 9999)
	for i := 0; i < len(contour); i++ {
		if contour[i].y > max.y {
			max.y = contour[i].y
		}
		if contour[i].x > max.x {
			max.x = contour[i].x
		}
		if contour[i].y < min.y {
			min.y = contour[i].y
		}
		if contour[i].x < min.x {
			min.x = contour[i].x
		}
	}
	var ret Rect
	ret.pos.x = min.x
	ret.pos.y = min.y
	ret.width = max.x - min.x
	ret.height = max.y - min.y
	return ret
}

func drawLine(p1, p2 Point, img *image.RGBA, col color.Color) {
	var dx, dy, e, slope int
	if p1.x > p2.x {
		p1.x, p1.y, p2.x, p2.y = p2.x, p2.y, p1.x, p1.y
	}

	dx, dy = p2.x-p1.x, p2.y-p1.y
	if dy < 0 {
		dy = -dy
	}

	switch {
	case p1.x == p2.x && p1.y == p2.y:
		img.Set(p1.x, p1.y, col)
	case p1.y == p2.y:
		for ; dx != 0; dx-- {
			img.Set(p1.x, p1.y, col)
			p1.x++
		}
		img.Set(p1.x, p1.y, col)
	case p1.x == p2.x:
		if p1.y > p2.y {
			p1.y, p2.y = p2.y, p1.y
		}
		for ; dy != 0; dy-- {
			img.Set(p1.x, p1.y, col)
			p1.y++
		}
		img.Set(p1.x, p1.y, col)
	case dx == dy:
		if p1.y < p2.y {
			for ; dx != 0; dx-- {
				img.Set(p1.x, p1.y, col)
				p1.x++
				p1.y++
			}
		} else {
			for ; dx != 0; dx-- {
				img.Set(p1.x, p1.y, col)
				p1.x++
				p1.y--
			}
		}
		img.Set(p1.x, p1.y, col)
	case dx > dy:
		if p1.y < p2.y {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img.Set(p1.x, p1.y, col)
				p1.x++
				e -= dy
				if e < 0 {
					p1.y++
					e += slope
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img.Set(p1.x, p1.y, col)
				p1.x++
				e -= dy
				if e < 0 {
					p1.y--
					e += slope
				}
			}
		}
		img.Set(p2.x, p2.y, col)
	default:
		if p1.y < p2.y {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img.Set(p1.x, p1.y, col)
				p1.y++
				e -= dx
				if e < 0 {
					p1.x++
					e += slope
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img.Set(p1.x, p1.y, col)
				p1.y--
				e -= dx
				if e < 0 {
					p1.x++
					e += slope
				}
			}
		}
		img.Set(p2.x, p2.y, col)
	}
}

func drawRectangle(r Rect, img *image.RGBA, col color.Color) {
	x1 := r.pos.x
	y1 := r.pos.y
	x2 := x1 + r.width
	y2 := y1 + r.height
	drawLine(NewPoint(x1, y1), NewPoint(x2, y1), img, col)
	drawLine(NewPoint(x1, y2), NewPoint(x2, y2), img, col)
	drawLine(NewPoint(x1, y1), NewPoint(x1, y2), img, col)
	drawLine(NewPoint(x2, y1), NewPoint(x2, y2), img, col)
}

func (r *Rect) Contains(r2 Rect) bool {
	return r2.pos.x > r.pos.x && r2.pos.y > r.pos.y && r2.pos.x+r2.width < r.pos.x+r.width && r2.pos.y+r2.height < r.pos.y+r.height
}

func (c *Contour) CalculateCenterOfMass() {
	tmpCenter := NewPoint(0, 0)
	for i := range c.points {
		tmpCenter.y += c.points[i].y
		tmpCenter.x += c.points[i].x
	}
	tmpCenter.y /= len(c.points)
	tmpCenter.x /= len(c.points)
	c.centerOfMass = tmpCenter
}

type SpriteSet struct {
	Contours []Contour
	img      image.Image
	binImage [][]int
	epsilon  float32
	drawtype int
}

func (ss *SpriteSet) LoadFile(filename string, threshold int, binanirizeInverted bool) {
	imgFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()
	ss.img, err = jpeg.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}
	ss.binImage = binarizeImage(ss.img, threshold, binanirizeInverted)
}

func (ss *SpriteSet) ProcessContours() {

	var NBD, LNBD Border
	var contours [][]Point

	LNBD.border_type = 1
	NBD.border_type = 1
	NBD.seq_num = 1

	var hierarchy []Node
	temp_node := NewNode(-1, -1, -1)
	temp_node.border = NBD
	hierarchy = append(hierarchy, temp_node)

	var p2 Point
	var border_start_found bool
	for y := 0; y < ss.img.Bounds().Max.Y; y++ {
		LNBD.seq_num = 1
		LNBD.border_type = 1
		for x := 0; x < ss.img.Bounds().Max.X; x++ {
			border_start_found = false
			if (ss.binImage[y][x] == 1 && x-1 < 0) || (ss.binImage[y][x] == 1 && ss.binImage[y][x-1] == 0) {
				NBD.border_type = 2
				NBD.seq_num += 1
				p2.SetPoint(x, y-1)
				border_start_found = true
			} else if x+1 < ss.img.Bounds().Max.X && (ss.binImage[y][x] >= 1 && ss.binImage[y][x+1] == 0) {
				NBD.border_type = 1
				NBD.seq_num += 1
				if ss.binImage[y][x] > 1 {
					LNBD.seq_num = ss.binImage[y][x]
					LNBD.border_type = hierarchy[LNBD.seq_num-1].border.border_type
				}
				p2.SetPoint(x, y+1)
				border_start_found = true
			}

			if border_start_found {
				temp_node.reset()
				if NBD.border_type == LNBD.border_type {
					temp_node.parent = hierarchy[LNBD.seq_num-1].parent
					temp_node.next_sibling = hierarchy[temp_node.parent-1].first_child
					hierarchy[temp_node.parent-1].first_child = NBD.seq_num
					temp_node.border = NBD
					hierarchy = append(hierarchy, temp_node)
				} else {
					if hierarchy[LNBD.seq_num-1].first_child != -1 {
						temp_node.next_sibling = hierarchy[LNBD.seq_num-1].first_child
					}

					temp_node.parent = LNBD.seq_num
					hierarchy[LNBD.seq_num-1].first_child = NBD.seq_num
					temp_node.border = NBD
					hierarchy = append(hierarchy, temp_node)
				}
				followBorder(&ss.binImage, x, y, &p2, NBD, &contours)
			}
			if math.Abs(float64(ss.binImage[y][x])) > 1 {
				LNBD.seq_num = int(math.Abs(float64(ss.binImage[y][x])))
				LNBD.border_type = hierarchy[LNBD.seq_num-1].border.border_type
			}
		}

	}

	var tmpContours []Contour
	for i := range contours {
		var tmpContour Contour
		tmpContour.boundingRect = BoundingRect(contours[i])
		tmpContour.points = contours[i]
		tmpContours = append(tmpContours, tmpContour)
	}
	for i := range tmpContours {
		contains := false
		for j := range tmpContours {
			if tmpContours[j].boundingRect.Contains(tmpContours[i].boundingRect) {
				contains = true
			}
		}
		if !contains {
			ss.Contours = append(ss.Contours, tmpContours[i])
		}
	}

}

var (
	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)

func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(sprite), gl.Ptr(sprite), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}

func (ss *SpriteSet) keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press {
		if k == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
		if k == glfw.KeyUp {
			ss.epsilon += 1
		}
		if k == glfw.KeyDown {
			if ss.epsilon > 0 {
				ss.epsilon -= 1
			}
		}
		if k == glfw.KeySpace {
			ss.epsilon = 0
		}
	}
}

func (ss *SpriteSet) DrawSprite() {
	var tex uint32

	width := ss.img.Bounds().Max.X
	height := ss.img.Bounds().Max.Y

	win := render.NewWindow(height, width, true)
	window := win.Window
	defer glfw.Terminate()
	window.SetKeyCallback(ss.keyCallBack)
	VAO := makeVao()
	program := win.Program
	gl.GenTextures(1, &tex)
	for !window.ShouldClose() {
		img := image.NewRGBA(ss.img.Bounds())
		draw.Draw(img, img.Bounds(), ss.img, image.Pt(0, 0), draw.Src)
		red := color.RGBA{255, 0, 0, 255}
		for i := range ss.Contours {
			if ss.epsilon == 0 {
				drawRectangle(ss.Contours[i].boundingRect, img, red)
			} else {
				tmppnt := RamerDouglasPeucker(ss.Contours[i].points, ss.epsilon)
				for j := 1; j < len(tmppnt); j++ {
					drawLine(tmppnt[j], tmppnt[j-1], img, red)
				}
				drawLine(tmppnt[0], tmppnt[len(tmppnt)-1], img, red)
			}
			ss.Contours[i].CalculateCenterOfMass()

			drawLine(NewPoint(ss.Contours[i].centerOfMass.x, ss.Contours[i].centerOfMass.y-5), NewPoint(ss.Contours[i].centerOfMass.x, ss.Contours[i].centerOfMass.y+5), img, red)
			drawLine(NewPoint(ss.Contours[i].centerOfMass.x-5, ss.Contours[i].centerOfMass.y), NewPoint(ss.Contours[i].centerOfMass.x+5, ss.Contours[i].centerOfMass.y), img, red)
		}

		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.Enable(gl.TEXTURE_2D)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tex)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
		worldTranslate := mgl32.Ident4()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
		projectTransform := mgl32.Ident4()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PerpendicularDistance(pt, lineStart, lineEnd Point) float32 {

	dx := float64(lineEnd.x - lineStart.x)
	dy := float64(lineEnd.y - lineStart.y)
	mag := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if mag > 0 {
		dx /= mag
		dy /= mag
	}

	pvx := float64(pt.x - lineStart.x)
	pvy := float64(pt.y - lineStart.y)

	pvdot := dx*pvx + dy*pvy

	dsx := pvdot * dx
	dsy := pvdot * dy

	ax := pvx - dsx
	ay := pvy - dsy

	return float32(math.Sqrt(math.Pow(ax, 2) + math.Pow(ay, 2)))
}

func RamerDouglasPeucker(pointList []Point, epsilon float32) []Point {
	var ret []Point

	if len(pointList) < 2 {
		panic("Not enough points to simplify")
	}

	dmax := float32(0)
	index := 0
	end := len(pointList) - 1
	for i := range pointList {
		d := PerpendicularDistance(pointList[i], pointList[0], pointList[end])
		if d > dmax {
			index = i
			dmax = d
		}
	}

	if dmax > epsilon {

		firstLine := pointList[0 : index+1]
		lastLine := pointList[index:]

		recResults1 := RamerDouglasPeucker(firstLine, epsilon)
		recResults2 := RamerDouglasPeucker(lastLine, epsilon)

		ret = append(recResults1, recResults2...)
		if len(ret) < 2 {
			panic("Problem assembling output")
		}
	} else {
		ret = nil
		ret = append(ret, pointList[0])
		ret = append(ret, pointList[end])
	}
	return ret
}
