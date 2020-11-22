package render

// TODO(mlesniak) Move Render to own module

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mlesniak/opengl/scene"
	"github.com/mlesniak/opengl/shader"
	_ "image/png"
	"log"
	"runtime"
)

type render struct {
	width  int
	height int
	camera *camera

	window *glfw.Window
}

type camera struct {
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3

	yaw   float32
	pitch float32
}

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func New(width, height int) *render {
	window := initializeGraphics(width, height)
	camera := initializeCamera()
	r := &render{
		width:  width,
		height: height,
		camera: camera,
		window: window,
	}

	r.setupInput()
	return r
}

func initializeCamera() *camera {
	return &camera{
		position: mgl32.Vec3{0, 3, 3},
		front:    mgl32.Vec3{0, 0, -1},
		up:       mgl32.Vec3{0, 1, 0},
		yaw:      float32(-90.0),
		pitch:    float32(0.0),
	}
}

func (r *render) Exit() {
	glfw.Terminate()
}

func (r *render) Render(scene *scene.Scene) {
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

	view := mgl32.LookAtV(r.camera.position, r.camera.position.Add(r.camera.front), r.camera.up)
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	projection := mgl32.Perspective(mgl32.DegToRad(60), float32(r.width)/float32(r.height), 0.1, 1000)
	projection = projection.Mul4(mgl32.Translate3D(0, 0, -3))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	lightPos := gl.GetUniformLocation(program, gl.Str("lightPos\x00"))

	var deltaTime float32 = 0
	var lastFrameTime float64 = 0

	for !r.window.ShouldClose() {
		currentFrameTime := glfw.GetTime()
		deltaTime = float32(currentFrameTime - lastFrameTime)
		lastFrameTime = currentFrameTime

		r.processKeyboardInput(r.window, deltaTime, r.camera)

		gl.ClearColor(0.39, 0.39, 0.39, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		gl.Uniform3f(lightPos, r.camera.position.X(), r.camera.position.Y()+3, r.camera.position.Z())

		// Update all matrices.
		view := mgl32.LookAtV(r.camera.position, r.camera.position.Add(r.camera.front), r.camera.up)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

		for i, v := range vaos {
			e := scene.Entities[i]
			gl.BindVertexArray(v)
			gl.Uniform4f(colorUniform, e.Color.X(), e.Color.Y(), e.Color.Z(), 1.0)
			gl.UniformMatrix4fv(modelUniform, 1, false, &e.Position[0])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(e.Vertices)/3))
		}

		r.window.SwapBuffers()
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
