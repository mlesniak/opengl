package scene

import "github.com/go-gl/mathgl/mgl32"

type Entity struct {
	Vertices []float32
	Position mgl32.Mat4
	Color    mgl32.Vec3

	// TODO(mlesniak) Not sure if viable?
	WithNormal bool
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
