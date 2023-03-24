package graphics

type AttributeLocation uint32

func (a AttributeLocation) GL() uint32 { return uint32(a) }
