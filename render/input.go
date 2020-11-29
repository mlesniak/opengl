package render

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"math"
)

type input struct {
	mode      uint32
	render    *render
	lastPress float64
}

func setupInput(r *render) *input {
	i := input{
		mode:   gl.FILL,
		render: r,
	}
	i.render.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	lastX := float64(i.render.width / 2)
	lastY := float64(i.render.height / 2)
	initialMove := true
	i.render.window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		i.processMouseInput(&initialMove, &lastX, &lastY, xpos, ypos)
	})

	i.render.window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		i.processMouseWheelInput(yoff)
	})

	return &i
}

func (i *input) processMouseWheelInput(yoff float64) {
	camSpeed := yoff * 0.2
	i.render.camera.position = i.render.camera.position.Add(i.render.camera.up.Mul(float32(camSpeed)))
	if i.render.camera.position.Y() <= 0.1 {
		i.render.camera.position = mgl32.Vec3{i.render.camera.position.X(), 0.1, i.render.camera.position.Z()}
	}
}

func (i *input) processMouseInput(initialMove *bool, lastX *float64, lastY *float64, xpos float64, ypos float64) {
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

	i.render.camera.yaw += float32(xoffset)
	i.render.camera.pitch -= float32(yoffset)
	if i.render.camera.pitch > 90 {
		i.render.camera.pitch = 90
	}
	if i.render.camera.pitch < -89 {
		i.render.camera.pitch = -89
	}
	v3 := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(i.render.camera.yaw))) * math.Cos(float64(mgl32.DegToRad(i.render.camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(i.render.camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(i.render.camera.yaw))) * math.Cos(float64(mgl32.DegToRad(i.render.camera.pitch)))),
	}
	i.render.camera.front = v3.Normalize()
}

func (i *input) processKeyboardInput(window *glfw.Window, deltaTime float32, cam *camera) {
	if window.GetKey(glfw.KeyR) == glfw.Press {
		if glfw.GetTime()-i.lastPress > 0.250 {
			if i.mode == gl.LINE {
				i.mode = gl.FILL
			} else {
				i.mode = gl.LINE
			}
			i.lastPress = glfw.GetTime()
			gl.PolygonMode(gl.FRONT_AND_BACK, i.mode)
		}
		return
	}

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	// Increase speed the farther away we are from the center.
	//distToCenter := cam.position.Sub(mgl32.Vec3{0, 0, 0}).Len()
	//centerFactor := distToCenter*distToCenter * 0.0

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
