package atlas

import (
	_ "embed"
	"image"
	"image/draw"
	"sort"

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
	entries []entry
)

type entry struct {
	image.Image
	volume, page, x, y, w, h float32
}

func Open() {
	shaders = []gfx.Shader{
		gfx.OpenVertexShader(vertexShaderSource),
		gfx.OpenFragmentShader(fragmentShaderSource),
	}
	program = gfx.OpenProgram(shaders...)
}

func Pack() {
	// Know Volume limit - hardcoded into fragment shader, the minimum required to be supported by all compatible hardware
	const maxVolumeCount = 16

	// Get page limits
	maxPageCount := gfx.Query(gl.MAX_ARRAY_TEXTURE_LAYERS)
	maxPageSize := gfx.Query(gl.MAX_TEXTURE_SIZE)

	// Make list of free spaces
	type space struct {
		x, y, w, h, page, volume int
	}

	spaces := []space{}

	for volume := 0; volume < maxVolumeCount; volume++ {
		for page := 0; page < maxPageCount; page++ {
			spaces = append(spaces, space{0, 0, maxPageSize, maxPageSize, page, volume})
		}
	}

	// Make list of volume sizes
	volumeSizes := [maxVolumeCount]struct {
		width     int
		height    int
		pageCount int
	}{}

	// Get list of entry indices, sorted by size
	indices := make([]int, len(entries))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		ib := entries[i].Bounds()
		jb := entries[j].Bounds()
		return ib.Dx()*ib.Dy() > jb.Dx()*jb.Dy()
	})

	// Arrange
	for _, index := range indices {
		ent := entries[index]
		packed := false

		b := ent.Bounds()
		// Check each space, from last to first
		for i, s := range spaces {
			if b.Dx() <= s.w && b.Dy() <= s.h {
				// Found the space; add the box to its top-left corner
				// |-------?-------|
				// |  box  ?       |
				// ?????????       |
				// |         space |
				// |_______________|
				ent.volume = float32(s.volume)
				ent.page = float32(s.page)
				ent.x = float32(s.x)
				ent.y = float32(s.y)
				ent.w = float32(b.Dx())
				ent.h = float32(b.Dy())
				packed = true

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
					spaces = append([]space{
						{s.x + b.Dx(), s.y, s.w - b.Dx(), s.h, s.page, s.volume},
					}, spaces...)
				} else if b.Dx() == s.w {
					// Space matches the box width
					// |---------------|
					// |      box      |
					// |_______________|
					// |     new space |
					// |_______________|
					spaces = append([]space{
						{s.x, s.y + b.Dy(), s.w, s.h - b.Dy(), s.page, s.volume},
					}, spaces...)
				} else {
					// Otherwise the box splits the space into two spaces
					// |-------|-----------|
					// |  box  | new space |
					// |_______|___________|
					// |     new space     |
					// |___________________|
					spaces = append([]space{
						{s.x, s.y + b.Dy(), s.w, s.h - b.Dy(), s.page, s.volume},
						{s.x + b.Dx(), s.y, s.w - b.Dx(), b.Dy(), s.page, s.volume},
					}, spaces...)
				}
				break
			}
		}
		if packed {
			entries[index] = ent
		} else {
			panic("insufficient storage capacity")
		}
	}

	// Create RGBAs
	rgbaVolumes := [maxVolumeCount][]*image.RGBA{}
	for volume, volumeSize := range volumeSizes {
		rgbaVolumes[volume] = make([]*image.RGBA, volumeSize.pageCount)
		for page := range rgbaVolumes[volume] {
			rgbaVolumes[volume][page] = image.NewRGBA(image.Rect(0, 0, volumeSize.width, volumeSize.height))
		}
	}

	// Copy data
	for _, ent := range entries {
		rgba := rgbaVolumes[int(ent.volume)][int(ent.page)]
		// Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
		// Draw aligns r.Min in dst with sp in src and then replaces the rectangle r in dst
		draw.Draw(
			rgba, // dst
			image.Rect(int(ent.x), int(ent.y), int(ent.x+ent.w), int(ent.y+ent.h)), // r
			pixels.FlipY(ent), // src
			ent.Bounds().Min,  // sp
			draw.Src,          // op
		)
	}

	// Create TextureArrays
	for _, rgbaPages := range rgbaVolumes {
		if len(rgbaPages) > 0 {
			volumes = append(volumes, gfx.OpenTextureArray(rgbaPages, false, false, false, false))
		}
	}

	// DEBUG - save atlas
	// for i, rgbaPages := range rgbaVolumes {
	// 	for j, rgbaPage := range rgbaPages {
	// 		pixels.SaveImage(fmt.Sprint("volume", i, "page", j, ".png"), rgbaPage)
	// 	}
	// }
}

func Close() {
	program.Close()

	for _, shader := range shaders {
		shader.Close()
	}
	shaders = shaders[:0]

	for _, volume := range volumes {
		volume.Close()
	}
	volumes = volumes[:0]

	for i := range entries {
		entries[i] = entry{}
	}
	entries = entries[:0]
}
