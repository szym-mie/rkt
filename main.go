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

type GameMode uint

const (
	EditMode GameMode = iota + 1
	PlayMode
)

var mainVehicle *rkt.Vehicle
var radius float32

const ctrlSpeed = 0.4

func ctrlVector(axis rkt.VecAxis, scale float32) rkt.Vec3 {
	comp := 1.0 * scale
	switch axis {
	case rkt.XAxis:
		return rkt.Vec3{X: comp, Y: 0.0, Z: 0.0}
	case rkt.YAxis:
		return rkt.Vec3{X: 0.0, Y: comp, Z: 0.0}
	case rkt.ZAxis:
		return rkt.Vec3{X: 0.0, Y: 0.0, Z: comp}
	}

	return rkt.Vec3{}
}

func onKey(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, mods glfw.ModifierKey) {
	if act == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		case glfw.KeySpace:
			mainVehicle.ApplyStage()
		case glfw.KeyBackspace:
			radius = 10.0
		case glfw.KeyEqual:
			radius *= 0.7
		case glfw.KeyMinus:
			radius *= 1.2
		case glfw.KeyA:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.XAxis, +ctrlSpeed)))
		case glfw.KeyD:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.XAxis, -ctrlSpeed)))
		case glfw.KeyW:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.YAxis, -ctrlSpeed)))
		case glfw.KeyS:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.YAxis, +ctrlSpeed)))
		case glfw.KeyQ:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.ZAxis, -ctrlSpeed)))
		case glfw.KeyE:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.ZAxis, +ctrlSpeed)))
		}
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

	rkt.InitTextureUnit(0)
	rkt.InitTextureUnit(1)
	rkt.InitTextureUnit(2)

	camera := rkt.NewCamera(1.0, 8000.0, 100.0)
	camera.SetViewport(w, h)
	camera.CaptureMouse(window)

	radius = 10.0

	rkt.LoadPkg("res/base.zip")

	mainVehicle = rkt.NewVehicle("test", rkt.NewPart("base/ctrl1"))
	p := mainVehicle.Parts
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoupa"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid101"))
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoupa"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid102"))
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoupa"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/adapt1015"))
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid153"))
	mainVehicle.Link()

	camera.Target = mainVehicle

	patch := rkt.NewPatch("base/patch/00")
	patch.Scale = 1600.0

	labelWelcome := rkt.NewLabelFor("base/font/anlg", "Welcome")
	labelToStage := rkt.NewLabelFor("base/font/anlg", "press Space to fire the next stage")

	qt := rkt.NewAxisAngleQuat(90.0, rkt.Vec3{X: 1.0, Y: 0.0, Z: 0.0})
	qt = qt.Product(rkt.NewAxisAngleQuat(-90, rkt.Vec3{X: 0.0, Y: 1.0, Z: 0.0}))

	labelWelcome.Scale = rkt.Vec2{X: 6.0, Y: 8.0}
	labelWelcome.Rot = qt
	labelWelcome.Pos.Y = 68.0
	labelWelcome.Pos.Z = -168.0 + 0.1

	labelToStage.Scale = rkt.Vec2{X: 2.5, Y: 5.0}
	labelToStage.Rot = rkt.NewAxisAngleQuat(90.0, rkt.Vec3{X: 1.0, Y: 0.0, Z: 0.0})
	labelToStage.Pos.Y = 25.0
	labelToStage.Pos.Z = -80.0

	n := mainVehicle.Stages
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

		//log.Printf("pz %v vz %v", v.Pos.Z, v.Vel.Z)

		dt := time.Millisecond * 50

		patch.Draw()
		for _, v := range rkt.Vehicles {
			if v == nil {
				break
			}
			v.Draw()
			v.Update(float32(dt.Seconds()))
		}

		// last because transparency
		labelWelcome.Draw()
		labelToStage.Draw()

		time.Sleep(dt)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
