package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var framebufferSizeCallback glfw.FramebufferSizeCallback = func(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
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

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "Learn Open GL", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("Failed to start gl", err)
	}
	gl.Viewport(0, 0, 800, 600)
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	// build and compile our shader program
	shaderEndChar := "\x00"
	// vertices
	triangleVertices := []float32{
		// x, y, z
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
		0, 0.5, 0, // top

	}

	// configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*4, gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)

	// vertexShader

	vertexShaderSource := `
  #version 330 core

  layout (location = 0) in vec3 aPos;

  void main() {
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0); 
  }
` + shaderEndChar
	compileShader(vertexShaderSource, gl.VERTEX_SHADER)

	// fragmentShader

	fragmentShaderSource := `
  #version 330 core
  out vec4 FragColor;
  void main() {
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
  }
` + shaderEndChar

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	for !window.ShouldClose() {
		// input
		processInput(window)
		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteProgram(program)
}
