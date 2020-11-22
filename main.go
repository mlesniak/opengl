package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	models "github.com/mlesniak/opengl/model"
	_ "image/png"
	"runtime"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

type Entity struct {
	Vertices   []float32
	Position   mgl32.Mat4
	WithNormal bool
	Color      mgl32.Vec3
}

type Scene struct {
	Entities []*Entity
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Add(e *Entity) {
	s.Entities = append(s.Entities, e)
}

func main() {
	// Create vertices.
	scene := New()
	scene.Add(Cube())
	scene.Add(Plane())

	// Display them.
	window := initializeGraphics()
	renderLoop(window, scene)
	glfw.Terminate()
}

func createScene() Scene {
	return Scene{}
}

func Cube() *Entity {
	return &Entity{
		Vertices:   models.CubeVertices,
		Position:   mgl32.Translate3D(+0.5, +0.5, 0),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}

func Plane() *Entity {
	model := mgl32.HomogRotate3DX(mgl32.DegToRad(-90))
	model = model.Mul4(mgl32.Scale3D(20, 20, 1))
	model = model.Mul4(mgl32.Translate3D(-0.5, -0.5, 0.0))

	return &Entity{
		Vertices:   models.Plane,
		Position:   model,
		WithNormal: true,
		Color:      mgl32.Vec3{0.29, 0.29, 0.29},
	}
}
