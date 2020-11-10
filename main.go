package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

const (
	width  = 1000
	height = 1000
	fps    = 60
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program, _, _ := initOpenGL()

	// Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)
	// Camera is at (4,3,3), in World Space
	// and looks at the origin
	// Head is up (set to 0,-1,0 to look upside-down)
	view := mgl32.LookAt(
		4, 3, 3,
		0, 0, 0,
		0, 1, 0,
	)
	model := mgl32.Ident4()
	//sm := mgl32.Scale3D(0.5, 0.5, 0.5)
	//ro := mgl32.HomogRotate3D(mgl32.DegToRad(45), [3]float32{0, -1, 0})
	//model = model.Mul4(sm).Mul4(ro)
	fmt.Printf("%v\n", model)

	mvp := projection.Mul4(view.Mul4(model))
	matrix := gl.GetUniformLocation(program, gl.Str("MVP\x00"))

	// Translation of C Code glfwSetInputMode(window, GLFW_STICKY_KEYS, GL_TRUE)
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	points := []float32{
		-1.0, -1.0, -1.0, // triangle 1 : begin
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0, -1.0, // triangle 2 : begin
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0, // triangle 2 : end
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	}

	var vertextbuffer uint32
	gl.GenBuffers(1, &vertextbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertextbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.ClearColor(0, 0, 0.4, 1)
	for !window.ShouldClose() {
		// Input handling
		key := window.GetKey(glfw.KeyEscape)
		if key == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.UniformMatrix4fv(matrix, 1, false, &mvp[0])

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertextbuffer)
		// 0 in index refers to location in layout for vec3
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3)
		gl.DisableVertexAttribArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	//gl.DetachShader(program, vid)
	//gl.DetachShader(program, fid)
}

func draw(window *glfw.Window, program uint32) {
	glfw.PollEvents()
	window.SwapBuffers()
}
