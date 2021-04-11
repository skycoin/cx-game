package cv

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"

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

type Joint struct {
	parentIndex int
	pivot       Point
	parentPivot Point
}

type Contour struct {
	points        []Point
	wrapPoints    []Point
	boundingRect  Rect
	centerOfMass  Point
	name          string
	img           image.RGBA
	selectedIndex int
	joint         Joint
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
	ret.pos.x = min.x - 1
	ret.pos.y = min.y - 1
	ret.width = max.x - min.x + 2
	ret.height = max.y - min.y + 2
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

func drawPoint(p Point, img *image.RGBA, col color.Color) {
	x1 := p.x - 5
	y1 := p.y - 5
	x2 := p.x + 5
	y2 := p.y + 5
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			img.Set(x, y, col)
		}
	}
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

type Queue struct {
	elements chan Point
}

func NewQueue(size int) *Queue {
	return &Queue{
		elements: make(chan Point, size),
	}
}

func (queue *Queue) Push(element Point) {
	select {
	case queue.elements <- element:
	default:
		panic("Queue full")
	}
}

func (queue *Queue) Pop(p *Point) bool {
	select {
	case *p = <-queue.elements:
		return true
	default:
		return false
	}
}

func MaskFromContour(contour []Point, width, height int, offset Point) [][]int {
	var ret [][]int = make([][]int, width)

	for i := 0; i < width; i++ {
		ret[i] = make([]int, height)
	}
	for i := range contour {
		ret[contour[i].x-offset.x][contour[i].y-offset.y] = 1
	}

	//	var stack []Point
	var stack = NewQueue(100)
	var spanAbove bool
	var spanBelow bool
	stack.Push(NewPoint(0, 0))
	//	stack.Push(NewPoint(width-1, 0))
	//	stack.Push(NewPoint(width-1, height-2))
	//	stack.Push(NewPoint(0, height-2))
	var p Point
	for stack.Pop(&p) {
		x1 := p.x
		for x1 >= 0 && ret[x1][p.y] == 0 {
			x1--
		}
		x1++
		spanAbove = false
		spanBelow = false
		for x1 < width && ret[x1][p.y] == 0 {
			ret[x1][p.y] = 1
			if !spanAbove && p.y > 0 && ret[x1][p.y] == 0 {
				stack.Push(NewPoint(x1, p.y-1))
				spanAbove = true
			} else if spanAbove && p.y > 0 && ret[x1][p.y-1] != 0 {
				spanAbove = false
			}
			if !spanBelow && p.y < height-1 && ret[x1][p.y+1] == 0 {
				stack.Push(NewPoint(x1, p.y+1))
				spanBelow = true
			} else if spanBelow && p.y < height-1 && ret[x1][p.y+1] != 0 {
				spanBelow = false
			}
			x1++
		}
	}
	return ret
}

const (
	JointTool = iota
	AddPointTool
	RemovePointTool
	SelectTool

//    West
)

type SpriteSet struct {
	Contours []Contour
	img      image.Image
	binImage [][]int
	mousePos Point

	selectedContour int
	tool            int
	drawType        int
	canDrag         bool
	dragOffset      Point
	tmpJoint        Joint
}

func (ss *SpriteSet) LoadFile(filename string, threshold int, binanirizeInverted bool) {
	imgFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()
	ss.img, err = png.Decode(imgFile)
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
		tmpContour.selectedIndex = -1
		tmpContours = append(tmpContours, tmpContour)
	}
	for i := range tmpContours {
		contains := false
		for j := range tmpContours {
			if tmpContours[j].boundingRect.Contains(tmpContours[i].boundingRect) {
				contains = true
			}
		}
		if !contains && tmpContours[i].boundingRect.height > 0 && tmpContours[i].boundingRect.width > 0 {
			tmpContours[i].ConvexHull()
			ss.Contours = append(ss.Contours, tmpContours[i])
		}
	}
	for i := range ss.Contours {
		imrect := image.Rect(0, 0, ss.Contours[i].boundingRect.width, ss.Contours[i].boundingRect.height)
		brect := image.Rect(ss.Contours[i].boundingRect.pos.x, ss.Contours[i].boundingRect.pos.y, ss.Contours[i].boundingRect.pos.x+ss.Contours[i].boundingRect.width, ss.Contours[i].boundingRect.pos.y+ss.Contours[i].boundingRect.height)
		ss.Contours[i].img = *image.NewRGBA(imrect)
		draw.Draw(&ss.Contours[i].img, imrect, ss.img, brect.Min, draw.Src)
		mask := MaskFromContour(ss.Contours[i].points, ss.Contours[i].boundingRect.width+1, ss.Contours[i].boundingRect.height+1, ss.Contours[i].boundingRect.pos)
		for y := 0; y < ss.Contours[i].boundingRect.height; y++ {
			for x := 0; x < ss.Contours[i].boundingRect.width; x++ {
				if mask[x][y] == 1 {
					ss.Contours[i].img.Set(x, y, color.Transparent)
				}
			}
		}
		ss.Contours[i].joint.parentIndex = -1
		fn := fmt.Sprintf("c:\\test\\%d_test.png", i+1)
		imgFile, err := os.Create(fn)
		if err != nil {
			log.Fatalln(err)
		}
		defer imgFile.Close()
		err = png.Encode(imgFile, &ss.Contours[i].img)
		if err != nil {
			log.Fatalln(err)
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

func euclidianDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.x-p2.x), 2) + math.Pow(float64(p1.y-p2.y), 2))
}

