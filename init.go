package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"io/ioutil"
)

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadFile("simple.fragmentshader")
	if err != nil {
		panic(err)
	}
	fid, _ := compileShader(bs, gl.FRAGMENT_SHADER)

	bs, err = ioutil.ReadFile("simple.vertexshader")
	if err != nil {
		panic(err)
	}
	vid, _ := compileShader(bs, gl.VERTEX_SHADER)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, fid)
	gl.AttachShader(prog, vid)
	gl.LinkProgram(prog)
	return prog
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic("unable to initialize glfw")
	}

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}
