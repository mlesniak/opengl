package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math/rand"
	"runtime"
	"time"
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
		4, 3, -3,
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

	colors := []float32{
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	}

	randomColors(colors)

	var vertextbuffer uint32
	gl.GenBuffers(1, &vertextbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertextbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var colorbuffer uint32
	gl.GenBuffers(1, &colorbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)

	gl.ClearColor(0, 0, 0.4, 1)
	for !window.ShouldClose() {
		// Input handling
		key := window.GetKey(glfw.KeyEscape)
		if key == glfw.Press {
			window.SetShouldClose(true)
		}

		// Random color on each frame.
		randomColors(colors)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.UniformMatrix4fv(matrix, 1, false, &mvp[0])

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertextbuffer)
		// 0 in index refers to location in layout for vec3
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3)
		gl.DisableVertexAttribArray(0)

		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
		gl.VertexAttribPointer(
			1, // layout(location = 1)
			3,
			gl.FLOAT,
			false,
			0,
			nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	//gl.DetachShader(program, vid)
	//gl.DetachShader(program, fid)
}

func randomColors(colors []float32) {
	rand.Seed(time.Now().UnixNano())
	num := 1
	for i := 0; i < num; i++ {
		j := rand.Intn(len(colors))
		colors[j] = rand.Float32()
	}
}

func draw(window *glfw.Window, program uint32) {
	glfw.PollEvents()
	window.SwapBuffers()
}
