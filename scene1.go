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
	base := Base()
	s.Add(base)

	s.Add(addTower(base))

	return s
}

func addTower(parent *scene.Entity) *scene.Entity {
	params := make(map[string]interface{})

	sx := rand.Float32() * 5
	sy := rand.Float32() * 5
	sz := rand.Float32() * 5
	params["scale"] = mgl32.Vec3{sz, sy, sz}
	scaleMatrix := mgl32.Scale3D(sx, sy, sz)

	pc := parent.Parameter["center"].(mgl32.Vec4)
	px := pc.X()
	py := pc.Y() + 0
	pz := pc.Z()
	params["position"] = mgl32.Vec3{px, py, pz}
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
		Color:      mgl32.Vec3{1, 1, 0},
	}
}

func Base() *scene.Entity {
	params := make(map[string]interface{})

	sx := rand.Float32() * 5
	sy := rand.Float32() * 5
	sz := rand.Float32() * 5
	params["scale"] = mgl32.Vec3{sz, sy, sz}
	scaleMatrix := mgl32.Scale3D(sx, sy, sz)

	px := float32(0)
	py := float32(0.5)
	pz := float32(0)
	params["position"] = mgl32.Vec3{px, py, pz}
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
