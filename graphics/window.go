package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var Window *glfw.Window

func initWindow(title string) {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.False)

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	w, err := glfw.CreateWindow(mode.Width, mode.Height, title, monitor, nil)
	if err != nil {
		panic(err)
	}
	w.MakeContextCurrent()

	w.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	Window = w

}
