package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	_ "image/png"
	"log"
	"runtime"
	//"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glfw/v3.3/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func main() {
	window := initialize()
	render(window)
	glfw.Terminate()
}

func render(window *glfw.Window) {
	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initialize() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}
