package emodule

import (
	"Solution/evisiable"
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
)

/*<CONSTRAfloat64S> 5 5 3 <CONSTRAfloat64S>
<BLOCK>
<block_A> {{0 0} {0 100} {100 100} {100 0}} <block_A>
<block_B> {{0 0} {0 50} {200 50} {200 0}} <block_B>
<block_C> {{0 0} {0 25} {120 25} {120 0}} <block_C>
<BLOCK>
<INSTANCE>
<A_1> block_A {10 10} R0 <A_1>
<A_2> block_A {210 10} MY <A_2>
<B> block_B {10 110} R0 <B>
<B/C> block_C {80 125} R0 <B/C>
<INSTANCE>
<NODE>
<A_1/N1> {60 60} <A_1/N1>
<A_2/N1> {160 60} <A_2/N1>
<B/N1> {60 135} <B/N1>
<B/C/N1> {180 135} <B/C/N1>
<NODE>
<FLY_LINE>
<A_1/N1 B/N1>
<A_1/N1 B/C/N1>
<A_2/N1 B/C/N1>
<FLY_LINE>*/

type Node struct {
	Tag      string
	HasChild bool
	Content  string
	Children []*Node
	Father   *Node
}

type XMLTree struct {
	Root *Node
}

func (p *Node) addToFather(fa *Node) {
	fa.HasChild = true
	fa.Children = append(fa.Children, p)
	p.Father = fa
}

func ReadXML(reader io.Reader) XMLTree {
	root := Node{
		Tag:      "root",
		HasChild: true,
		Children: make([]*Node, 0),
		Father:   nil,
	}
	stack := make([]*Node, 0)
	stack = append(stack, &root)
	bufreader := bufio.NewReader(reader)
	contentBuilder := strings.Builder{}
	for {
		ch, _, err := bufreader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if ch == '<' {
			builder := strings.Builder{}
			for {
				nch, _, err := bufreader.ReadRune()
				if nch == '>' {
					break
				}
				if err == io.EOF {
					break
				}
				builder.WriteRune(nch)
			}
			node := Node{
				Tag: builder.String(),
			}
			if stack[len(stack)-1].Tag == node.Tag {
				node.Content = contentBuilder.String()
				node.Children = stack[len(stack)-1].Children
				stack = stack[:len(stack)-1]
				node.addToFather(stack[len(stack)-1])
				contentBuilder.Reset()
			} else {
				if stack[len(stack)-1].Tag == "FLY_LINE" {
					node.addToFather(stack[len(stack)-1])
				} else {
					stack = append(stack, &node)
				}
			}
		} else {
			contentBuilder.WriteRune(ch)
		}
	}
	return XMLTree{
		Root: &root,
	}
}

func PrintXML() {
	file, err := os.Open("./emodule/test.xml")
	if err != nil {
		fmt.Println(err)
	}
	xml := ReadXML(file)
	node := xml.Root
	PrintNode(node, 0)
}
func PrintNode(p *Node, cnt int) {
	for i := 0; i < cnt; i++ {
		fmt.Print("    ")
	}
	fmt.Println(p.Tag + ":" + p.Content)
	fmt.Println("-------------------")
	for _, v := range p.Children {
		PrintNode(v, cnt+1)
	}
}

type Constrafloat64 struct {
	Height    float64
	Width     float64
	LineWidth float64
}

func NewConstrafloat64(str string) *Constrafloat64 {
	ans := Constrafloat64{}
	fmt.Sscan(str, &ans.Height, &ans.Height, &ans.Width, &ans.LineWidth)
	return &ans
}

var BlockSet map[string]*Block
var InstanceSet map[string]*Instance
var NodeSet map[string]*NodeOfBlock
var FlyLineSet map[string]*FlyLine

func init() {
	BlockSet = make(map[string]*Block)
	InstanceSet = make(map[string]*Instance)
	NodeSet = make(map[string]*NodeOfBlock)
	FlyLineSet = make(map[string]*FlyLine)
}

type Block struct {
	Name          string
	PointSet      []Point
	LineFragments []LineFragment
	Rotators      map[string]*Rotator
}

func NewBlock(name string, str string) *Block {
	ans := Block{
		LineFragments: make([]LineFragment, 0),
	}
	ans.Name = name
	builder := strings.Builder{}
	for _, v := range str {
		if v != '{' && v != '}' {
			builder.WriteRune(v)
		}
	}
	str = builder.String()
	numbers := make([]float64, 0)
	str = strings.TrimSpace(str)
	str = str + " "
	prefix := 0
	suffix := 0
	flag := false
	for _, v := range str {
		if v >= '0' && v <= '9' {
			tmp := v - '0'
			if !flag {
				prefix = prefix * 10
				prefix += int(tmp)
			} else {
				suffix = suffix * 10
				suffix += int(tmp)
			}
		} else if v == '.' {
			flag = true
		} else {
			numbers = append(numbers, TwoInteger2Float(prefix, suffix))
			prefix = 0
			suffix = 0
			flag = false
		}
	}
	for i := 0; i+1 < len(numbers); i += 2 {
		ans.PointSet = append(ans.PointSet, NewPoint(numbers[i], numbers[i+1]))
	}
	ans.Rotators = make(map[string]*Rotator)
	ans.Rotators["R0"] = Rotate(&ans, "R0")
	ans.Rotators["MX"] = Rotate(&ans, "MX")
	ans.Rotators["MY"] = Rotate(&ans, "MY")
	ans.Rotators["R180"] = Rotate(&ans, "R180")
	BlockSet[ans.Name] = &ans
	return &ans
}

