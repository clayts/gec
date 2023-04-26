package atlas

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"sort"

	geo "github.com/clayts/gec/geometry"
	gfx "github.com/clayts/gec/graphics"
	"github.com/clayts/gec/pixels"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	//go:embed shaders/vertex.glsl
	vertexShaderSource string

	//go:embed shaders/fragment.glsl
	fragmentShaderSource string

	shaders []gfx.Shader
	program gfx.Program
	volumes []gfx.TextureArray
	entries []*entry
)

func Open() {
	shaders = []gfx.Shader{
		gfx.OpenVertexShader(vertexShaderSource),
		gfx.OpenFragmentShader(fragmentShaderSource),
	}
	program = gfx.OpenProgram(shaders...)

	pack()
}

func pack() {
	// Sort images by size
	sort.Slice(entries, func(i, j int) bool {
		ib := entries[i].Bounds()
		jb := entries[j].Bounds()
		return ib.Dx()*ib.Dy() > jb.Dx()*jb.Dy()
	})

	// Get page limits
	maxPageCount := gfx.Query(gl.MAX_ARRAY_TEXTURE_LAYERS)
	maxPageSize := gfx.Query(gl.MAX_TEXTURE_SIZE)

	// Know Volume limit - hardcoded into fragment shader, the minimum required to be supported by all compatible hardware
	const maxVolumeCount = 16

	fmt.Println("texture space:", maxVolumeCount, "x", maxPageCount, "x", "(", maxPageSize, "x", maxPageSize, ")")

	// Make list of free spaces, with the first spaces at the end
	type space struct {
		x, y, w, h, page, volume int
	}

	spaces := []space{}

	for volume := maxVolumeCount - 1; volume >= 0; volume-- {
		for page := maxPageCount - 1; page >= 0; page-- {
			spaces = append(spaces, space{0, 0, maxPageSize, maxPageSize, page, volume})
		}
	}

	// Make list of volume sizes
	volumeSizes := [maxVolumeCount]struct {
		width     int
		height    int
		pageCount int
	}{}

	// Arrange
	for _, img := range entries {
		b := img.Bounds()

		// Check each space, from last to first
		for i := len(spaces) - 1; i >= -1; i-- {
			if i == -1 {
				panic("unable to pack images")
			}
			s := spaces[i]
			if b.Dx() <= s.w && b.Dy() <= s.h {
				// Found the space; add the box to its top-left corner
				// |-------?-------|
				// |  box  ?       |
				// ?????????       |
				// |         space |
				// |_______________|
				img.volume = float32(s.volume)
				img.page = float32(s.page)
				img.x = float32(s.x)
				img.y = float32(s.y)

				// Keep track of volume sizes
				if s.x+b.Dx() > volumeSizes[s.volume].width {
					volumeSizes[s.volume].width = s.x + b.Dx()
				}
				if s.y+b.Dy() > volumeSizes[s.volume].height {
					volumeSizes[s.volume].height = s.y + b.Dy()
				}
				if s.page+1 > volumeSizes[s.volume].pageCount {
					volumeSizes[s.volume].pageCount = s.page + 1
				}

				// Remove the space
				spaces = append(spaces[:i], spaces[i+1:]...)

				if b.Dx() == s.w && b.Dy() == s.h {
					// Space matches the box exactly
					// |---------------|
					// |  box          |
					// |               |
					// |         space |
					// |_______________|
				} else if b.Dy() == s.h {
					// Space matches the box height
					// |-------|---------------|
					// |  box  |     new space |
					// |_______|_______________|
					spaces = append(spaces, space{s.x + b.Dx(), s.y, s.w - b.Dx(), s.h, s.page, s.volume})
				} else if b.Dx() == s.w {
					// Space matches the box width
					// |---------------|
					// |      box      |
					// |_______________|
					// |     new space |
					// |_______________|
					spaces = append(spaces, space{s.x, s.y + b.Dy(), s.w, s.h - b.Dy(), s.page, s.volume})
				} else {
					// Otherwise the box splits the space into two spaces
					// |-------|-----------|
					// |  box  | new space |
					// |_______|___________|
					// |     new space     |
					// |___________________|
					spaces = append(spaces, space{s.x + b.Dx(), s.y, s.w - b.Dx(), b.Dy(), s.page, s.volume})
					spaces = append(spaces, space{s.x, s.y + b.Dy(), s.w, s.h - b.Dy(), s.page, s.volume})
				}
				break
			}
		}
	}

	// Create RGBAs
	rgbas := [maxVolumeCount][]*image.RGBA{}
	for volume, volumeSize := range volumeSizes {
		rgbas[volume] = make([]*image.RGBA, volumeSize.pageCount)
		for page := range rgbas[volume] {
			rgbas[volume][page] = image.NewRGBA(image.Rect(0, 0, volumeSize.width, volumeSize.height))
		}
	}

	// Copy data
	for _, img := range entries {
		rgba := rgbas[int(img.volume)][int(img.page)]
		// Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
		// Draw aligns r.Min in dst with sp in src and then replaces the rectangle r in dst
		draw.Draw(
			rgba, // dst
			image.Rect(int(img.x), int(img.y), int(img.x)+img.Bounds().Dx(), int(img.y)+img.Bounds().Dy()), // r
			pixels.FlipY(img), // src
			img.Bounds().Min,  // sp
			draw.Src,          // op
		)
	}

	// Create TextureArrays
	for _, volumeRGBAs := range rgbas {
		if len(volumeRGBAs) > 0 {
			volumes = append(volumes, gfx.NewTextureArray(volumeRGBAs, false, false, false, false))
		}
	}
}

