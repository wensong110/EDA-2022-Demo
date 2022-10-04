package emodule

var NodeFindPoint map[*NodeOfBlock]*ImportantPoint

func init() {
	NodeFindPoint = make(map[*NodeOfBlock]*ImportantPoint)
}
