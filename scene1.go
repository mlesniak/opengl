package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
)

func Scene1() *scene.Scene {
	s := scene.New()
	s.Add(Cube())
	s.Add(Cube2())
	return s
}

func Cube() *scene.Entity {
	return &scene.Entity{
		Vertices:   model.CubeVertices,
		Position:   mgl32.Translate3D(+0.5, +0.5, 0).Mul4(mgl32.Scale3D(0.2, 5, 1)),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}

func Cube2() *scene.Entity {
	return &scene.Entity{
		Vertices:   model.CubeVertices,
		Position:   mgl32.Translate3D(+1, +0.5, 0),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}
