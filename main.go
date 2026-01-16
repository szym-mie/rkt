package main

import (
	"log"
	"runtime"

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

func onKey(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, mods glfw.ModifierKey) {
	log.Printf("key: %s\n", glfw.GetKeyName(key, sc))
	if key == glfw.KeyQ {
		log.Printf("== quit ==\n")
		w.SetShouldClose(true)
	}
}

const w = 320 * 3
const h = 240 * 3

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

	gl.ClearColor(0.2, 0.7, 0.8, 0.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearDepth(1.0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	f := float64(w)/h - 1
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 10.0)

	ctrl1 := rkt.LoadPartDef("res/ctrl1.part.json")
	solid2 := rkt.LoadPartDef("res/solid2.part.json")
	decoupa := rkt.LoadPartDef("res/decoupa.part.json")

	v := rkt.NewVehicle("test", rkt.NewPartNode(ctrl1.New()))
	v.PartTree.
		AttachBelow(solid2.New()).
		AttachBelow(decoupa.New()).
		AttachBelow(solid2.New())

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Translatef(4.0, 0.0, -8.0)
		gl.Rotatef(30.0, 1.0, 0.0, 0.0)
		gl.Rotatef(90.0, 0.0, 1.0, 0.0)

		v.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