func (ss *SpriteSet) cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	ss.mousePos.x = int(xpos)
	ss.mousePos.y = int(ypos)
	if ss.canDrag {
		ss.Contours[ss.selectedContour].wrapPoints[ss.Contours[ss.selectedContour].selectedIndex] = NewPoint(ss.mousePos.x+ss.dragOffset.x, ss.mousePos.y+ss.dragOffset.y)
	}
}

func (ss *SpriteSet) addPoint(p Point) {
	ss.Contours[ss.selectedContour].wrapPoints = append(ss.Contours[ss.selectedContour].wrapPoints, p)
	sort.Slice(ss.Contours[ss.selectedContour].wrapPoints, func(i int, j int) bool {
		tmpvp1 := NewPoint(ss.Contours[ss.selectedContour].wrapPoints[i].x-ss.Contours[ss.selectedContour].centerOfMass.x, ss.Contours[ss.selectedContour].wrapPoints[i].y-ss.Contours[ss.selectedContour].centerOfMass.y)
		tmpan1 := math.Atan2(float64(tmpvp1.x), float64(tmpvp1.y))
		tmpvp2 := NewPoint(ss.Contours[ss.selectedContour].wrapPoints[j].x-ss.Contours[ss.selectedContour].centerOfMass.x, ss.Contours[ss.selectedContour].wrapPoints[j].y-ss.Contours[ss.selectedContour].centerOfMass.y)
		tmpan2 := math.Atan2(float64(tmpvp2.x), float64(tmpvp2.y))
		return tmpan1 < tmpan2
	})
}

func (r *Rect) isPointInside(p Point) bool {
	return p.x >= r.pos.x && p.y >= r.pos.y && p.x <= r.pos.x+r.width && p.y <= r.pos.y+r.height
}

func (ss *SpriteSet) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if ss.tool == SelectTool {
			ss.selectClosestPoint(ss.mousePos)
		}
		if ss.tool == AddPointTool {
			ss.addPoint(ss.mousePos)
		}
		if ss.tool == JointTool {
			ss.AddJoint()
		}
	}
	if action == glfw.Release {
		ss.canDrag = false
	}
}

func (ss *SpriteSet) AddJoint() {
	ss.selectContour(ss.mousePos)
	if ss.selectedContour != -1 {
		if ss.tmpJoint.parentIndex == -1 {
			ss.tmpJoint.parentIndex = ss.selectedContour
			ss.tmpJoint.parentPivot = NewPoint(ss.mousePos.x-ss.Contours[ss.selectedContour].boundingRect.pos.x, ss.mousePos.y-ss.Contours[ss.selectedContour].boundingRect.pos.y)
		} else {
			ss.tmpJoint.pivot = NewPoint(ss.mousePos.x-ss.Contours[ss.selectedContour].boundingRect.pos.x, ss.mousePos.y-ss.Contours[ss.selectedContour].boundingRect.pos.y)
			ss.Contours[ss.selectedContour].joint = ss.tmpJoint
			ss.tmpJoint.parentIndex = -1
		}
	}
}

func (ss *SpriteSet) keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press {
		if k == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
		if k == glfw.KeyJ {
			w.SetTitle("Joint tool. Create pivot point on parent")
			ss.tool = JointTool
			ss.drawType = 2
			ss.tmpJoint.parentIndex = -1
		}
		if k == glfw.KeyA {
			w.SetTitle("Add point tool")
			ss.tool = AddPointTool
			ss.drawType = 1
		}
		if k == glfw.KeyS {
			w.SetTitle("Select tool")
			ss.tool = SelectTool
			ss.drawType = 1
		}
		if k == glfw.KeyD {
			ss.drawType = 1 - ss.drawType
		}
		if k == glfw.KeyE {
			ss.eraseSelectedPoint()
			ss.drawType = 1
		}
	}
}

