package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	_ "image/png"
	"runtime"
)

// TODO(mlesniak) add plane
// TODO(mlesniak) add diffuse lightning

const windowWidth = 800
const windowHeight = 600

// TODO(mlesniak) show plane first
// TODO(mlesniak) struct
// TODO(mlesniak) General refactoring for render loop based on objects

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func main() {
	window := initializeWindow()
	render(window)
	glfw.Terminate()
}
