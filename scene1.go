package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	"math/rand"
	"time"
)

func Scene1() *scene.Scene {
	s := scene.New(time.Now().UnixNano())

	base := Base()

	s.Add(base)
	return s
}

func Base() *scene.Entity {
	params := make(map[string]interface{})

	x := rand.Float32() * 5
	y := rand.Float32() * 5
	z := rand.Float32() * 5
	params["scale.x"] = x
	params["scale.y"] = y
	params["scale.z"] = z
	scaleMatrix := mgl32.Scale3D(x, y, z)

	center := mgl32.Vec4{0, 0, 0, 1}
	m := mgl32.Ident4().
		Mul4(scaleMatrix).
		Mul4(mgl32.Translate3D(0, 0.5, 0))

	params["center"] = m.Mul4x1(center)

	return &scene.Entity{
		Name:       "Base",
		Parameter:  params,
		Model:      m,
		Vertices:   model.CubeVertices,
		WithNormal: true,
		Color:      mgl32.Vec3{1, 0, 0},
	}
}