func orient(a, b, c Point) int {
	v := (b.y-a.y)*(c.x-b.x) - (b.x-a.x)*(c.y-b.y)
	if v == 0 {
		return 0
	} else if v > 0 {
		return 1
	} else {
		return 2
	}
}

func (c *Contour) ConvexHull() {
	m := len(c.points)
	if m < 3 {
		return
	}
	l := 0

	var next []int
	for i := 0; i < m; i++ {
		next = append(next, -1)
	}

	for i := 1; i < m; i++ {
		if c.points[i].x < c.points[l].x {
			l = i
		}
	}
	p := l
	var q int
	for {
		q = (p + 1) % m
		for i := 0; i < m; i++ {
			if orient(c.points[p], c.points[i], c.points[q]) == 2 {
				q = i
			}
		}
		next[p] = q
		p = q
		if p == l {
			break
		}
	}
	for i := 0; i < m; i++ {
		if next[i] != -1 {
			c.wrapPoints = append(c.wrapPoints, c.points[i])
		}
	}

}

func drawContour(points []Point, img *image.RGBA, col color.Color, selectedPoint int, selectedPointColor color.Color) {
	if len(points) > 1 {
		for j := 1; j < len(points); j++ {
			drawLine(points[j], points[j-1], img, col)
		}
		drawLine(points[0], points[len(points)-1], img, col)
		for j := range points {
			if j == selectedPoint {
				drawPoint(points[j], img, selectedPointColor)
			} else {
				drawPoint(points[j], img, col)
			}
		}
	}
}

/*func drawContourDeb(points []Point, img *image.RGBA, col color.Color, selectedPoint int, selectedPointColor color.Color, f int, s int) {
	if len(points) > 1 {
		for j := 1; j < len(points); j++ {
			drawLine(points[j], points[j-1], img, col)
		}
		d := color.RGBA{255, 255, 0, 255}
		drawLine(points[0], points[len(points)-1], img, col)
		for j := range points {
			if j == selectedPoint {
				drawPoint(points[j], img, selectedPointColor)
			} else if j == f {
				drawPoint(points[j], img, d)
			} else if j == s {
				drawPoint(points[j], img, d)
			} else {
				drawPoint(points[j], img, col)
			}
		}
	}
}
*/
func (ss *SpriteSet) eraseSelectedPoint() {
	if ss.Contours[ss.selectedContour].selectedIndex > -1 {
		copy(ss.Contours[ss.selectedContour].wrapPoints[ss.Contours[ss.selectedContour].selectedIndex:], ss.Contours[ss.selectedContour].wrapPoints[ss.Contours[ss.selectedContour].selectedIndex+1:])
		ss.Contours[ss.selectedContour].wrapPoints[len(ss.Contours[ss.selectedContour].wrapPoints)-1] = NewPoint(0, 0)
		ss.Contours[ss.selectedContour].wrapPoints = ss.Contours[ss.selectedContour].wrapPoints[:len(ss.Contours[ss.selectedContour].wrapPoints)-1]
	}
}

func (ss *SpriteSet) selectContour(p Point) {
	ss.selectedContour = -1
	for i := range ss.Contours {
		if ss.Contours[i].boundingRect.isPointInside(p) {
			ss.selectedContour = i
		}
	}
}

func (ss *SpriteSet) selectClosestPoint(p Point) {
	selectedPoint := -1
	for i := range ss.Contours {
		//		if ss.Contours[i].boundingRect.isPointInside(p) {
		for j := range ss.Contours[i].wrapPoints {
			dist := euclidianDistance(ss.Contours[i].wrapPoints[j], p)
			if dist < 7 {
				selectedPoint = j
				ss.selectedContour = i
			}
		}
		//		}
	}
	if selectedPoint > -1 {
		ss.Contours[ss.selectedContour].selectedIndex = selectedPoint
		ss.dragOffset = NewPoint(ss.Contours[ss.selectedContour].wrapPoints[selectedPoint].x-ss.mousePos.x, ss.Contours[ss.selectedContour].wrapPoints[selectedPoint].y-ss.mousePos.y)
		ss.canDrag = true
	} else {
		selected := -1
		for c := range ss.Contours {
			if ss.Contours[c].boundingRect.isPointInside(p) {
				selected = c
			}
		}
		if selected > -1 {
			ss.selectedContour = selected
			ss.Contours[ss.selectedContour].selectedIndex = -1
		}
	}
}

