package art

import (
	"image"
	"image/draw"
	"sort"

	geo "github.com/clayts/gec/geometry"
	gfx "github.com/clayts/gec/graphics"
	auto "github.com/clayts/gec/graphics/automatic"
	"github.com/clayts/gec/pixels"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Atlas struct {
	volumes []gfx.TextureArray
	entries []entry
	buffers []auto.Buffer
}

func OpenAtlas(buffers int) *Atlas {
	atl := &Atlas{}
	atl.buffers = make([]auto.Buffer, buffers)
	for i := range atl.buffers {
		atl.buffers[i] = openBuffer()
	}

	return atl
}

func (atl *Atlas) Pack() {
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
	indices := make([]int, len(atl.entries))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		ib := atl.entries[i].Bounds()
		jb := atl.entries[j].Bounds()
		return ib.Dx()*ib.Dy() > jb.Dx()*jb.Dy()
	})

	// Arrange
	for _, index := range indices {
		ent := atl.entries[index]
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
			atl.entries[index] = ent
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
	for _, ent := range atl.entries {
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
			atl.volumes = append(atl.volumes, gfx.OpenTextureArray(rgbaPages, false, false, false, false))
		}
	}

	// DEBUG - save atlas
	// for i, rgbaPages := range rgbaVolumes {
	// 	for j, rgbaPage := range rgbaPages {
	// 		pixels.SaveImage(fmt.Sprint("volume", i, "page", j, ".png"), rgbaPage)
	// 	}
	// }
}

func (atl *Atlas) Buffers(bufferIndices ...int) struct {
	Clear func()
	Draw  func(camera geo.Transform)
} {
	return struct {
		Clear func()
		Draw  func(camera geo.Transform)
	}{
		Clear: func() {
			for _, bufferIndex := range bufferIndices {
				atl.buffers[bufferIndex].Clear(false, true)
			}
		},
		Draw: func(camera geo.Transform) {
			if len(bufferIndices) > 0 {
				w, h := gfx.Window.GetSize()
				program.SetUniform(program.UniformLocation("screenSize"), [2]float32{float32(w), float32(h)})

				program.SetUniform(program.UniformLocation("cameraTransform"), camera.Inverse())

				textureUnits := make([]gfx.TextureUnit, len(atl.volumes))
				for i, volume := range atl.volumes {
					u := gfx.TextureUnit(i)
					u.SetTextureArray(volume)
					textureUnits[i] = u
				}
				program.SetUniform(program.UniformLocation("textureArray"), textureUnits)
				for _, bufferIndex := range bufferIndices {
					atl.buffers[bufferIndex].Draw(program)
				}
			}
		},
	}
}

func (atl *Atlas) Close() {
	for _, volume := range atl.volumes {
		volume.Close()
	}
	atl.volumes = atl.volumes[:0]

	for i := range atl.entries {
		atl.entries[i] = entry{}
	}
	atl.entries = atl.entries[:0]
	for i, buf := range atl.buffers {
		buf.Close()
		atl.buffers[i] = auto.Buffer{}
	}
	atl.buffers = atl.buffers[:0]
}
