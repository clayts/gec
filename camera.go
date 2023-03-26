package gec

import (
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
)

type Camera struct {
	Transform geometry.Transform
}

func (c Camera) Shape() geometry.Shape {
	width, height := graphics.Window.GetSize()
	return c.Transform.Rectangle(geometry.R(geometry.V(0, 0), geometry.V(float64(width), float64(height))))
}
