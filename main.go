package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	_ "image/png"
	"log"
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

	// Send vertices to GPU memory for later consumption.
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

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
	gl.UseProgram(shaderProgram)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.39, 0.39, 0.39, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

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

void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
} 
` + "\x00"
