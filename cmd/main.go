package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/DonMatano/learnOpenGLGo/lib"
	openglshader "github.com/DonMatano/learnOpenGLGo/openGlShader"
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

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Learn Open GL", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("Failed to start gl", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("Running opengl Version %s", version)
	gl.Viewport(0, 0, 800, 600)
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	// vertices
	// triangleVertices := []float32{
	// 	// x, y, z   Colours
	// 	0.5, -0.5, 0, 1, 0, 0, // bottom right
	// 	-0.5, -0.5, 0, 0, 1, 0, // bottom left
	// 	0, 0.5, 0, 0, 0, 1, // bottom left
	// }

	// indices := []uint32{
	// 	0, 1, 3,
	// 	1, 2, 3,
	// }
	//
	// Rectangle
	rectangleVertices := []float32{
		// positions    //colours     // texture
		0.5, 0.5, 0, 1, 0, 0, 1, 1, // top right
		0.5, -0.5, 0, 0, 1, 0, 1, 0, // bottom right
		-0.5, -0.5, 0, 0, 0, 1, 0, 0, // bottom left
		-0.5, 0.5, 0, 1, 1, 0, 0, 1, // bottom left
	}

	rectangleIndices := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	shader, err := openglshader.NewShader("shaders/vertexShader.glsl", "shaders/fragmentShader.glsl")
	if err != nil {
		log.Printf("Error getting Shader: \n %v", err)
	}

	// configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(rectangleVertices)*4, gl.Ptr(rectangleVertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(rectangleIndices)*4, gl.Ptr(rectangleIndices), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)
	// color attribute
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)
	// texture attribute
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	// load and create a texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// Set texture wrapping
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// Set texture filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// load image, create texture and generate mipmaps
	data, err := lib.LoadImage("resources/textures/container.jpg")
	if err != nil {
		log.Fatalln("Failed to load texture", err)
	}
	log.Println("data received", data.Rect.Size())
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(data.Rect.Size().X), int32(data.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	log.Println("Generated MipMap")
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindVertexArray(0)
	// Wireframe
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	for !window.ShouldClose() {
		// input
		processInput(window)
		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// bind Texture
		gl.BindTexture(gl.TEXTURE_2D, texture)

		shader.Use()

		// update uniform ourColour
		// timeValue := glfw.GetTime()
		// greenValue := math.Sin(timeValue)
		// vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColour"+shaderEndChar))
		// gl.Uniform4f(vertexColorLocation, 0, float32(greenValue), 0, 1)
		gl.BindVertexArray(vao)
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(nil))

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
	shader.Delete()
}
