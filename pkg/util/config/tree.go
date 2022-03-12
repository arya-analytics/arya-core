package config

const LeafSeparator = "."

type Tree struct {
	Base string
}

func (t Tree) Leaf(name string) string {
	return t.Base + LeafSeparator + name
}

func (t Tree) SubTree(base string) Tree {
	return Tree{Base: t.Leaf(base)}
}
