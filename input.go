package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"math"
)

func processMouse(initialMove *bool, xpos float64, ypos float64) {
	if *initialMove {
		lastX = xpos
		lastY = ypos
		*initialMove = false
	}

	xoffset := xpos - lastX
	yoffset := ypos - lastY
	lastX = xpos
	lastY = ypos
	sensitivity := 0.1
	xoffset = xoffset * sensitivity
	yoffset = yoffset * sensitivity

	yaw += float32(xoffset)
	pitch -= float32(yoffset)
	if pitch > 0.1 {
		pitch = 0.1
	}
	if pitch < -89 {
		pitch = -89
	}
	v3 := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
	}
	camFront = v3.Normalize()
}

func processKeyboardInput(window *glfw.Window, deltaTime float32) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	var camSpeed = 10 * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camPos = camPos.Add(camFront.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camPos = camPos.Sub(camFront.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		camPos = camPos.Add(camUp.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		camPos = camPos.Sub(camUp.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		x := camFront.Cross(camUp).Normalize().Mul(camSpeed)
		camPos = camPos.Sub(x)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		x := camFront.Cross(camUp).Normalize().Mul(camSpeed)
		camPos = camPos.Add(x)
	}

	if camPos.Y() <= 0.1 {
		camPos = mgl32.Vec3{camPos.X(), 0.1, camPos.Z()}
	}
}
