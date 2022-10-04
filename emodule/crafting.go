package emodule

import (
	"Solution/evisiable"
	"fmt"
	"image/color"

	"github.com/liyue201/gostl/ds/set"
)

/*
hint
type Instance struct {
	Name          string
	InstanceChain []*Instance
	Parts         []*Instance
	Block         *Block
	Offset        Point
	Rotate        string
	RotatedBlock  *Block
}

type Block struct {
	Name          string
	PointSet      []Point
	LineFragments []LineFragment
	Rotators      map[string]*Rotator
}

*/

var ImportantPointSet *set.Set
var ImportantLineSet *set.MultiSet
var BlockFindInstace map[*Block][]*Instance

func init() {
	BlockFindInstace = make(map[*Block][]*Instance)
	ImportantLineSet = set.NewMultiSet(set.WithGoroutineSafe())
	ImportantPointSet = set.New(set.WithGoroutineSafe(), set.WithKeyComparator(func(a, b interface{}) int {
		pointA := a.(*ImportantPoint)
		pointB := b.(*ImportantPoint)
		if pointA.X == pointB.X && pointA.Y == pointB.Y {
			return 0
		}
		if pointA.X != pointB.X {
			if pointA.X < pointB.X {
				return -1
			} else {
				return 1
			}
		}
		if pointA.Y < pointB.Y {
			return -1
		} else {
			return 1
		}
	}))
}

type ImportantPoint struct {
	X             float64
	Y             float64
	BlockOwner    *Block
	InstanceOwner *Instance
	NodeRef       *NodeOfBlock
}

type ImportantLine struct {
	A *ImportantPoint
	B *ImportantPoint
}

func (p *ImportantPoint) Draw(ctx evisiable.Canvas) {
	if p.NodeRef != nil {
		ctx.DrawPoint(int(p.X), int(p.Y), 8)
	} else {
		ctx.DrawPoint(int(p.X), int(p.Y))
	}
}

func ImprotInstaces() {
	for _, instance := range InstanceSet {
		block := instance.Block
		_, hasBlock := BlockFindInstace[block]
		if hasBlock {
			BlockFindInstace[block] = append(BlockFindInstace[block], instance)
		} else {
			BlockFindInstace[block] = make([]*Instance, 0)
			BlockFindInstace[block] = append(BlockFindInstace[block], instance)
		}
	}
}

func AddInstanceBorderLine() error {
	for _, instance := range InstanceSet {
		pointSet := make([]Point, 0)
		block := instance.Block
		rotator := block.Rotators[instance.Rotate]
		if rotator == nil {
			return fmt.Errorf("don't has rotate option")
		}
		for _, v := range block.PointSet {
			pointSet = append(pointSet, rotator.f(v))
		}
		for i := 0; i < len(pointSet)-1; i++ {
			a := NewImportantPoint(NewPoint(pointSet[i].X+instance.Offset.X, pointSet[i].Y+instance.Offset.Y), instance)
			b := NewImportantPoint(NewPoint(pointSet[i+1].X+instance.Offset.X, pointSet[i+1].Y+instance.Offset.Y), instance)
			NewImportantLine(a, b)
		}
		len := len(pointSet)
		fmt.Println(len)
		a := NewImportantPoint(NewPoint(pointSet[0].X+instance.Offset.X, pointSet[0].Y+instance.Offset.Y), instance)
		b := NewImportantPoint(NewPoint(pointSet[len-1].X+instance.Offset.X, pointSet[len-1].Y+instance.Offset.Y), instance)
		NewImportantLine(a, b)
	}
	return nil
}

func NewImportantLine(a, b *ImportantPoint) *ImportantLine {
	ans := ImportantLine{
		A: a,
		B: b,
	}
	ImportantLineSet.Insert(&ans)
	return &ans
}

func (p *ImportantLine) Draw(ctx evisiable.Canvas) {
	ctx.DrawLine(int(p.A.X), int(p.A.Y), int(p.B.X), int(p.B.Y), color.NRGBA{0, 0, 255, 255})
}

