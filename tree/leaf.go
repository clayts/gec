package tree

import geo "github.com/clayts/gec/geometry"

type Leaf struct {
	shape    geo.Shape
	tree     *Tree
	branch   *branch
	index    int
	Contents interface{}
}

func (t *Tree) NewLeaf() *Leaf {
	return &Leaf{tree: t}
}

func (l *Leaf) Enable() *Leaf {
	if l.branch == nil {
		l.tree.add(l)
	}
	return l
}

func (l *Leaf) Disable() *Leaf {
	if l.branch != nil {
		l.tree.remove(l)
	}
	return l
}

func (l *Leaf) SetShape(s geo.Shape) *Leaf {
	l.shape = s
	if l.branch != nil {
		l.tree.remove(l)
		l.tree.add(l)
	}
	return l
}

func (l *Leaf) SetContents(c interface{}) *Leaf {
	l.Contents = c
	return l
}
