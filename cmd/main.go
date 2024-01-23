package main

import (
	"fmt"
	"log"
	"math"
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
		fmt.Println("Closing window...")
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
		0.5, 0.5, 0, // top right
		0.5, -0.5, 0, // bottom right
		-0.5, -0.5, 0, // bottom left
		-0.5, 0.5, 0, // top left
	}

	indices := []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	// configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*4, gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// vertexShader

	vertexShaderSource := `
  #version 330 core

  layout (location = 0) in vec3 aPos;
  out vec4 vertexColor;

  void main() {
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0); 
  }
` + shaderEndChar
	compileShader(vertexShaderSource, gl.VERTEX_SHADER)

	// fragmentShader

	fragmentShaderSource := `
  #version 330 core
  out vec4 FragColor;
  uniform vec4 ourColour;
  void main() {
    FragColor = ourColour;
  }
` + shaderEndChar

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}
	// Wireframe
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	for !window.ShouldClose() {
		// input
		processInput(window)
		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)
		// update uniform ourColour
		timeValue := glfw.GetTime()
		greenValue := math.Sin(timeValue)
		vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColour"+shaderEndChar))
		gl.Uniform4f(vertexColorLocation, 0, float32(greenValue), 0, 1)
		gl.BindVertexArray(vao)
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(nil))

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
	gl.DeleteProgram(program)
}
