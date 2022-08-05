package emodule

type Point struct {
	X float64
	Y float64
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

type LineFragment struct {
	ParentId  int
	ParentRef *Line
	Start     Point
	End       Point
}

type Line struct {
	Id        int
	Width     float64
	Frangment []*LineFragment
}

type Rect struct {
	LeftTop     Point
	LeftButtom  Point
	RightTop    Point
	RightButtom Point
}

func MinOfInt(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxOfInt(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func TwoInteger2Float(a, b int) float64 {
	ans := float64(a)
	bb := b
	cnt := 0
	for bb > 0 {
		bb /= 10
		cnt++
	}
	floatb := float64(b)
	for cnt > 0 {
		floatb /= 10
		cnt--
	}
	ans += floatb
	return ans
}

func NewRect(a Point, b Point) Rect {
	minX := MinOfInt(a.X, b.X)
	minY := MinOfInt(a.Y, b.Y)
	maxX := MaxOfInt(a.X, b.X)
	maxY := MaxOfInt(a.Y, b.Y)
	leftTop := NewPoint(minX, minY)
	rightButtom := NewPoint(maxX, maxY)
	return Rect{
		LeftTop:     leftTop,
		RightButtom: rightButtom,
		LeftButtom:  NewPoint(rightButtom.X, leftTop.Y),
		RightTop:    NewPoint(leftTop.X, rightButtom.Y),
	}
}

func (p *LineFragment) GetMarginRect() Rect {
	a := p.Start
	b := p.End
	minX := MinOfInt(a.X, b.X)
	minY := MinOfInt(a.Y, b.Y)
	maxX := MaxOfInt(a.X, b.X)
	maxY := MaxOfInt(a.Y, b.Y)
	leftTop := NewPoint(minX, minY)
	rightButtom := NewPoint(maxX, maxY)
	return NewRect(NewPoint(leftTop.X-p.ParentRef.Width, leftTop.Y-p.ParentRef.Width),
		NewPoint(rightButtom.X+p.ParentRef.Width, rightButtom.Y+p.ParentRef.Width))
}
