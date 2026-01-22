package main

import (
	"log"
	"runtime"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	rkt "szymm.org/rkt/src"
)

func init() {
	log.Printf("== init ==\n")
	runtime.LockOSThread()
}

var v *rkt.Vehicle
var radius float32

func onKey(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, mods glfw.ModifierKey) {
	log.Printf("key: %s\n", glfw.GetKeyName(key, sc))
	if key == glfw.KeyQ {
		log.Printf("== quit ==")
		w.SetShouldClose(true)
	}
	if key == glfw.KeyS && act == glfw.Press {
		log.Printf("== stage ==")
		v.ApplyStage()
	}
	if key == glfw.KeyBackspace && act == glfw.Press {
		radius = 10.0
	}
	if key == glfw.KeyEqual && act == glfw.Press {
		radius *= 0.7
	}
	if key == glfw.KeyMinus && act == glfw.Press {
		radius *= 1.2
	}
}

const scale = 3
const w = 320 * scale
const h = 240 * scale

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(w, h, "rkt", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetKeyCallback(onKey)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.2, 0.7, 0.8, 0.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearDepth(1.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	camera := rkt.NewCamera(1.0, 8000.0, 100.0)
	camera.SetViewport(w, h)
	camera.CaptureMouse(window)

	radius = 10.0

	rkt.LoadPkg("res/base.zip")

	v = rkt.NewVehicle("test", rkt.NewPart("base/ctrl1"))
	p := v.Parts
	p = v.AttachBelow(p, rkt.NewPart("base/solid2"))
	p = v.AttachBelow(p, rkt.NewPart("base/decoupa"))
	v.AddStage()
	p = v.AttachBelow(p, rkt.NewPart("base/solid2"))
	camera.Target = v

	patch := rkt.NewPatch("base/geom/patch00")
	patch.Scale = 1600.0
	patch.Pos.Z = -10.6

	labelPressQ := rkt.NewLabelFor("base/font/anlg", "press 'q' to quit")
	labelPressS := rkt.NewLabelFor("base/font/anlg", "press 's' to fire next stage")
	labelUseZoom := rkt.NewLabelFor("base/font/anlg", "press '+/-' to zoom in/out")
	labelUseMouse := rkt.NewLabelFor("base/font/anlg", "move mouse to look around")

	labelPressQ.Scale = 1.0
	labelPressS.Scale = 1.0
	labelUseZoom.Scale = 1.0
	labelUseMouse.Scale = 1.0

	labelPressQ.Pos.Y = 6.0
	labelPressQ.Pos.Z = -6.0
	labelPressS.Pos.Y = 5.0
	labelPressS.Pos.Z = -6.0
	labelUseZoom.Pos.Y = 4.0
	labelUseZoom.Pos.Z = -6.0
	labelUseMouse.Pos.Y = 3.0
	labelUseMouse.Pos.Z = -6.0

	n := v.Stages
	for n != nil {
		name := "<sep>"
		if n.Part != nil {
			name = n.Part.GetName()
		}
		log.Printf("%v, ", name)
		n = n.Next
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera.Radius = radius

		x, y := window.GetCursorPos()
		mousePos := rkt.Vec2{X: float32(x), Y: float32(y)}
		camera.Update(mousePos)
		camera.Apply()

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		log.Printf("pz %v vz %v", v.Pos.Z, v.Vel.Z)

		dt := time.Millisecond * 50

		patch.Draw()
		v.Draw()

		labelPressQ.Draw()
		labelPressS.Draw()
		labelUseZoom.Draw()
		labelUseMouse.Draw()

		v.Update(float32(dt.Seconds()))

		time.Sleep(dt)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
