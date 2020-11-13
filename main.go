package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	_ "image/png"
	"log"
	"math"
	"runtime"
	"strings"
	//"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

func main() {
	window := initialize()
	render(window)
	glfw.Terminate()
}

func render(window *glfw.Window) {
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	// ... a VAO that stores our vertex attribute configuration and which VBO to use.
	// Usually when you have multiple objects you want to draw, you first generate/configure
	// all the VAOs (and thus the required VBO and attribute pointers) and store those for
	// later use. The moment we want to draw one of our objects, we take the corresponding VAO,
	// bind it, then draw the object and unbind the VAO again.
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Send vertices to GPU memory for later consumption.
	// 0. copy our vertices array in a buffer for OpenGL to use
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	// 1. then set the vertex attributes pointers
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	// Compile shaders.
	var vertexShader uint32
	vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, csources, nil)
	free()
	gl.CompileShader(vertexShader)

	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(logs))
		log.Printf("failed to compile %v: %v", vertexShaderSource, logs)
	}

	var fragmentShader uint32
	fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, csources, nil)
	free()
	gl.CompileShader(fragmentShader)

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(logs))
		log.Printf("failed to compile %v: %v", fragmentShaderSource, logs)
	}

	// Create program combining shaders.
	var shaderProgram uint32
	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	// Every shader and rendering call after glUseProgram will now use this program object (and thus the shaders).

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.39, 0.39, 0.39, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// 2. use our shader program when we want to render an object
		//
		// To draw our objects of choice, OpenGL provides us with the
		// glDrawArrays function that draws primitives using the currently
		// active shader, the previously defined vertex attribute
		// configuration and with the VBO's vertex data (indirectly bound
		// via the VAO).
		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vao)

		// Set uniform variable in shader.
		t := glfw.GetTime()
		green := (math.Sin(t) / 2) + 0.5
		vertexColorLoc := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor\x00"))
		gl.Uniform4f(vertexColorLoc, 0, float32(green), 0, 1)

		// Draw triangles.
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func initialize() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
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

var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"

var fragmentShaderSource = `
#version 330 core
out vec4 FragColor;

uniform vec4 ourColor;

void main()
{
    FragColor = ourColor;
} 
` + "\x00"