func ImportNodes() error {
	for _, node := range NodeSet {
		instance := node.InstanceRef
		blockRef := instance.Block
		importantPoint := &ImportantPoint{}
		importantPoint.BlockOwner = blockRef
		importantPoint.InstanceOwner = instance
		importantPoint.NodeRef = node
		importantPoint.X = node.Pos.X
		importantPoint.Y = node.Pos.Y
		x := node.Pos.X
		y := node.Pos.Y
		x -= node.InstanceRef.Offset.X
		y -= node.InstanceRef.Offset.Y
		blockPos := blockRef.Rotators[instance.Rotate].f(NewPoint(x, y))
		AddInportantPointInBlock(blockRef, blockPos)
		// fmt.Println(importantPoint.X, importantPoint.Y, "#######")
		// fmt.Print("{")
		// for it := ImportantPointSet.Begin(); !it.Equal(ImportantLineSet.Last()); it.Next() {
		// 	fmt.Print("(", it.Value().(*ImportantPoint).X, it.Value().(*ImportantPoint).Y, ")", ",")
		// }
		// fmt.Print("}\n")
		pos := ImportantPointSet.Find(importantPoint)
		if !pos.IsValid() {
			return fmt.Errorf("should has ImportantPoint %v", importantPoint)
		}
		ImportantPointSet.Erase(importantPoint)
		ImportantPointSet.Insert(importantPoint)
		NodeFindPoint[node] = importantPoint
	}
	return nil
}

func NewImportantPoint(p Point, instance *Instance) *ImportantPoint {
	importantPoint := &ImportantPoint{}
	if instance != nil {
		importantPoint.BlockOwner = instance.Block
	}
	importantPoint.InstanceOwner = instance
	importantPoint.NodeRef = nil
	importantPoint.X = p.X
	importantPoint.Y = p.Y
	ImportantPointSet.Insert(importantPoint)
	return importantPoint
}

func AddInportantPointInBlock(block *Block, point Point) error {
	_, hasBlock := BlockFindInstace[block]
	if !hasBlock {
		return fmt.Errorf("can't find block to add point")
	}
	for _, v := range BlockFindInstace[block] {
		tempPoint := NewPoint(point.X, point.Y)
		tempPoint = v.Block.Rotators[v.Rotate].f(tempPoint)
		tempPoint.X = v.Offset.X + tempPoint.X
		tempPoint.Y = v.Offset.Y + tempPoint.Y
		NewImportantPoint(tempPoint, v)
	}
	return nil
}

func GenPath(a *ImportantPoint, b *ImportantPoint, cmd int) []*ImportantLine {
	ans := make([]*ImportantLine, 0)
	if cmd == 1 {
		if a.X == b.X || a.Y == b.Y {
			tmp := &ImportantLine{
				A: a,
				B: b,
			}
			ans = append(ans, tmp)
			return ans
		}
		mid := NewImportantPoint(NewPoint(a.X, b.Y), nil)
		one := &ImportantLine{
			A: a,
			B: mid,
		}
		two := &ImportantLine{
			A: mid,
			B: b,
		}
		ans = append(ans, one)
		ans = append(ans, two)
		return ans
	}
	if cmd == 2 {
		if a.X == b.X || a.Y == b.Y {
			tmp := &ImportantLine{
				A: a,
				B: b,
			}
			ans = append(ans, tmp)
			return ans
		}
		mid := NewImportantPoint(NewPoint(b.X, a.Y), nil)
		one := &ImportantLine{
			A: a,
			B: mid,
		}
		two := &ImportantLine{
			A: mid,
			B: b,
		}
		ans = append(ans, one)
		ans = append(ans, two)
		return ans
	}
	return ans
}

func LineParallelJudge(x *ImportantLine) bool {
	for it := ImportantLineSet.Begin(); it.IsValid(); it.Next() {
		if TwoLineParallelJudge(it.Value().(*ImportantLine), x, 5) {
			return true
		}
	}
	return false
}

func AbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TwoLineParallelJudge(a, b *ImportantLine, param int) bool {
	type1 := -1
	type2 := -1
	if a.A.X == a.B.X {
		type1 = 1
	} else {
		type1 = 2
	}
	if b.A.X == b.B.X {
		type2 = 1
	} else {
		type2 = 2
	}
	if type1 != type2 {
		return false
	}
	if type1 == 1 {
		if AbsFloat64(a.A.X-b.A.X) > float64(param) {
			return false
		}
	} else {
		if AbsFloat64(a.A.Y-b.A.Y) > float64(param) {
			return false
		}
	}
	if type1 == 1 {
		if MaxOfInt(a.A.Y, a.B.Y) > MinOfInt(b.A.Y, b.B.Y) || MaxOfInt(b.A.Y, b.B.Y) > MinOfInt(a.A.Y, a.B.Y) {
			return true
		} else {
			return false
		}
	} else {
		if MaxOfInt(a.A.X, a.B.X) > MinOfInt(b.A.X, b.B.X) || MaxOfInt(b.A.X, b.B.X) > MinOfInt(a.A.X, a.B.X) {
			return true
		} else {
			return false
		}
	}
}

func RouteBlock() {

}
