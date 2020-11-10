package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"io/ioutil"
	"strings"
)

func initOpenGL() (uint32, uint32, uint32) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("VERSION %s\n", version)

	bs, err := ioutil.ReadFile("simple.fragmentshader")
	if err != nil {
		panic(err)
	}
	fid, err := compileShader(bs, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("simple.fragmentshader")
		panic(err)
	}

	bs, err = ioutil.ReadFile("simple.vertexshader")
	if err != nil {
		panic(err)
	}
	vid, err := compileShader(bs, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("simple.vertexshader")
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, fid)
	gl.AttachShader(program, vid)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		fmt.Printf("failed to link program: %v\n", log)
		panic("bye")
	}

	gl.DeleteShader(fid)
	gl.DeleteShader(vid)

	return program, fid, vid
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic("unable to initialize glfw")
	}

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}
