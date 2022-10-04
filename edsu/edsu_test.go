package edsu_test

import (
	"Solution/edsu"
	"testing"
)

func TestDSU(t *testing.T) {
	dsu := edsu.NewDSU()
	dsu.AddItem(1)
	dsu.AddItem(2)
	dsu.AddItem(3)
	t.Log(dsu.Mp)
	dsu.Merge(1, 2)
	t.Log(dsu.Mp)
	dsu.Merge(1, 3)
	t.Log(dsu.Mp)
	t.Log(dsu.FindLeader(1))
	t.Log(dsu.FindLeader(2))
	t.Log(dsu.FindLeader(3))
}
