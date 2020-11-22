package main

// TODO(mlesniak) Move render to own module

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	models "github.com/mlesniak/opengl/model"
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

type Entity struct {
	Vertices   []float32
	Position   mgl32.Mat4
	WithNormal bool
	Color      mgl32.Vec4
}

var geometries = make([]*Entity, 0)

func AddGeometry(object *Entity) {
	// TODO(mlesniak) Generate object here

	geometries = append(geometries, object)
}

func Cube() *Entity {
	//var cube uint32
	//gl.GenVertexArrays(1, &cube)
	//gl.BindVertexArray(cube)
	//log.Print("Created plane")
	//
	//// Send cubeData to GPU memory for later consumption.
	//// Copy our cubeData array (including other data such as colors, texture information)
	//// in a memory buffer for OpenGL to use.
	//var cubeData uint32
	//gl.GenBuffers(1, &cubeData)
	//gl.BindBuffer(gl.ARRAY_BUFFER, cubeData)
	//gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.CubeVertices), gl.Ptr(&models.CubeVertices[3]), gl.STATIC_DRAW)
	//log.Print("Created cubeData and stored vertices")
	//
	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(6*SizeFloat32), nil)
	//gl.EnableVertexAttribArray(0)
	//gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(6*SizeFloat32), gl.PtrOffset(3*SizeFloat32))
	//gl.EnableVertexAttribArray(1)
	//log.Print("Bound vertices to location")

	return &Entity{
		Vertices:   models.CubeVertices,
		Position:   mgl32.Ident4(),
		WithNormal: true,
		Color:      mgl32.Vec4{1, 0, 0, 1},
	}

}

func renderLoop(window *glfw.Window, scene Scene) {
	//AddGeometry(Cube())

	// Store vaos for geometries.
	vaos := make([]uint32, len(geometries))
	for i, e := range geometries {
		var vaoIndex uint32
		gl.GenBuffers(1, &vaoIndex)
		gl.BindBuffer(gl.ARRAY_BUFFER, vaoIndex)
		gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(e.Vertices), gl.Ptr(&e.Vertices[0]), gl.STATIC_DRAW)

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

	//cube := makeCube()
	plane := makePlane()
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

		gl.Uniform3f(lightPos, camPos.X(), camPos.Y()+3, camPos.Z())

		// Update all matrices.
		view := mgl32.LookAtV(camPos, camPos.Add(camFront), camUp)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		for i, v := range vaos {
			g := geometries[i]
			gl.BindVertexArray(v)
			gl.Uniform4f(colorUniform, g.Color.X(), g.Color.Y(), g.Color.Z(), g.Color.W())
			gl.UniformMatrix4fv(modelUniform, 1, false, &geometries[i].Position[0])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(geometries[i].Vertices)/3))
		}

		gl.BindVertexArray(plane)
		model = mgl32.HomogRotate3DX(mgl32.DegToRad(-90))
		model = model.Mul4(mgl32.Scale3D(20, 20, 1))
		model = model.Mul4(mgl32.Translate3D(-0.5, -0.5, 0.0))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.Uniform4f(colorUniform, 0.2, 0.2, 0.2, 1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.Plane)/3))

		// Use scene...
		// objects
		//
		// list of vertices, normals
		// color is ignored for now.
		//
		// primitives are loaded using vaos.
		//		must vaos be initialized for each TYPE before?
		// object describes type, position, scale

		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(0, 0.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(0, 1.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(0, 2.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(-1, 2.5, 0))
		//model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(-2, 2.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(-2, 0.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(-2, 1.5, 0))
		//model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))
		//
		//gl.BindVertexArray(cube)
		//gl.Uniform3fv(colorUniform, 1, &models.CubeVertices[0])
		//model = mgl32.Ident4()
		//model = model.Mul4(mgl32.Translate3D(-2, 2.5, 0))
		//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(models.CubeVertices)/3))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func makePlane() uint32 {
	var plane uint32
	gl.GenVertexArrays(1, &plane)
	gl.BindVertexArray(plane)
	log.Print("Created plane")

	var planeData uint32
	gl.GenBuffers(1, &planeData)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeData)
	gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.Plane), gl.Ptr(&models.Plane[3]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(6*SizeFloat32), nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(6*SizeFloat32), gl.PtrOffset(3*SizeFloat32))
	gl.EnableVertexAttribArray(1)
	return plane
}

func makeCube() uint32 {
	var cube uint32
	gl.GenVertexArrays(1, &cube)
	gl.BindVertexArray(cube)
	log.Print("Created plane")

	// Send cubeData to GPU memory for later consumption.
	// Copy our cubeData array (including other data such as colors, texture information)
	// in a memory buffer for OpenGL to use.
	var cubeData uint32
	gl.GenBuffers(1, &cubeData)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeData)
	gl.BufferData(gl.ARRAY_BUFFER, SizeFloat32*len(models.CubeVertices), gl.Ptr(&models.CubeVertices[3]), gl.STATIC_DRAW)
	log.Print("Created cubeData and stored vertices")

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(6*SizeFloat32), nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(6*SizeFloat32), gl.PtrOffset(3*SizeFloat32))
	gl.EnableVertexAttribArray(1)
	log.Print("Bound vertices to location")
	return cube
}

func createProgram() uint32 {
	var program uint32
	program = gl.CreateProgram()

	// Compile shaders.
	vertexShader, err := shader.Compile(shader.Vertex, "data/vertex.shader")
	if err != nil {
		log.Fatalf("error compiling vertex shader: %v", err)
	}
	fragmentShader, err := shader.Compile(shader.Fragment, "data/fragment.shader")
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
	return program
}
