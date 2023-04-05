package sprites

import (
	"image"
	"image/draw"
	"sort"

	geo "github.com/clayts/gec/geometry"
	gfx "github.com/clayts/gec/graphics"
	ren "github.com/clayts/gec/renderer"
)

type Sheet struct {
	renderer     *ren.InstanceRenderer
	textureArray gfx.TextureArray
	sources      []struct {
		location [3]float32
		size     [2]float32
		offset   geo.Vector
		image    image.Image
	}
}

func NewSheet() *Sheet {
	sh := &Sheet{}
	return sh
}

func (sh *Sheet) Clear() {
	if sh.renderer == nil {
		return
	}
	sh.renderer.ClearInstances()
}

func (sh *Sheet) Render(camera geo.Transform) {
	if sh.renderer == nil {
		return
	}

	w, h := gfx.Window.GetSize()
	program.SetUniform(screenSizeUniformLocation, [2]float32{float32(w), float32(h)})

	inverse := camera.Inverse()
	program.SetUniform(cameraTransformUniformLocation, [2][3]float32{
		{float32(inverse[0][0]), float32(inverse[0][1]), float32(inverse[0][2])},
		{float32(inverse[1][0]), float32(inverse[1][1]), float32(inverse[1][2])},
	})

	program.SetUniform(textureArrayUniformLocation, gfx.TextureUnit(0).WithSetTextureArray(sh.textureArray))

	sh.renderer.Render()
}

func (sh *Sheet) Delete() {
	if sh.renderer == nil {
		return
	}
	sh.sources = nil
	sh.renderer.Delete()
	sh.textureArray.Delete()
}

func (sh *Sheet) initialize() {
	if sh.renderer != nil {
		return
	}
	sh.renderer = ren.NewInstanceRenderer(program, gfx.TRIANGLE_STRIP, []string{"position"}, "dstTransform", "dstDepth", "srcLocation", "srcSize")
	sh.renderer.Draw(0, 0)
	sh.renderer.Draw(0, 1)
	sh.renderer.Draw(1, 0)
	sh.renderer.Draw(1, 1)
	sh.pack()
}

func (sh *Sheet) pack() {
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
	locations := make([][3]int, len(sh.sources))

	// Sort data
	boxes := make([]struct {
		index int
		w, h  int
	}, len(sh.sources))
	for i, src := range sh.sources {
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
					spaces[i] = space
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
					spaces[i] = space
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
	for i, src := range sh.sources {
		l := locations[i]
		sh.sources[i].location = [3]float32{float32(l[0]), float32(l[1]), float32(l[2])}
		rgba := rgbas[l[2]]
		// Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
		// Draw aligns r.Min in dst with sp in src and then replaces the rectangle r in dst
		draw.Draw(
			rgba, // dst
			image.Rect(l[0], l[1], l[0]+src.image.Bounds().Dx(), l[1]+src.image.Bounds().Dy()), // r
			src.image,              // src
			src.image.Bounds().Min, // sp
			draw.Src,               // op
		)
	}

	// // DEBUG ONLY
	// for _, target := range rgbas {
	// 	f, err := os.Create(strconv.Itoa(rand.Int()) + ".jpg")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	if err = jpeg.Encode(f, target, nil); err != nil {
	// 		log.Printf("failed to encode: %v", err)
	// 	}

	// }

	// Create TextureArray
	sh.textureArray = gfx.NewTextureArray(rgbas, false, false, false, false)
}
