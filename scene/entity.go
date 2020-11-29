package scene

import "github.com/go-gl/mathgl/mgl32"

type Entity struct {
	Parent   *Entity
	Children []*Entity

	Name      string
	Parameter map[string]interface{}

	Vertices   []float32
	Model      mgl32.Mat4
	Color      mgl32.Vec3
	WithNormal bool
}

func (e *Entity) Add(c *Entity) {
	c.Parent = e
	e.Children = append(e.Children, c)
}

func (e *Entity) AddFloat32(name string, value float32) {
	e.Parameter[name] = value
}
