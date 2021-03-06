package hexagen

import (
	"fmt"
	"image/color"
	"math"
)

func adjacent(a, b FaceAddr) bool {
	if a.Orientation == b.Orientation {
		return false
	}

	if a.X == b.X && a.Y == b.Y {
		return true
	}

	if a.X == b.X {
		return int(math.Abs(float64(a.Y-b.Y))) == 1
	}
	if a.Y == b.Y {
		return int(math.Abs(float64(a.X-b.X))) == 1
	}
	return false
}

type mat struct {
	a, b, c, d float64
}

func (m mat) rmul(v vec) vec {
	return vec{m.a*v.x + m.b*v.y, m.c*v.x + m.d*v.y}
}

type vec struct {
	x, y float64
}

func (v vec) add(w vec) vec {
	return vec{v.x + w.x, v.y + w.y}
}

func (v vec) String() string {
	return fmt.Sprintf("{% 1.1f, % 1.1f}", v.x, v.y)
}

func resolve(x, y float64) FaceAddr {
	v := vec{x * 4, y * 4}

	m := mat{1, -1 / math.Sqrt(3), 0, 2 / math.Sqrt(3)}

	w := m.rmul(v)
	w.x += 1

	// in welchem dreieck?
	o := Orientation(w.x-math.Floor(w.x)+w.y-math.Floor(w.y) < 1)

	addr := FaceAddr{X: int(math.Floor(w.x)), Y: int(math.Floor(w.y)), Orientation: o}

	return addr
}

func inhexagon(addr FaceAddr) bool {
	acc := true

	acc = acc && (addr.X > 0 || addr.Y > 0) // exclude (0, 0)
	acc = acc && (addr.X >= 0)
	acc = acc && (addr.Y >= 0)

	acc = acc && !(addr.X == 1 && addr.Y == 0 && addr.Orientation == Up)
	acc = acc && !(addr.X == 0 && addr.Y == 1 && addr.Orientation == Up)

	acc = acc && (addr.X < 4)
	acc = acc && (addr.Y < 4)

	acc = acc && !(addr.X == 3 && addr.Y == 2 && addr.Orientation == Down)
	acc = acc && !(addr.X == 2 && addr.Y == 3 && addr.Orientation == Down)

	acc = acc && (addr.X < 3 || addr.Y < 3) // exclude (3, 3)

	return acc
}

// Generate takes some data and a (pixel) width and generates a Grid that can be encoded into an image.
func Generate(data []byte, width float64) *Grid {
	var g Grid
	g.m = make(map[FaceAddr]color.CMYK, 0)
	g.w = width

	cur := FaceAddr{2, 0, true}
	prev := FaceAddr{2, -1, false}
	delta := cur.Sub(prev)

	for _, b := range data {

		for j := 0; j < 2; j++ {
			//fmt.Println(cur)
			input := (b & 1) > 0
			b >>= 1

			n := next(cur, cur.Sub(delta), input)

			prev = cur
			cur = n
			delta = cur.Sub(prev)

			cur = wrap(cur)

			col := g.m[cur]
			//col.A = 0xff

			col.C += b & 1
			b >>= 1

			col.M += b & 1
			b >>= 1

			col.Y += b & 1
			b >>= 1

			g.m[cur] = col
		}
	}

	var max float64

	for _, col := range g.m {
		max = math.Max(
			math.Max(float64(col.C), float64(col.M)),
			math.Max(float64(col.Y), float64(max)),
		)
	}

	scale := byte(0xff / max)

	for addr, col := range g.m {
		col.C *= scale
		col.M *= scale
		col.Y *= scale

		g.m[addr] = col
	}

	return &g
}