func (p *Block) Rotate(how string) *Block {
	ans := &Block{
		Name:     p.Name + "_" + how,
		PointSet: make([]Point, 0),
	}
	rotator := p.Rotators[how]
	if rotator == nil {
		return nil
	}
	for _, v := range p.PointSet {
		ans.PointSet = append(ans.PointSet, rotator.f(v))
	}
	return ans
}

func (p *Block) PointInBlock(point Point) bool {
	cnt := 0
	for i := 0; i < len(p.PointSet)-1; i++ {
		pa := p.PointSet[i]
		pb := p.PointSet[i+1]
		if point.Y > MinOfInt(pa.Y, pb.Y) && point.Y < MaxOfInt(pa.Y, pb.Y) {
			nowX := (point.Y-pa.Y)/(pb.Y-pa.Y)*(pb.X-pa.X) + pa.X
			if nowX > point.X {
				cnt++
			}
		}
	}
	return cnt%2 == 1
}

func (p *Block) SpawnAPointInBlock() (Point, error) {
	mountedPoint := p.PointSet[0]
	var dx [4]int = [4]int{1, -1, -1, +1}
	var dy [4]int = [4]int{1, -1, +1, -1}
	for i := 0; i < 4; i++ {
		nowP := NewPoint(mountedPoint.X+float64(dx[i]), mountedPoint.Y+float64(dy[i]))
		if p.PointInBlock(nowP) {
			return nowP, nil
		}
	}
	return NewPoint(0, 0), fmt.Errorf("can't spawn point in block")
}

type Instance struct {
	Name          string
	InstanceChain []*Instance
	Parts         []*Instance
	Block         *Block
	Offset        Point
	Rotate        string
	RotatedBlock  *Block
}

func string2float64(s string) float64 {
	prefix := 0
	suffix := 0
	flag := false
	for _, v := range s {
		if v >= '0' && v <= '9' {
			if !flag {
				prefix *= 10
				prefix += int(v - '0')
			} else {
				suffix *= 10
				suffix += int(v - '0')
			}
		}
		if v == '.' {
			flag = true
		}
	}
	return TwoInteger2Float(prefix, suffix)
}

func NewInstance(name string, str string) (*Instance, error) {
	ans := Instance{}
	ans.Name = name
	str = strings.TrimSpace(str)
	strs := strings.Split(str, " ")
	if len(strs) < 4 {
		return nil, fmt.Errorf("wrong instance format")
	}
	blockName := strs[0]
	_, hasBlock := BlockSet[blockName]
	if !hasBlock {
		return nil, fmt.Errorf("don't has block %s", blockName)
	}
	ans.Block = BlockSet[blockName]
	offsetPointStr := strs[1] + " " + strs[2]
	builder := strings.Builder{}
	for _, v := range offsetPointStr {
		if v != '{' && v != '}' {
			builder.WriteRune(v)
		}
	}
	offsetPointStr = builder.String()
	offsetPointStr = strings.TrimSpace(offsetPointStr)
	points := strings.Split(offsetPointStr, " ")
	if len(points) < 2 {
		fmt.Println(len(points))
		for _, v := range points {
			fmt.Println(v)
		}
		return nil, fmt.Errorf("wrong instance format")
	}
	ans.Offset = NewPoint(string2float64(points[0]), string2float64(points[1]))
	ans.Rotate = strs[3]
	//struct consturct
	layerNames := strings.Split(ans.Name, "/")
	builder.Reset()
	for i, v := range layerNames {
		if i != 0 {
			builder.WriteRune('/')
		}
		builder.WriteString(v)
		if i != len(layerNames)-1 {
			_, hasInstance := InstanceSet[builder.String()]
			if hasInstance {
				ans.InstanceChain = append(ans.InstanceChain, InstanceSet[builder.String()])
			} else {
				fmt.Println(builder.String())
				panic(fmt.Errorf("Instance Constractor Error"))
			}
		}
		if i == len(layerNames)-2 {
			_, hasInstance := InstanceSet[builder.String()]
			if hasInstance {
				InstanceSet[builder.String()].Parts = append(InstanceSet[builder.String()].Parts, &ans)
			} else {
				panic(fmt.Errorf("Instance Constractor Error"))
			}
		}
	}
	InstanceSet[ans.Name] = &ans
	return &ans, nil
}

type RotateFunc func(Point) Point

type Rotator struct {
	f RotateFunc
}

