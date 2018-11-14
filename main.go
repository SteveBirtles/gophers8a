package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 1024
	height = 768
)

var (
	vao         uint32
	cameraX	    float32 = -0.5
	cameraY		float32 = 0.0
	zoom        float32 = 1.0
	power  		float32 = 1.0
	infinity  	float32 = 1e32
	iterations  int32 = 100

	frames            = 0
	frameLength       float32
	second            = time.Tick(time.Second)
	windowTitlePrefix = "Mandelbrot with Shaders"
)

func init() {
	runtime.LockOSThread()
}


func compileShader(filename string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to read file %v", filename)
	}

	shaderSource, free := gl.Strs(string(source) + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, shaderSource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil

}


func main() {

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, windowTitlePrefix, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(0)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	vertexShader, err := compileShader("shader.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader("shader.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	iterationsVar := gl.GetUniformLocation(program, gl.Str("iterations\x00"))
	cameraVar := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	zoomVar := gl.GetUniformLocation(program, gl.Str("zoom\x00"))
	powerVar := gl.GetUniformLocation(program, gl.Str("power\x00"))
	infinityVar := gl.GetUniformLocation(program, gl.Str("infinity\x00"))

	gl.GenVertexArrays(1, &vao)

	gl.ClearColor(0,0,0,1)

	for !window.ShouldClose() {

		frameStart := time.Now()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			cameraX -= frameLength*zoom
		}
		if window.GetKey(glfw.KeyRight) == glfw.Press {
			cameraX += frameLength*zoom
		}
		if window.GetKey(glfw.KeyUp) == glfw.Press {
			cameraY += frameLength*zoom
		}
		if window.GetKey(glfw.KeyDown) == glfw.Press {
			cameraY -= frameLength*zoom
		}
		if window.GetKey(glfw.KeySpace) == glfw.Press {
			zoom /= 1+frameLength
		}
		if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
			zoom *= 1+frameLength
		}

		if window.GetKey(glfw.KeyA) == glfw.Press {
			power /= 1+frameLength
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			power *= 1+frameLength
		}

		if window.GetKey(glfw.KeyW) == glfw.Press {
			iterations += 1
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			iterations -= 1
			if iterations < 1 { iterations = 1 }
		}

		gl.UseProgram(program)
		gl.Uniform1i(iterationsVar, iterations)
		gl.Uniform2f(cameraVar, cameraX, cameraY)
		gl.Uniform1f(zoomVar, zoom)
		gl.Uniform1f(powerVar, power)
		gl.Uniform1f(infinityVar, infinity)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)

		gl.UseProgram(0)

		glfw.PollEvents()
		window.SwapBuffers()

		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}
		frameLength = float32(time.Since(frameStart).Seconds())

	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	glfw.Terminate()

}
