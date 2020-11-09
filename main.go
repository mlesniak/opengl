package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	width  = 1000
	height = 1000
	fps    = 60
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	initOpenGL()

	// Translation of C Code glfwSetInputMode(window, GLFW_STICKY_KEYS, GL_TRUE)
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	for !window.ShouldClose() {
		// Input handling
		key := window.GetKey(glfw.KeyEscape)
		if key == glfw.Press {
			break
		}

		gl.Clear(gl.COLOR_BUFFER_BIT)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func draw(window *glfw.Window, program uint32) {
	glfw.PollEvents()
	window.SwapBuffers()
}
