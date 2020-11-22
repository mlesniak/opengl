package main

// TODO(mlesniak) Move render to own module

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/shader"
	_ "image/png"
	"log"
)

var camPos = mgl32.Vec3{0, 3, 3}
var camFront = mgl32.Vec3{0, 0, -1}
var camUp = mgl32.Vec3{0, 1, 0}

var lastX = float64(windowWidth / 2)
var lastY = float64(windowHeight / 2)

var yaw = float32(-90.0)
var pitch = float32(0.0)

func renderLoop(window *glfw.Window, scene *Scene) {
	vaos := make([]uint32, len(scene.Entities))
	for i, e := range scene.Entities {
		var vaoIndex uint32
		gl.GenVertexArrays(1, &vaoIndex)
		gl.BindVertexArray(vaoIndex)

		var vaoData uint32
		gl.GenBuffers(1, &vaoData)
		gl.BindBuffer(gl.ARRAY_BUFFER, vaoData)
		gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(e.Vertices), gl.Ptr(e.Vertices), gl.STATIC_DRAW)

		var stride = 3
		if e.WithNormal {
			stride = 6
		}
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(stride*SizeFloat32), nil)
		gl.EnableVertexAttribArray(0)
		if e.WithNormal {
			gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(stride*SizeFloat32), gl.PtrOffset(3*SizeFloat32))
			gl.EnableVertexAttribArray(1)
		} else {
			// TODO(mlesniak) Fixed normal value?
		}

		vaos[i] = vaoIndex
	}

	program := createProgram()

	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	projection := mgl32.Perspective(mgl32.DegToRad(60), windowWidth/float32(windowHeight), 0.1, 1000)
	projection = projection.Mul4(mgl32.Translate3D(0, 0, -3))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	log.Print("Matrices initialized")

	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	lightPos := gl.GetUniformLocation(program, gl.Str("lightPos\x00"))

	var deltaTime float32 = 0
	var lastFrameTime float64 = 0

	log.Print("Starting rendering loop")
	for !window.ShouldClose() {
		currentFrameTime := glfw.GetTime()
		deltaTime = float32(currentFrameTime - lastFrameTime)
		lastFrameTime = currentFrameTime

		processKeyboardInput(window, deltaTime)

		gl.ClearColor(0.39, 0.39, 0.39, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		gl.Uniform3f(lightPos, camPos.X(), camPos.Y()+3, camPos.Z())

		// Update all matrices.
		view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

		for i, v := range vaos {
			e := scene.Entities[i]
			gl.BindVertexArray(v)
			gl.Uniform4f(colorUniform, e.Color.X(), e.Color.Y(), e.Color.Z(), 1.0)
			gl.UniformMatrix4fv(modelUniform, 1, false, &e.Position[0])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(e.Vertices)/3))
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createProgram() uint32 {
	var program uint32
	program = gl.CreateProgram()

	vertexShader, err := shader.Compile(shader.Vertex, "data/vertex.shader")
	if err != nil {
		log.Fatalf("error compiling vertex shader: %v", err)
	}
	fragmentShader, err := shader.Compile(shader.Fragment, "data/fragment.shader")
	if err != nil {
		log.Fatalf("error compiling fragment shader: %v", err)
	}

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	gl.UseProgram(program)

	return program
}
