package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	"time"
)

// TODO(mlesniak) Expect center automatically?

func Scene1() *scene.Scene {
	s := scene.New(time.Now().UnixNano())
	base := Base()
	base.Add(Top(base))

	s.Add(base)
	return s
}

func Top(parent *scene.Entity) *scene.Entity {
	// TODO(mlesniak) render this
	params := make(map[string]float32)

	m := mgl32.Ident4().
		Mul4(mgl32.Translate3D(
			parent.Parameter["center.x"],
			parent.Parameter["center.y"]+2,
			parent.Parameter["center.z"]))

	return &scene.Entity{
		Name:      "Base",
		Parameter: params,

		Model:      m,
		Vertices:   model.CubeVertices,
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}

func Base() *scene.Entity {
	params := make(map[string]float32)

	m := mgl32.Ident4().Mul4(mgl32.Translate3D(-3, 1, 0))

	originalCenter := mgl32.Vec4{0, 0, 0, 1}
	c := m.Mul4x1(originalCenter)
	params["center.x"] = c.X()
	params["center.y"] = c.Y()
	params["center.z"] = c.Z()

	return &scene.Entity{
		Name:      "Base",
		Parameter: params,

		Model:      m,
		Vertices:   model.CubeVertices,
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}
