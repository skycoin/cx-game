package spine

import "fmt"

// alpha non-premultiplied color
type Color struct{ R, G, B, A float32 }

func (c Color) WithAlpha(a float32) Color {
	c.A = a
	return c
}

func u8(v float32) uint8 {
	if v <= 0 {
		return 0
	} else if v >= 1 {
		return 0xFF
	}
	return uint8(v * 0xFF)
}

func u32(v float32) uint32 {
	if v <= 0 {
		return 0
	} else if v >= 1 {
		return 0xFFFF
	}
	return uint32(v * 0xFFFF)
}

// create color from alpha-non-premultiplied r, g, b, a
func RGBA(r, g, b, a float32) Color { return Color{r, g, b, a} }

// RGBA returns alpha premultiplied color in range[0..0xFFFF]
func (c Color) RGBA() (r, g, b, a uint32) {
	c.R *= c.A
	c.G *= c.A
	c.B *= c.A
	return u32(c.R), u32(c.G), u32(c.B), u32(c.A)
}

// NRGBA32 returns alpha non-premultiplied color components
func (c Color) NRGBA32() (r, g, b, a float32) {
	return c.R, c.G, c.B, c.A
}

// NRGBA64 returns alpha non-premultiplied color components
func (c Color) NRGBA64() (r, g, b, a float64) {
	return float64(c.R), float64(c.G), float64(c.B), float64(c.A)
}

// RGBA32 returns alpha premultiplied color components
func (c Color) RGBA32() (r, g, b, a float32) {
	return c.R * c.A, c.G * c.A, c.B * c.A, c.A
}

// RGBA64 returns alpha premultiplied color components
func (c Color) RGBA64() (r, g, b, a float64) {
	return float64(c.R * c.A), float64(c.G * c.A), float64(c.B * c.A), float64(c.A)
}

func (c Color) String() string {
	return fmt.Sprintf("RGBA{%.2f,%.2f,%.2f,%.2f}", c.R, c.G, c.B, c.A)
}

func lerpColor(a, b Color, p float32) Color {
	if p <= 0 {
		return a
	} else if p >= 1 {
		return b
	}
	// TODO: optimize
	return Color{
		lerp(a.R, b.R, p),
		lerp(a.G, b.G, p),
		lerp(a.B, b.B, p),
		lerp(a.A, b.A, p),
	}
}
