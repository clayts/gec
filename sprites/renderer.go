package sprites

import (
	"image"
	"image/draw"
	"sort"

	gfx "github.com/clayts/gec/graphics"
	ren "github.com/clayts/gec/renderer"
)

type Renderer struct {
	renderer     *ren.InstanceRenderer
	textureArray gfx.TextureArray
	sources      []struct {
		location [3]float32
		size     [2]float32
		image    image.Image
	}
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	return r
}

func (r *Renderer) Clear() {
	if r.renderer == nil {
		return
	}
	r.renderer.ClearInstances()
}

func (r *Renderer) Render() {
	if r.renderer == nil {
		return
	}
	program.SetUniform(textureArrayUniformLocation, gfx.TextureUnit(0).WithSetTextureArray(r.textureArray))
	w, h := gfx.Window.GetSize()
	program.SetUniform(screenSizeUniformLocation, [2]float32{float32(w), float32(h)})

	r.renderer.Render()
}

func (r *Renderer) Delete() {
	if r.renderer == nil {
		return
	}
	r.sources = nil
	r.renderer.Delete()
	r.textureArray.Delete()
}

func (r *Renderer) initialize() {
	if r.renderer != nil {
		return
	}
	r.renderer = ren.NewInstanceRenderer(program, gfx.TRIANGLE_STRIP, []string{"position"}, "dstTransform", "dstDepth", "srcLocation", "srcSize")
	r.renderer.Draw(0, 0)
	r.renderer.Draw(0, 1)
	r.renderer.Draw(1, 0)
	r.renderer.Draw(1, 1)
	r.pack()
}

func (r *Renderer) pack() {
	// Make list of free spaces
	spaces := []struct {
		x, y, w, h, z int
	}{}
	maxZ := gfx.MaximumTextureArrayLayers()
	maxWH := gfx.MaximumTextureSize()
	for i := maxZ - 1; i >= 0; i-- {
		spaces = append(spaces, struct{ x, y, w, h, z int }{0, 0, maxWH, maxWH, i})
	}

	// Make list of locations
	locations := make([][3]int, len(r.sources))

	// Sort data
	boxes := make([]struct {
		index int
		w, h  int
	}, len(r.sources))
	for i, src := range r.sources {
		v := boxes[i]
		v.w = src.image.Bounds().Dx()
		v.h = src.image.Bounds().Dy()
		v.index = i
		boxes[i] = v
	}
	sort.Slice(boxes, func(i, j int) bool {
		di := boxes[i]
		dj := boxes[j]
		return di.w*di.h > dj.w*dj.h
	})

	width, height, depth := 0, 0, 0

	// Determine locations
	for _, box := range boxes {
		for i := len(spaces) - 1; i >= -1; i-- {
			if i == -1 {
				panic("maximum texture size exceeded")
			}
			space := spaces[i]
			if box.w <= space.w && box.h <= space.h {
				// found the space; add the box to its top-left corner
				// |-------?-------|
				// |  box  ?       |
				// ?????????       |
				// |         space |
				// |_______________|
				l := [3]int{space.x, space.y, space.z}
				locations[box.index] = l
				if width < space.x+box.w {
					width = space.x + box.w
				}
				if height < space.y+box.h {
					height = space.y + box.h
				}
				if depth < space.z+1 {
					depth = space.z + 1
				}
				if box.w == space.w && box.h == space.h {
					// space matches the box exactly; remove it
					// |---------------|
					// |  box          |
					// |               |
					// |         space |
					// |_______________|
					spaces = append(spaces[:i], spaces[i+1:]...)
				} else if box.h == space.h {
					// space matches the box height; update it accordingly
					// |-------|---------------|
					// |  box  | updated space |
					// |_______|_______________|
					space.x += box.w
					space.w -= box.w
					spaces[i] = space
				} else if box.w == space.w {
					// space matches the box width; update it accordingly
					// |---------------|
					// |      box      |
					// |_______________|
					// | updated space |
					// |_______________|
					space.y += box.h
					space.h -= box.h
				} else {
					// otherwise the box splits the space into two spaces
					// |-------|-----------|
					// |  box  | new space |
					// |_______|___________|
					// | updated space     |
					// |___________________|
					spaces = append(spaces, struct{ x, y, w, h, z int }{space.x + box.w, space.y, space.w - box.w, box.h, space.z})
					space.y += box.h
					space.h -= box.h
				}
				break
			}
		}
	}

	// Create RGBAs
	rgbas := make([]*image.RGBA, depth)
	for i := range rgbas {
		rgbas[i] = image.NewRGBA(image.Rect(0, 0, width, height))
	}

	// Copy data to RGBAs and update sources
	for i, src := range r.sources {
		l := locations[i]
		r.sources[i].location = [3]float32{float32(l[0]), float32(l[1]), float32(l[2])}
		rgba := rgbas[l[2]]
		draw.Draw(rgba, image.Rect(l[0], l[1], l[0]+src.image.Bounds().Dx(), l[1]+src.image.Bounds().Dy()), src.image, image.Point{}, draw.Src)
	}

	// // DEBUG ONLY
	// for i, target := range rgbas {
	// 	f, err := os.Create(strconv.Itoa(i) + ".jpg")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	if err = jpeg.Encode(f, target, nil); err != nil {
	// 		log.Printf("failed to encode: %v", err)
	// 	}

	// }

	// Create TextureArray
	r.textureArray = gfx.NewTextureArray(rgbas, false, false, false, false)
}
