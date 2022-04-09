package perlin

import (
	"math"
	"math/rand"

	"github.com/seehuhn/mt19937"
)

var (
	_grad3 = [][]float32{
		{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
		{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
		{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1}}
)

type Perlin3D struct {
	gradient_array *[64 * 64 * 32]uint8

	xs int //x sample density
	zs int //y sample density

	ssize int
	xsize int
	zsize int
}

func NewPerlin3D(seed int64) Perlin3D {
	ga := new([64 * 64 * 32]uint8)
	rng := rand.New(mt19937.New())
	rng.Seed(seed)
	for i := 0; i < 64*64*32; i++ {
		ga[i] = uint8(rng.Uint64() % 12)
	}
	return Perlin3D{
		xs:             8,
		zs:             4,
		ssize:          64 * 64 * 32,
		xsize:          64,
		zsize:          32,
		gradient_array: ga,
	}
}

func (p *Perlin3D) dot(g []float32, x, y, z float32) float32 {
	return g[0]*x + g[1]*y + g[2]*z
}

func (p *Perlin3D) mix(a, b, t float32) float32 {
	return a + t*(b-a) //optimized version
}

func (p *Perlin3D) fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

func (p *Perlin3D) get_gradient(x, y, z int) uint8 {
	x = x % 64 //replace with bitmask
	y = y % 64
	z = z % 32

	if x+y*64+z*64*64 >= p.ssize {
		return 0
	}

	return p.gradient_array[x+y*64+z*64*64]
}

func (p *Perlin3D) base(x, y, z float32) float32 {
	x /= 8.0 //replace with multiplication
	y /= 8.0
	z /= 4.0
	X := int(math.Floor(float64(x)))
	Y := int(math.Floor(float64(y)))
	Z := int(math.Floor(float64(z)))

	x = x - float32(X)
	y = y - float32(Y)
	z = z - float32(Z)

	gi000 := p.get_gradient(X+0, Y+0, Z+0)
	gi001 := p.get_gradient(X+0, Y+0, Z+1)
	gi010 := p.get_gradient(X+0, Y+1, Z+0)
	gi011 := p.get_gradient(X+0, Y+1, Z+1)

	gi100 := p.get_gradient(X+1, Y+0, Z+0)
	gi101 := p.get_gradient(X+1, Y+0, Z+1)
	gi110 := p.get_gradient(X+1, Y+1, Z+0)
	gi111 := p.get_gradient(X+1, Y+1, Z+1)

	n000 := p.dot(_grad3[gi000], x, y, z)
	n100 := p.dot(_grad3[gi100], x-1, y, z)
	n010 := p.dot(_grad3[gi010], x, y-1, z)
	n110 := p.dot(_grad3[gi110], x-1, y-1, z)
	n001 := p.dot(_grad3[gi001], x, y, z-1)
	n101 := p.dot(_grad3[gi101], x-1, y, z-1)
	n011 := p.dot(_grad3[gi011], x, y-1, z-1)
	n111 := p.dot(_grad3[gi111], x-1, y-1, z-1)

	u := p.fade(x)
	v := p.fade(y)
	w := p.fade(z)

	//u := x;
	//v := y;
	//w := z;

	nx00 := p.mix(n000, n100, u)
	nx01 := p.mix(n001, n101, u)
	nx10 := p.mix(n010, n110, u)
	nx11 := p.mix(n011, n111, u)
	nxy0 := p.mix(nx00, nx10, v)
	nxy1 := p.mix(nx01, nx11, v)
	nxyz := p.mix(nxy0, nxy1, w)

	return nxyz * 0.707106781 //-1 to 1
}

func (p *Perlin3D) Noise(x, y, z float32) float32 {
	return p.base(x, y, z)
}

func (p *Perlin3D) one_over_f(x, y, z float32) float32 {
	var tmp float32
	tmp = 0.0
	tmp += p.base(x, y, z)
	tmp += 0.50 * p.base(2*x, 2*y, 2*z)
	tmp += 0.25 * p.base(4*x, 4*y, 2*z)
	return tmp
}
