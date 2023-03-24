package graphics

type Attribute struct {
	name     string
	location AttributeLocation
	l, m, n  int32
}

func (a Attribute) Name() string                { return a.name }
func (a Attribute) Location() AttributeLocation { return a.location }
func (a Attribute) Size() (int32, int32, int32) { return a.l, a.m, a.n }
