package evisiable

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Drawable interface {
	Draw(ctx *Canvas)
}

type Canvas struct {
	Height   int
	Width    int
	Img      *image.NRGBA
	DrawItem []*Drawable
}

func NewCanvas(height, width int) *Canvas {
	ans := &Canvas{
		Height: height,
		Width:  width,
		Img:    image.NewNRGBA(image.Rect(0, 0, height, width)),
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			ans.Img.Set(i, j, color.White)
		}
	}
	return ans
}
func abs(x int) int {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

func (p *Canvas) DrawPoint(x, y int, weight ...int) {
	w := 4
	if len(weight) != 0 {
		w = weight[0]
	}
	for i := 0; i < w; i++ {
		for j := 0; j < 4; j++ {
			p.Img.Set(i, j, color.Black)
		}
	}
}

func (p *Canvas) DrawLine(x0, y0, x1, y1 int, color color.Color) {
	dx := abs(x1 - x0)

	dy := abs(y1 - y0)

	sx, sy := 1, 1

	if x0 >= x1 {

		sx = -1

	}

	if y0 >= y1 {

		sy = -1

	}

	err := dx - dy

	for {

		p.Img.Set(x0, y0, color)

		if x0 == x1 && y0 == y1 {

			return

		}

		e2 := err * 2

		if e2 > -dy {

			err -= dy

			x0 += sx

		}

		if e2 < dx {

			err += dx

			y0 += sy

		}

	}
}

func (p *Canvas) ColorifyZone(x, y int, targetColor color.Color) {
	type Pix struct {
		X int
		Y int
	}
	var dx [4]int = [4]int{1, -1, 0, 0}
	var dy [4]int = [4]int{0, 0, 1, -1}
	pixs := make([]Pix, 0)
	initColor := p.Img.At(x, y)
	pixs = append(pixs, Pix{x, y})
	p.Img.Set(pixs[0].X, pixs[0].Y, targetColor)
	for len(pixs) > 0 {
		now := pixs[0]
		pixs = pixs[1:]
		//fmt.Println(now)
		for i := 0; i < 4; i++ {
			nxt := now
			nxt.X += dx[i]
			nxt.Y += dy[i]
			if nxt.X >= 0 && nxt.X < p.Height && nxt.Y >= 0 && nxt.Y < p.Width {
				if p.Img.At(nxt.X, nxt.Y) == initColor && p.Img.At(nxt.X, nxt.Y) != targetColor {
					p.Img.Set(nxt.X, nxt.Y, targetColor)
					pixs = append(pixs, nxt)
				}
			}
		}
	}
}

func (p *Canvas) DrawAllItem() {
	for _, v := range p.DrawItem {
		(*v).Draw(p)
	}
}

func (p *Canvas) Clear() {
	p.Img = image.NewNRGBA(image.Rect(0, 0, p.Height, p.Width))
	for i := 0; i < p.Height; i++ {
		for j := 0; j < p.Width; j++ {
			p.Img.Set(i, j, color.White)
		}
	}
}

func (p *Canvas) Save(path string) {
	p.DrawAllItem()
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	png.Encode(file, p.Img)
}
