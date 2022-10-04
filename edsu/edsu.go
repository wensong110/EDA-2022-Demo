package edsu

import "fmt"

type Any interface{}

type DSU struct {
	Mp map[Any]Any
}

func NewDSU() *DSU {
	ans := DSU{}
	ans.Mp = make(map[Any]Any)
	return &ans
}

func (p *DSU) AddItem(item Any) error {
	_, hasItem := p.Mp[item]
	if hasItem {
		return fmt.Errorf("there has is the item %v", item)
	}
	p.Mp[item] = item
	return nil
}

func (p *DSU) FindLeader(item Any) (Any, error) {
	_, hasItem := p.Mp[item]
	if !hasItem {
		return nil, fmt.Errorf("there hasn't is the item %v", item)
	}
	if p.Mp[item] == item {
		return item, nil
	}
	ret1, ret2 := p.FindLeader(p.Mp[item])
	p.Mp[item] = ret1
	return ret1, ret2
}

func (p *DSU) Merge(item1, item2 Any) error {
	f1, err := p.FindLeader(item1)
	if err != nil {
		return err
	}
	f2, err := p.FindLeader(item2)
	if err != nil {
		return err
	}
	p.Mp[f1] = f2
	return nil
}
