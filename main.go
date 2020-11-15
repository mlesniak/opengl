package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/mlesniak/opengl/shader"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"
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
	// Coordinates
	vertices := []float32{
		// positions          // colors           // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom let
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top let
	}

	// use colors here, too.

	// Use an EBO to reference above coordinates.
	indices := []int32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
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

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// 1. then set the vertex attributes pointers
	// position attributes.
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)
	// For colors
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// For texture
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	texture, _ := newTexture("square.png")
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Compile shaders.
	vertexShader, err := shader.Compile(shader.Vertex, "vertex.shader")
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.Compile(shader.Fragment, "fragment.shader")
	if err != nil {
		panic(err)
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
		gl.UseProgram(shaderProgram)

		// Draw triangles.
		// To draw our objects of choice, OpenGL provides us with the
		// glDrawArrays function that draws primitives using the currently
		// active shader, the previously defined vertex attribute
		// configuration and with the VBO's vertex data (indirectly bound
		// via the VAO).
		//
		//gl.BindVertexArray(vao)
		//gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// Use ebo / index
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

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

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
