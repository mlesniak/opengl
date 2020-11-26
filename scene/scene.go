package scene

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Entity struct {
	Parent   *Entity
	Children []*Entity

	Name      string
	Parameter map[string]float32

	Vertices   []float32
	Model      mgl32.Mat4
	Color      mgl32.Vec3
	WithNormal bool
}

type Scene struct {
	Seed int64

	Entities []*Entity
}

func New(seed int64) *Scene {
	return &Scene{
		Seed: seed,
	}
}

func (e *Entity) Add(c *Entity) {
	c.Parent = e
	e.Children = append(e.Children, c)
}

func (s *Scene) Add(e *Entity) {
	if e.Parameter == nil {
		e.Parameter = make(map[string]float32)
	}

	s.Entities = append(s.Entities, e)
}
