package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	models "github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/shader"
	_ "image/png"
	"log"
	"math"
	"runtime"
)

// TODO(mlesniak) add plane
// TODO(mlesniak) add diffuse lightning

const windowWidth = 800
const windowHeight = 600

// TODO(mlesniak) struct
var camPos = mgl32.Vec3{0, 1, 3}
var camFront = mgl32.Vec3{0, 0, -1}
var camUp = mgl32.Vec3{0, 1, 0}

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func initializeWindow() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initializeWindow glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return window
}

func main() {
	window := initializeWindow()
	render(window)
	glfw.Terminate()
}

func render(window *glfw.Window) {
	// General openGL configuration.
	gl.Enable(gl.DEPTH_TEST)

	// VAOs  store vertex attribute configuration and which VBOs to use.
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	log.Print("Created vao")

	// Send cube to GPU memory for later consumption.
	// Copy our cube array (including other data such as colors, texture information)
	// in a memory buffer for OpenGL to use.
	var cube uint32
	gl.GenBuffers(1, &cube)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube)
	gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.Cube), gl.Ptr(models.Cube), gl.STATIC_DRAW)
	log.Print("Created cube and stored vertices")

	// Set the vertex attributes pointers, i.e. configure where vertex pointers are located
	// to be used in location 0 in the vertex shader.
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*SizeFloat32), nil)
	gl.EnableVertexAttribArray(0)
	log.Print("Bound vertices to location")

	// Compile shaders.
	vertexShader, err := shader.Compile(shader.Vertex, "vertex.shader")
	if err != nil {
		log.Fatalf("error compiling vertex shader: %v", err)
	}
	fragmentShader, err := shader.Compile(shader.Fragment, "fragment.shader")
	if err != nil {
		log.Fatalf("error compiling fragment shader: %v", err)
	}
	log.Print("Compiled shader")

	// Create a program to combine shaders.
	var program uint32
	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	log.Print("Created program")

	// Create matrices for view.
	model := mgl32.Ident4()
	model = mgl32.HomogRotate3D(mgl32.DegToRad(-55), mgl32.Vec3{1, 0, 0})
	// Inject matrix into every shader based on the name.
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	projection := mgl32.Perspective(mgl32.DegToRad(60), windowWidth/float32(windowHeight), 0.1, 100)
	projection = projection.Mul4(mgl32.Translate3D(0, 0, -3))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	log.Print("Matrices initialized")

	// Configuration for movement.
	yaw := float32(-90.0)
	pitch := float32(0.0)
	initialMove := true

	lastX := float64(windowWidth / 2)
	lastY := float64(windowHeight / 2)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		processMouse(&initialMove, lastX, xpos, lastY, ypos, yaw, pitch)
	})

	var deltaTime float32 = 0
	var lastFrame float64 = 0
	log.Print("Starting rendering loop")
	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame - lastFrame)
		lastFrame = currentFrame

		processKeyboard(window, deltaTime)

		gl.ClearColor(0.39, 0.39, 0.39, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our predefined vertex and fragment shaders to render
		// vertext data wrapped into the vao.
		gl.UseProgram(program)

		// Update all matrices.
		view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processMouse(initialMove *bool, lastX float64, xpos float64, lastY float64, ypos float64, yaw float32, pitch float32) {
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
	if pitch > 89 {
		pitch = 89
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

func processKeyboard(window *glfw.Window, deltaTime float32) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	var camSpeed = 5 * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camPos = camPos.Add(camFront.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camPos = camPos.Sub(camFront.Mul(camSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		x := camFront.Cross(camUp).Normalize().Mul(camSpeed)
		camPos = camPos.Sub(x)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		x := camFront.Cross(camUp).Normalize().Mul(camSpeed)
		camPos = camPos.Add(x)
	}
}
