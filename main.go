package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	_ "image/png"
	"runtime"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

type Scene struct {
}

func main() {
	window := initializeGraphics()
	setupInput(window)

	// Procedural generation happens here.
	scene := createScene()

	renderLoop(window, scene)
	glfw.Terminate()
}

func createScene() Scene {
	return Scene{}
}
