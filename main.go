package main

import (
	"log"
	"runtime"

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

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(320, 240, "rkt", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetKeyCallback(onKey)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0.2, 0.0, 0.0, 0.0)

	ctrl1 := rkt.LoadPartDef("res/ctrl1.part.json")
	// solid2 := rkt.LoadPartDef("res/solid2.part.json")
	// decoupa := rkt.LoadPartDef("res/decoupa.part.json")

	v := rkt.NewVehicle("test", rkt.NewPartNode(ctrl1.New()))
	// v.PartTree.AttachBelow(solid2.New())

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		v.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
