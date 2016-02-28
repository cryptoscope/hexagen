package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Orientation bool

func (o Orientation) String() string {
	if o == Up {
		return "Up"
	}

	return "Down"
}

const (
	Down Orientation = false
	Up               = true
)

type FaceAddr struct {
	X, Y        int
	Orientation Orientation
}

func (a FaceAddr) Sub(b FaceAddr) FaceAddr {
	return FaceAddr{X: a.X - b.X, Y: a.Y - b.Y, Orientation: a.Orientation != b.Orientation}
}

func (addr FaceAddr) String() string {
	return fmt.Sprintf("{% 1d, % 1d, %v}", addr.X, addr.Y, addr.Orientation)
}

type Grid map[FaceAddr]color.CMYK

func (g Grid) At(x, y int) color.Color {
	height := width * math.Sqrt(3) / 2
	padding := int(width-height) / 2

	y -= padding
	if y < 0 || y > int(float64(width)*math.Sqrt(3)/2) {
		return color.Transparent
	}

	y = int(float64(width)/2*math.Sqrt(3)) - y
	addr := resolve(float64(x)/width, float64(y)/width)

	if !inhexagon(addr) {
		return color.Transparent
	}

	if col, ok := g[addr]; ok {
		return col
	}

	return color.White
}

func (g Grid) ColorModel() color.Model {
	return color.NRGBAModel
}

func (g Grid) Bounds() image.Rectangle {
	return image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(width), Y: int(width)}}
}