func Rotate(block *Block, option string) *Rotator {
	minX := 1e15
	minY := 1e15
	maxX := -1e15
	maxY := -1e15
	for _, v := range block.PointSet {
		if v.X > maxX {
			maxX = v.X
		}
		if v.X < minX {
			minX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
		if v.Y < minY {
			minY = v.Y
		}
	}
	if option == "R0" {
		return &Rotator{
			f: func(p Point) Point {
				return p
			},
		}
	}
	if option == "MY" {
		return &Rotator{
			f: func(p Point) Point {
				return NewPoint(maxX-(p.X-minX)-(maxX-minX), p.Y)
			},
		}
	}
	if option == "MX" {
		return &Rotator{
			f: func(p Point) Point {
				return NewPoint(p.X, maxY-(p.Y-minY)-(maxY-minY))
			},
		}
	}
	if option == "R180" {
		return &Rotator{
			f: func(p Point) Point {
				return NewPoint(maxX-(p.X-minX)-(maxX-minX), maxY-(p.Y-minY)-(maxY-minY))
			},
		}
	}
	return nil
}

func (p *Block) Draw(ctx *evisiable.Canvas) {
	for i := 0; i < len(p.PointSet)-1; i++ {
		ctx.DrawLine(int(p.PointSet[i].X), int(p.PointSet[i].Y), int(p.PointSet[i+1].X), int(p.PointSet[i+1].Y), color.Black)
	}
	last := len(p.PointSet) - 1
	ctx.DrawLine(int(p.PointSet[0].X), int(p.PointSet[0].Y), int(p.PointSet[last].X), int(p.PointSet[last].Y), color.Black)
	colorPoint, err := p.SpawnAPointInBlock()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	ctx.ColorifyZone(int(colorPoint.X), int(colorPoint.Y), color.RGBA{255, 125, 125, 255})
}

func (p *Instance) Draw(ctx *evisiable.Canvas) {
	ctx.AddOffset(int(p.Offset.X), int(p.Offset.Y))
	p.Block.Rotate(p.Rotate).Draw(ctx)
	ctx.AddOffset(int(-p.Offset.X), int(-p.Offset.Y))
	for _, v := range p.Parts {
		v.Draw(ctx)
	}
}

type NodeOfBlock struct {
	Name         string
	Pos          Point
	InstanceRef  *Instance
	InstanceName string
	NodeName     string
}

func NewNode(name string, str string) (*NodeOfBlock, error) {
	ans := NodeOfBlock{}
	ans.Name = name
	slashPos := strings.LastIndex(name, "/")
	ans.InstanceName = name[:slashPos]
	ans.NodeName = name[slashPos+1:]
	_, hasInstance := InstanceSet[ans.InstanceName]
	if hasInstance {
		ans.InstanceRef = InstanceSet[ans.InstanceName]
	} else {
		return nil, fmt.Errorf("can't construct Node because there isn't Instance named %v", ans.InstanceName)
	}
	builder := strings.Builder{}
	for _, v := range str {
		if v != '{' && v != '}' {
			builder.WriteRune(v)
		}
	}
	str = builder.String()
	str = strings.TrimSpace(str)
	points := strings.Split(str, " ")
	if len(points) < 2 {
		fmt.Println(len(points))
		for _, v := range points {
			fmt.Println(v)
		}
		return nil, fmt.Errorf("can't construct Node because format error %v", ans.InstanceName)
	}
	ans.Pos = NewPoint(string2float64(points[0]), string2float64(points[1]))
	NodeSet[ans.Name] = &ans
	return &ans, nil
}

type FlyLine struct {
	NodeAName string
	NodeBName string
	NodeARef  *NodeOfBlock
	NodeBRef  *NodeOfBlock
}

func NewFlyLine(name string) (*FlyLine, error) {
	ans := FlyLine{}
	names := strings.Split(name, " ")
	if len(names) != 2 {
		return nil, fmt.Errorf("flyline format error")
	}
	ans.NodeAName = names[0]
	ans.NodeBName = names[1]
	_, hasNode := NodeSet[ans.NodeAName]
	if !hasNode {
		return nil, fmt.Errorf("can't find Node %v", ans.NodeAName)
	}
	_, hasNode = NodeSet[ans.NodeBName]
	if !hasNode {
		return nil, fmt.Errorf("can't find Node %v", ans.NodeBName)
	}
	ans.NodeARef = NodeSet[ans.NodeAName]
	ans.NodeBRef = NodeSet[ans.NodeBName]
	FlyLineSet[name] = &ans
	return &ans, nil
}

func (p *NodeOfBlock) Draw(ctx evisiable.Canvas) {
	ctx.DrawPoint(int(p.Pos.X), int(p.Pos.Y), 5)
}

func (p *FlyLine) Draw(ctx evisiable.Canvas) {
	ctx.DrawLine(int(p.NodeARef.Pos.X), int(p.NodeARef.Pos.Y), int(p.NodeBRef.Pos.X), int(p.NodeBRef.Pos.Y), color.RGBA{20, 200, 20, 255})
}