func Close() {
	program.Close()

	for _, shader := range shaders {
		shader.Delete()
	}
	shaders = shaders[:0]

	for _, volume := range volumes {
		volume.Delete()
	}
	volumes = volumes[:0]

	for i := range entries {
		entries[i] = nil
	}
	entries = entries[:0]
}

type entry struct {
	image.Image
	volume, page, x, y float32
}

type Sprite struct {
	entry     *entry
	region    geo.Rectangle
	Transform geo.Transform
	Depth     float32
}

func MakeSprite(source image.Image) Sprite {
	ent := &entry{Image: source}
	entries = append(entries, ent)
	bounds := source.Bounds()
	return Sprite{
		entry:     ent,
		region:    geo.R(geo.V(float64(bounds.Min.X), float64(bounds.Min.Y)), geo.V(float64(bounds.Max.X), float64(bounds.Max.Y))),
		Transform: geo.T(),
		Depth:     0,
	}
}

func (spr Sprite) WithRegion(region geo.Rectangle) Sprite {
	if !spr.region.Contains(region) {
		panic("region out of bounds")
	}
	spr.region = region
	return spr
}

func (spr Sprite) WithTransform(transform geo.Transform) Sprite {
	spr.Transform = transform
	return spr
}

func (spr Sprite) WithDepth(depth float32) Sprite {
	spr.Depth = depth
	return spr
}

func (spr Sprite) Region() geo.Rectangle { return spr.region }

func (spr Sprite) Image() image.Image { return spr.entry }

func (spr Sprite) Shape() geo.Shape {
	return spr.Transform.Rectangle(geo.R(geo.V(0, 0), spr.region.Size()))
}

func (spr Sprite) Instance() []float32 {
	size := spr.region.Size()
	return []float32{
		float32(spr.Transform[0][0]), float32(spr.Transform[0][1]), float32(spr.Transform[0][2]), float32(spr.Transform[1][0]), float32(spr.Transform[1][1]), float32(spr.Transform[1][2]),
		spr.Depth,
		spr.entry.x + float32(spr.region.Min.X), spr.entry.y + float32(spr.region.Min.Y), spr.entry.page, spr.entry.volume,
		float32(size.X), float32(size.Y),
	}
}

type Renderer struct {
	renderer *gfx.Renderer
}

func OpenRenderer() Renderer {
	ren := Renderer{}
	ren.renderer = gfx.OpenRenderer(
		gfx.TRIANGLE_STRIP,
		[]string{"position"},
		[]string{"dstTransform", "dstDepth", "srcLocation", "srcSize"},
		program,
	)
	ren.renderer.SetVertices(
		0, 0,
		0, 1,
		1, 0,
		1, 1,
	)
	return ren
}

func (ren Renderer) Render(camera geo.Transform) {
	w, h := gfx.Window.GetSize()
	program.SetUniform(program.UniformLocation("screenSize"), [2]float32{float32(w), float32(h)})

	inverse := camera.Inverse()
	program.SetUniform(program.UniformLocation("cameraTransform"), [2][3]float32{
		{float32(inverse[0][0]), float32(inverse[0][1]), float32(inverse[0][2])},
		{float32(inverse[1][0]), float32(inverse[1][1]), float32(inverse[1][2])},
	})

	textureUnits := make([]gfx.TextureUnit, len(volumes))
	for i, volume := range volumes {
		textureUnits[i] = gfx.TextureUnit(i).WithSetTextureArray(volume)
	}
	program.SetUniform(program.UniformLocation("textureArray"), textureUnits)

	ren.renderer.Render()
}

func (ren Renderer) SetInstances(instances ...float32) {
	ren.renderer.SetInstances(instances...)
}

func (ren Renderer) Close() {
	ren.renderer.Close()
}
