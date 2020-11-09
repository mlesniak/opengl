package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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
	program := initOpenGL()

	// Translation of C Code glfwSetInputMode(window, GLFW_STICKY_KEYS, GL_TRUE)
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	points := []float32{
		-1.0, -1.0, 0,
		1, -1, 0,
		0, 1, 0,
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
			break
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertextbuffer)
		// 0 in index refers to location in layout for vec3
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DisableVertexAttribArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func draw(window *glfw.Window, program uint32) {
	glfw.PollEvents()
	window.SwapBuffers()
}
