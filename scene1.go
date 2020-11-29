package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/scene"
	"math/rand"
	"time"
)

func Scene1() *scene.Scene {
	seed := time.Now().UnixNano()
	s := scene.New(seed)
	s.Add(Base())
	return s
}

func Base() *scene.Entity {
	params := make(map[string]interface{})

	params["scale.x"] = rand.Float32() * 5
	params["scale.y"] = rand.Float32() * 5
	params["scale.z"] = rand.Float32() * 5
	scaleMatrix := mgl32.Scale3D(
		params["scale.x"].(float32),
		params["scale.y"].(float32),
		params["scale.z"].(float32))

	px := float32(0)
	py := float32(0.5)
	pz := float32(0)
	params["position.x"] = px
	params["position.y"] = py
	params["position.z"] = pz
	positionMatrix := mgl32.Translate3D(px, py, pz)

	center := mgl32.Vec4{0, 0, 0, 1}
	m := mgl32.Ident4().
		Mul4(scaleMatrix).
		Mul4(positionMatrix)

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
