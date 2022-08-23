package evisiable_test

import (
	"Solution/evisiable"
	"image/color"
	"testing"
)

func TestDrawLine(t *testing.T) {
	canvas := evisiable.NewCanvas(500, 500)
	canvas.DrawLine(100, 100, 400, 400, color.Black)
	canvas.Save("../1.png")
}

func TestDrawRect(t *testing.T) {
	canvas := evisiable.NewCanvas(500, 500)
	canvas.DrawLine(100, 100, 100, 400, color.Black)
	canvas.DrawLine(100, 100, 400, 100, color.Black)
	canvas.DrawLine(100, 400, 400, 400, color.Black)
	canvas.DrawLine(400, 100, 400, 400, color.Black)
	canvas.ColorifyZone(101, 101, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
	canvas.Save("../1.png")
}
