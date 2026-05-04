package main

import (
	"log"
	"runtime"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v3.3-core/gl"
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

var pause bool
var focusIndex uint
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
		case glfw.KeyP:
			pause = !pause
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
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.ZAxis, -ctrlSpeed*10)))
		case glfw.KeyE:
			mainVehicle.Ang = mainVehicle.Ang.Add(mainVehicle.Rot.Rotate(ctrlVector(rkt.ZAxis, +ctrlSpeed*10)))
		case glfw.KeyF2:
			focusIndex = (focusIndex + 1) % rkt.VehiclesIndex
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

	// fog := [4]float32{0.2, 0.7, 0.8, 0.0}
	// gl.Enable(gl.FOG)
	// gl.Fogi(gl.FOG_MODE, gl.LINEAR)
	// gl.Fogi(gl.FOG_COORD_SRC, gl.FRAGMENT_DEPTH)
	// gl.Fogfv(gl.FOG_COLOR, &fog[0])
	// gl.Fogf(gl.FOG_DENSITY, 0.05)
	// gl.Hint(gl.FOG_HINT, gl.NICEST)
	// gl.Fogf(gl.FOG_START, 1.0)
	// gl.Fogf(gl.FOG_END, 5000.0)

	// fp, _ := os.Open("bmlcube.bml")
	// bml, err := rkt.ReadBML(fp)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println(bml.Header.ElemCount)
	// 	log.Println(bml.Header.Externs)
	// 	log.Println(bml.Header.Attribs)
	// }

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	rkt.InitTextureUnit(0)
	rkt.InitTextureUnit(1)
	rkt.InitTextureUnit(2)

	rkt.LoadPkg("res/base.zip")

	hud := rkt.NewHud()
	hud.SetViewport(w, h)

	camera := rkt.NewCamera(1.0, 8000.0, 10.0)
	camera.SetViewport(w, h)
	camera.CaptureMouse(window)

	radius = 10.0

	mainVehicle = rkt.NewVehicle("test", rkt.NewPart("base/pod10"))
	p := mainVehicle.Parts
	mainVehicle.AttachAbove(p, rkt.NewPart("base/chute05"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoup10"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid101"))
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoup10"))
	mainVehicle.AddStage()
	p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid101"))
	// p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoup10"))
	// mainVehicle.AddStage()
	// p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid102"))
	// p = mainVehicle.AttachBelow(p, rkt.NewPart("base/decoup10"))
	// mainVehicle.AddStage()
	// p = mainVehicle.AttachBelow(p, rkt.NewPart("base/adapt1015"))
	// p = mainVehicle.AttachBelow(p, rkt.NewPart("base/solid153"))
	mainVehicle.Link()

	camera.Target = mainVehicle

	patch00 := rkt.NewPatch("base/patch/00")
	patch00.Scale = 1600.0
	// patchInf := rkt.NewPatch("base/patch/inf")
	// patchInf.Scale = 1600.0

	rkt.InitDraw()
	rkt.SetLineColor(1.0, 0.0, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera.Radius = radius
		camera.Target = rkt.Vehicles[focusIndex]

		x, y := window.GetCursorPos()
		mousePos := rkt.Vec2{X: float32(x), Y: float32(y)}
		camera.SetProjection()
		camera.Update(mousePos)
		rkt.ActivePV = &camera.PVMatrixPair

		dt := time.Millisecond * 25

		patch00.Draw()
		// patchInf.Draw()
		for _, v := range rkt.Vehicles {
			if v == nil {
				break
			}
			v.Draw()
			if !pause {
				v.Update(float32(dt.Seconds()))
			}
		}

		rkt.ActivePV = &hud.PVMatrixPair
		hud.Draw(mainVehicle.Rot)

		time.Sleep(dt)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
