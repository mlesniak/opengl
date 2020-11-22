package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"math"
)

func (r *render) setupInput() {
	r.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	lastX := float64(r.width / 2)
	lastY := float64(r.height / 2)
	initialMove := true
	r.window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		r.processMouseInput(&initialMove, &lastX, &lastY, xpos, ypos)
	})

	r.window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		r.processMouseWheelInput(yoff)
	})
}

func (r *render) processMouseWheelInput(yoff float64) {
	camSpeed := yoff * 0.2
	r.camera.position = r.camera.position.Add(r.camera.up.Mul(float32(camSpeed)))
	if r.camera.position.Y() <= 0.1 {
		r.camera.position = mgl32.Vec3{r.camera.position.X(), 0.1, r.camera.position.Z()}
	}
}

func (r *render) processMouseInput(initialMove *bool, lastX *float64, lastY *float64, xpos float64, ypos float64) {
	if *initialMove {
		*lastX = xpos
		*lastY = ypos
		*initialMove = false
	}

	xoffset := xpos - *lastX
	yoffset := ypos - *lastY
	*lastX = xpos
	*lastY = ypos
	sensitivity := 0.1
	xoffset = xoffset * sensitivity
	yoffset = yoffset * sensitivity

	r.camera.yaw += float32(xoffset)
	r.camera.pitch -= float32(yoffset)
	if r.camera.pitch > 0.1 {
		r.camera.pitch = 0.1
	}
	if r.camera.pitch < -89 {
		r.camera.pitch = -89
	}
	v3 := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(r.camera.yaw))) * math.Cos(float64(mgl32.DegToRad(r.camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(r.camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(r.camera.yaw))) * math.Cos(float64(mgl32.DegToRad(r.camera.pitch)))),
	}
	r.camera.front = v3.Normalize()
}

func (r *render) processKeyboardInput(window *glfw.Window, deltaTime float32, cam *camera) {
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
