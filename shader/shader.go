package shader

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"io/ioutil"
	"strings"
)

const (
	Vertex   = gl.VERTEX_SHADER
	Fragment = gl.FRAGMENT_SHADER
)

func Compile(shaderType uint32, filename string) (uint32, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("error reading shader file %s: %v", filename, err)
	}

	var shaderLoc uint32
	shaderLoc = gl.CreateShader(shaderType)
	csources, free := gl.Strs(string(bs) + "\x00")
	gl.ShaderSource(shaderLoc, 1, csources, nil)
	free()
	gl.CompileShader(shaderLoc)

	var status int32
	gl.GetShaderiv(shaderLoc, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderLoc, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderLoc, logLength, nil, gl.Str(logs))
		return 0, fmt.Errorf("failed to compile %s: %v", filename, logs)
	}

	return shaderLoc, nil
}
