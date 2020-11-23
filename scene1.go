package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	"math/rand"
)

func Scene1() *scene.Scene {
	s := scene.New()
	s.Add(Cube())
	s.Add(Cube2())
	s.Add(Random(20))
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

func Random(triangles int) *scene.Entity {
	vs := make([]float32, triangles*9+9)
	for i := 0; i < triangles*9; i++ {
		vs[i+0] = rand.Float32() * 10
		vs[i+1] = rand.Float32() * 10
		vs[i+2] = rand.Float32() * 10
		vs[i+3] = rand.Float32() * 10
		vs[i+4] = rand.Float32() * 10
		vs[i+5] = rand.Float32() * 10
		vs[i+6] = rand.Float32() * 10
		vs[i+7] = rand.Float32() * 10
		vs[i+8] = rand.Float32() * 10
	}

	return &scene.Entity{
		Vertices:   vs,
		Position:   mgl32.Translate3D(-10, 0, -10),
		WithNormal: false,
		Color:      mgl32.Vec3{1, 1, 0},
	}
}
