package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	models "github.com/mlesniak/opengl/model"
	"github.com/mlesniak/opengl/shader"
	_ "image/png"
	"log"
)

// TODO(mlesniak) add plane
// TODO(mlesniak) add diffuse lightning

// TODO(mlesniak) show plane first
// TODO(mlesniak) struct
// TODO(mlesniak) General refactoring for render loop based on objects

var camPos = mgl32.Vec3{0, 1, 3}
var camFront = mgl32.Vec3{0, 0, -1}
var camUp = mgl32.Vec3{0, 1, 0}

func render(window *glfw.Window) {
	// General openGL configuration.
	gl.Enable(gl.DEPTH_TEST)

	// VAOs  store vertex attribute configuration and which VBOs to use.
	//var cube uint32
	//gl.GenVertexArrays(1, &cube)
	//gl.BindVertexArray(cube)
	//log.Print("Created plane")

	// Send cubeData to GPU memory for later consumption.
	// Copy our cubeData array (including other data such as colors, texture information)
	// in a memory buffer for OpenGL to use.
	//var cubeData uint32
	//gl.GenBuffers(1, &cubeData)
	//gl.BindBuffer(gl.ARRAY_BUFFER, cubeData)
	//gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.Cube), gl.Ptr(models.Cube), gl.STATIC_DRAW)
	//log.Print("Created cubeData and stored vertices")
	//
	//// Set the vertex attributes pointers, i.e. configure where vertex pointers are located
	//// to be used in location 0 in the vertex shader.
	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*SizeFloat32), nil)
	//gl.EnableVertexAttribArray(0)
	//log.Print("Bound vertices to location")

	var plane uint32
	gl.GenVertexArrays(1, &plane)
	gl.BindVertexArray(plane)
	log.Print("Created plane")

	var planeData uint32
	gl.GenBuffers(1, &planeData)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeData)
	gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.Plane), gl.Ptr(&models.Plane[3]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*SizeFloat32), nil)
	gl.EnableVertexAttribArray(0)

	var program uint32
	program = gl.CreateProgram()

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
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	log.Print("Created program")

	gl.UseProgram(program)

	model := mgl32.Ident4()
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

	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform3fv(colorUniform, 1, &models.Plane[0])
	//gl.Uniform3f(colorUniform, 1.0, 1, 0)

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

	window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		camSpeed := -yoff * 0.2
		camPos = camPos.Add(camUp.Mul(float32(camSpeed)))
	})

	var deltaTime float32 = 0
	var lastFrame float64 = 0
	log.Print("Starting rendering loop")
	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame - lastFrame)
		lastFrame = currentFrame

		processKeyboard(window, deltaTime)

		gl.ClearColor(0.19, 0.19, 0.19, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our predefined vertex and fragment shaders to render
		// vertext data wrapped into the plane.
		gl.UseProgram(program)

		// Update all matrices.
		view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		//gl.BindVertexArray(cube)
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.Cube)/3))

		gl.BindVertexArray(plane)
		model = mgl32.HomogRotate3DX(mgl32.DegToRad(-90))
		model = model.Mul4(mgl32.Scale3D(20, 20, 1))
		model = model.Mul4(mgl32.Translate3D(-0.5, -0.5, 0.0))
		//model = model.Mul4(mgl32.Scale3D(10, 10, 10))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.Plane)/3))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
