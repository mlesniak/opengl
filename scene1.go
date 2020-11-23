package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	"math/rand"
	"time"
)

// TODO(mlesniak) How can I use scenes to create other scenes?
// TODO(mlesniak) Add a scene means applying the matrix to everything?

func Scene1() *scene.Scene {
	rand.Seed(time.Now().UnixNano())
	s := scene.New()
	s.Add(Base())
	s.Add(Cross())
	return s
}

func Base() *scene.Entity {
	return &scene.Entity{
		Vertices: model.CubeVertices,
		Position: mgl32.Ident4().
			Mul4(mgl32.Scale3D(4, 1, 1)).
			Mul4(mgl32.Translate3D(0, 0.5, 0)),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}

func Cross() *scene.Entity {
	return &scene.Entity{
		Vertices: model.CubeVertices,
		Position: mgl32.Ident4().
			Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-90))).
			Mul4(mgl32.Scale3D(4, 1, 1)).
			Mul4(mgl32.Translate3D(0, 0.5, -1)),
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}
