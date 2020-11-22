package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"math"
)

var lastX = float64(windowWidth / 2)
var lastY = float64(windowHeight / 2)

func setupInput(window *glfw.Window, cam *camera) {
	initialMove := true
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		processMouseInput(&initialMove, cam, xpos, ypos)
	})

	window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		processMouseWheelInput(cam, yoff)
	})
}

func processMouseWheelInput(cam *camera, yoff float64) {
	camSpeed := yoff * 0.2
	cam.position = cam.position.Add(cam.up.Mul(float32(camSpeed)))

	if cam.position.Y() <= 0.1 {
		cam.position = mgl32.Vec3{cam.position.X(), 0.1, cam.position.Z()}
	}
}

func processMouseInput(initialMove *bool, cam *camera, xpos float64, ypos float64) {
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

	cam.yaw += float32(xoffset)
	cam.pitch -= float32(yoffset)
	if cam.pitch > 0.1 {
		cam.pitch = 0.1
	}
	if cam.pitch < -89 {
		cam.pitch = -89
	}
	v3 := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(cam.yaw))) * math.Cos(float64(mgl32.DegToRad(cam.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.yaw))) * math.Cos(float64(mgl32.DegToRad(cam.pitch)))),
	}
	cam.front = v3.Normalize()
}

func processKeyboardInput(window *glfw.Window, deltaTime float32, cam *camera) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	var camSpeed = 10 * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cam.position = cam.position.Add(cam.front.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cam.position = cam.position.Sub(cam.front.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		cam.position = cam.position.Add(cam.up.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		cam.position = cam.position.Sub(cam.up.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		x := cam.front.Cross(cam.up).Normalize().Mul(camSpeed)
		cam.position = cam.position.Sub(x)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		x := cam.front.Cross(cam.up).Normalize().Mul(camSpeed)
		cam.position = cam.position.Add(x)
	}

	if cam.position.Y() <= 0.1 {
		cam.position = mgl32.Vec3{cam.position.X(), 0.1, cam.position.Z()}
	}
}
