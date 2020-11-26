package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/render"
	"github.com/mlesniak/opengl/scene"
	_ "image/png"
)

func main() {
	s := Scene1()
	s.Add(Plane())

	r := render.New(1000, 600)
	r.Render(s)
	r.Exit()

}

func Plane() *scene.Entity {
	m := mgl32.HomogRotate3DX(mgl32.DegToRad(-90))
	m = m.Mul4(mgl32.Scale3D(20, 20, 1))
	m = m.Mul4(mgl32.Translate3D(-0.5, -0.5, 0.0))

	return &scene.Entity{
		Vertices:   model.Plane,
		Model:      m,
		WithNormal: true,
		Color:      mgl32.Vec3{0.29, 0.29, 0.29},
	}
}
