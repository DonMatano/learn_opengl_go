package openglshader

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type shader struct {
	progId uint32
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	sourceInC, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, sourceInC, nil)
	free()
	gl.CompileShader(shader)
	var hasSuccessfullyCompiled int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &hasSuccessfullyCompiled)
	if hasSuccessfullyCompiled == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var hasSuccessfullyCompiled int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &hasSuccessfullyCompiled)
	if hasSuccessfullyCompiled == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func NewShader(vertexFilePath, fragmentFilePath string) (*shader, error) {
	vertexCode, err := os.ReadFile(vertexFilePath)
	if err != nil {
		return nil, err
	}

	fragmentCode, err := os.ReadFile(fragmentFilePath)
	if err != nil {
		return nil, err
	}
	progId, err := newProgram(string(vertexCode), string(fragmentCode))
	if err != nil {
		return nil, err
	}

	return &shader{progId: progId}, nil
}

func (sh *shader) Use() {
	gl.UseProgram(sh.progId)
}

func (sh *shader) Delete() {
	gl.DeleteProgram(sh.progId)
}