func isAllDrawn(drawnPoints []Point) bool {
	for i := range drawnPoints {
		if drawnPoints[i].x == -1 && drawnPoints[i].y == -1 {
			return false
		}
	}
	return true
}

func (ss *SpriteSet) DrawJoints(img *image.RGBA) {
	drawnPositions := make([]Point, len(ss.Contours))
	width := ss.img.Bounds().Max.X
	//height = ss.img.Bounds().Max.Y
	for i := range drawnPositions {
		drawnPositions[i] = NewPoint(-1, -1)
	}
	//first we draw all independent contours
	for i := range ss.Contours {
		if ss.Contours[i].joint.parentIndex == -1 {
			p := NewPoint(ss.Contours[i].boundingRect.pos.x+width, ss.Contours[i].boundingRect.pos.y)
			imrect := image.Rect(p.x, p.y, p.x+ss.Contours[i].img.Bounds().Max.X, p.y+ss.Contours[i].img.Bounds().Max.Y)
			draw.Draw(img, imrect, &ss.Contours[i].img, image.Pt(0, 0), draw.Over)
			drawnPositions[i] = p
		}
	}
	for !isAllDrawn(drawnPositions) {
		for i := range ss.Contours {
			if ss.Contours[i].joint.parentIndex != -1 {
				if drawnPositions[ss.Contours[i].joint.parentIndex].x != -1 && drawnPositions[i].x == -1 {
					x := drawnPositions[ss.Contours[i].joint.parentIndex].x + ss.Contours[i].joint.parentPivot.x - ss.Contours[i].joint.pivot.x
					y := drawnPositions[ss.Contours[i].joint.parentIndex].y + ss.Contours[i].joint.parentPivot.y - ss.Contours[i].joint.pivot.y
					imrect := image.Rect(x, y, x+ss.Contours[i].img.Bounds().Max.X, y+ss.Contours[i].img.Bounds().Max.Y)
					draw.Draw(img, imrect, &ss.Contours[i].img, image.Pt(0, 0), draw.Over)
					drawnPositions[i] = NewPoint(x, y)
				}
			}
		}
	}
}

func (ss *SpriteSet) DrawSprite() {
	var tex uint32
	ss.drawType = 2
	width := ss.img.Bounds().Max.X
	height := ss.img.Bounds().Max.Y

	win := render.NewWindow(height, width*2, true)
	window := win.Window
	defer glfw.Terminate()
	window.SetKeyCallback(ss.keyCallBack)
	window.SetMouseButtonCallback(ss.mouseButtonCallback)
	window.SetCursorPosCallback(ss.cursorPosCallback)
	VAO := makeVao()
	program := win.Program
	gl.GenTextures(1, &tex)
	for !window.ShouldClose() {
		imgRect := image.Rect(0, 0, width*2, height)
		img := image.NewRGBA(imgRect)
		for x := 0; x < imgRect.Max.X; x++ {
			for y := 0; y < imgRect.Max.Y; y++ {
				img.Set(x, y, color.White)
			}
		}
		imr1 := image.Rect(0, 0, ss.img.Bounds().Max.X, ss.img.Bounds().Max.Y)
		draw.Draw(img, imr1, ss.img, image.Pt(0, 0), draw.Src)

		red := color.RGBA{255, 0, 0, 255}
		green := color.RGBA{0, 255, 0, 255}
		blue := color.RGBA{0, 0, 255, 255}

		for i := range ss.Contours {
			if ss.drawType == 0 {
				if ss.selectedContour == i {
					drawRectangle(ss.Contours[i].boundingRect, img, green)
				} else {
					drawRectangle(ss.Contours[i].boundingRect, img, red)
				}
			} else if ss.drawType == 1 {
				if ss.selectedContour == i {
					drawContour(ss.Contours[i].wrapPoints, img, green, ss.Contours[i].selectedIndex, blue)
				} else {
					drawContour(ss.Contours[i].wrapPoints, img, red, ss.Contours[i].selectedIndex, blue)
				}
			} else if ss.drawType == 2 {
				if ss.selectedContour == i {
					drawRectangle(ss.Contours[i].boundingRect, img, green)
				} else {
					drawRectangle(ss.Contours[i].boundingRect, img, red)
				}
				ss.DrawJoints(img)
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

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Max.X), int32(img.Rect.Max.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
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
