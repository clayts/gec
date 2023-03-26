package graphics

import (
	"runtime"
	"time"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func Initialize(title string) {
	runtime.LockOSThread()

	// GLFW
	initWindow(title)

	// GL
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// wait for window to be properly initialized
	start := time.Now()
	for w := 1; w == 1; w, _ = Window.GetSize() {
		if time.Since(start).Seconds() > 1 {
			panic("could not initialize window")
		}
	}
}

func Clear(color, depth, stencil bool) {
	if color {
		if depth {
			if stencil {
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
			} else {
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			}
		} else {
			if stencil {
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
			} else {
				gl.Clear(gl.COLOR_BUFFER_BIT)
			}
		}
	} else {
		if depth {
			if stencil {
				gl.Clear(gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
			} else {
				gl.Clear(gl.DEPTH_BUFFER_BIT)
			}
		} else {
			if stencil {
				gl.Clear(gl.STENCIL_BUFFER_BIT)
			} else {
				return
			}
		}
	}
}

func Render() {
	glfw.PollEvents()
	Window.SwapBuffers()
}

func Delete() {
	glfw.Terminate()
}

func MaximumTextureArrayLayers() int {
	var v int32
	gl.GetIntegerv(gl.MAX_ARRAY_TEXTURE_LAYERS, &v)
	return int(v)
}

func MaximumTextureSize() int {
	var v int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &v)
	return int(v)
}

func PanicOnError() {
	e := gl.GetError()
	switch e {
	case gl.INVALID_ENUM:
		panic("invalid enum")
	case gl.INVALID_VALUE:
		panic("invalid value")
	case gl.INVALID_OPERATION:
		panic("invalid operation")
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		panic("invalid framebuffer operation")
	case gl.OUT_OF_MEMORY:
		panic("out of memory")
	case gl.STACK_UNDERFLOW:
		panic("stack underflow")
	case gl.STACK_OVERFLOW:
		panic("stack overflow")
	}
}
