package perlin

import (
	"math"
	"math/rand"

	"github.com/seehuhn/mt19937"
)

type Perlin2D struct {
	gradient_array []uint8
	grad           [][2]float32

	ssize int
	xsize int
	xs    int

	xscale   float32
	grad_max int
}

func NewPerlin2D(seed int64, _x, _xs, _grad_max int) Perlin2D {
	_xsize := _x / _xs
	_ssize := _xsize * _xsize
	_xscale := 1.0 / float32(_xs)
	rng := rand.New(mt19937.New())
	rng.Seed(seed)
	gr := make([][2]float32, _grad_max)
	ga := make([]uint8, _ssize)
	for i := 0; i < _ssize; i++ {
		ga[i] = uint8(rng.Uint64() % 12)
	}
	for i := 0; i < _grad_max; i++ {
		x := 2*mrandf(rng) - 1
		y := 2*mrandf(rng) - 1

		len := math.Sqrt(x*x + y*y)
		x /= len
		y /= len

		gr[i][0] = float32(x)
		gr[i][1] = float32(y)
	}
	return Perlin2D{
		xs:             _xs,
		ssize:          _ssize,
		xsize:          _xsize,
		grad_max:       _grad_max,
		xscale:         float32(_xscale),
		gradient_array: ga,
		grad:           gr,
	}
}

func mrandf(r *rand.Rand) float64 {
	return float64(r.Uint32()) / 4294967295.0
}

func (p *Perlin2D) dot(g [2]float32, x, y float32) float32 {
	return g[0]*x + g[1]*y
}

func (p *Perlin2D) mix(a, b, t float32) float32 {
	return a + t*(b-a) //optimized version
}

func (p *Perlin2D) fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

func (p *Perlin2D) get_gradient(x, y int) uint8 {
	x = x % p.xsize //replace with bitmask
	y = y % p.xsize

	if x+y*p.xsize >= p.ssize {
		return 0
	}

	return p.gradient_array[x+y*p.xsize]
}

func (p *Perlin2D) Base(x, y float32) float32 {
	x *= p.xscale //replace with multiplication
	y *= p.xscale
	//get grid point
	X := int(math.Floor(float64(x)))
	Y := int(math.Floor(float64(y)))

	x = x - float32(X)
	y = y - float32(Y)

	gi00 := p.get_gradient(X+0, Y+0)
	gi01 := p.get_gradient(X+0, Y+1)
	gi10 := p.get_gradient(X+1, Y+0)
	gi11 := p.get_gradient(X+1, Y+1)

	// Calculate noise contributions from each of the eight corners
	n00 := p.dot(p.grad[gi00], x, y)
	n10 := p.dot(p.grad[gi10], x-1, y)
	n01 := p.dot(p.grad[gi01], x, y-1)
	n11 := p.dot(p.grad[gi11], x-1, y-1)
	// Compute the fade curve value for each of x, y, z

	u := p.fade(x)
	v := p.fade(y)
	//	float u = x;
	//	float v = y;

	// Interpolate along x the contributions from each of the corners
	nx00 := p.mix(n00, n10, u)
	nx10 := p.mix(n01, n11, u)
	// Interpolate the four results along y
	nxy := p.mix(nx00, nx10, v)

	if nxy < -1 || nxy > 1 {
		return 0
	}
	return nxy //-1 to 1
}

func (p *Perlin2D) Noise(x, y float32) float32 {
	return p.Base(x, y)
}

func (p *Perlin2D) One_over_f(x, y float32) float32 {
	var tmp float32 = 0
	tmp += p.Base(x, y)
	tmp += 0.50 * p.Base(2*x, 2*y)
	tmp += 0.25 * p.Base(4*x, 4*y)
	tmp += 0.125 * p.Base(8*x, 8*y)
	tmp += 0.0625 * p.Base(16*x, 16*y)
	return tmp
}

func (p *Perlin2D) One_over_f_pers(x, y, persistence float32) float32 {
	var tmp float32 = 0.0
	var m float32 = 1.0

	tmp = p.Base(x, y)

	m *= persistence

	tmp += m * p.Base(2*x, 2*y)
	m *= persistence

	tmp += m * p.Base(4*x, 4*y)
	m *= persistence

	tmp += m * p.Base(8*x, 8*y)

	m *= persistence
	tmp += m * p.Base(16*x, 16*y)
	return tmp
}

//order 0 is base
func (p *Perlin2D) order(x, y, persistence float32, order int) float32 {
	var tmp float32 = 0.0
	var m float32 = 1.0
	var b float32 = 1.0

	for i := 0; i <= order; i++ {
		tmp += p.Base(b*x, b*y)
		m *= persistence
		b *= 2
	}
	return tmp
}

func (p *Perlin2D) abs(x, y float32) float32 {
	tmp := float64(p.Base(x, y))
	return float32(math.Sqrt(tmp * tmp))
}
