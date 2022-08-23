package emodule

import (
	"bufio"
	"fmt"
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

func init() {
	BlockSet = make(map[string]*Block)
	InstanceSet = make(map[string]*Instance)
}

type Block struct {
	Name     string
	PointSet []Point
}

func NewBlock(name string, str string) *Block {
	ans := Block{}
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
	BlockSet[ans.Name] = &ans
	return &ans
}

func (p *Block) Rotate(how string) *Block {
	//TODO
	return nil
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
	if len(strs) < 3 {
		return nil, fmt.Errorf("wrong instance format")
	}
	blockName := strs[0]
	_, hasBlock := BlockSet[blockName]
	if !hasBlock {
		return nil, fmt.Errorf("don't has block %s", blockName)
	}
	ans.Block = BlockSet[blockName]
	offsetPointStr := strs[1]
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
		return nil, fmt.Errorf("wrong instance format")
	}
	ans.Offset = NewPoint(string2float64(points[0]), string2float64(points[0]))
	ans.Rotate = strs[2]
	//struct consturct
	layerNames := strings.Split(ans.Name, "/")
	builder.Reset()
	for i, v := range layerNames {
		if i != 0 {
			builder.WriteRune('/')
		}
		builder.WriteString(v)
		_, hasInstance := BlockSet[builder.String()]
		if hasInstance {
			ans.InstanceChain = append(ans.InstanceChain, InstanceSet[builder.String()])
		}
	}

	return &ans, nil
}
