package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	_ "image/png"
	"runtime"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func main() {
	// Create vertices.
	s := scene.New()
	s.Add(Cube())
	s.Add(Plane())

	// Display them.
	// TODO(mlesniak) Display object with new(windows, scene) instead of call to Render()
	window := initializeGraphics()
	Render(window, s)
	glfw.Terminate()
}

func Cube() *scene.Entity {
	return &scene.Entity{
		Vertices:   model.CubeVertices,
		Position:   mgl32.Translate3D(+0.5, +0.5, 0),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}

func Plane() *scene.Entity {
	m := mgl32.HomogRotate3DX(mgl32.DegToRad(-90))
	m = m.Mul4(mgl32.Scale3D(20, 20, 1))
	m = m.Mul4(mgl32.Translate3D(-0.5, -0.5, 0.0))

	return &scene.Entity{
		Vertices:   model.Plane,
		Position:   m,
		WithNormal: true,
		Color:      mgl32.Vec3{0.29, 0.29, 0.29},
	}
}
