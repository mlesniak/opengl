package main

import (
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
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
	program := initOpenGL()

	// Translation of C Code glfwSetInputMode(window, GLFW_STICKY_KEYS, GL_TRUE)
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	for !window.ShouldClose() {
		// Input handling
		key := window.GetKey(glfw.KeyEscape)
		if key == glfw.Press {
			break
		}

		// Render loop.
		t := time.Now()
		draw(window, program)
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func draw(window *glfw.Window, program uint32) {
	glfw.PollEvents()
	window.SwapBuffers()
}
