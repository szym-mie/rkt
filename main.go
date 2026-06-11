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

var pause bool = true
var focusIndex uint
var camera *rkt.Camera

// var mainVehicle *rkt.Vehicle
var airplane *rkt.L39
var speedup int8 = 4
var radius float32 = 10
var radiusCh float32

var pitch, roll, yaw, throttle float32
var pitchCh, rollCh, yawCh, throttleCh float32

type onKeyAction func(w *glfw.Window)

var onKeyPressMap = map[glfw.Key]onKeyAction{
	glfw.KeyEscape:       func(w *glfw.Window) { w.SetShouldClose(true) },
	glfw.KeyP:            func(_ *glfw.Window) { pause = !pause },
	glfw.KeyLeftBracket:  func(_ *glfw.Window) { radiusCh = -0.8 },
	glfw.KeyRightBracket: func(_ *glfw.Window) { radiusCh = +1.0 },
	glfw.KeyEqual:        func(_ *glfw.Window) { throttleCh = +1.0 },
	glfw.KeyMinus:        func(_ *glfw.Window) { throttleCh = -1.0 },
	glfw.KeyA:            func(_ *glfw.Window) { rollCh = +1.0 },
	glfw.KeyD:            func(_ *glfw.Window) { rollCh = -1.0 },
	glfw.KeyW:            func(_ *glfw.Window) { pitchCh = +1.0 },
	glfw.KeyS:            func(_ *glfw.Window) { pitchCh = -1.0 },
	glfw.KeyX:            func(_ *glfw.Window) { yawCh = +1.0 },
	glfw.KeyZ:            func(_ *glfw.Window) { yawCh = -1.0 },
	glfw.KeyF2:           func(_ *glfw.Window) { camera.ToOrbit() },
	glfw.KeyF3:           func(_ *glfw.Window) { camera.ToFlyBy() },
	glfw.KeyF4:           func(_ *glfw.Window) { camera.ToFixed() },
}

var onKeyReleaseMap = map[glfw.Key]onKeyAction{
	glfw.KeyEqual:        func(_ *glfw.Window) { throttleCh = 0.0 },
	glfw.KeyMinus:        func(_ *glfw.Window) { throttleCh = 0.0 },
	glfw.KeyLeftBracket:  func(_ *glfw.Window) { radiusCh = 0.0 },
	glfw.KeyRightBracket: func(_ *glfw.Window) { radiusCh = 0.0 },
	glfw.KeyA:            func(_ *glfw.Window) { rollCh = 0.0 },
	glfw.KeyD:            func(_ *glfw.Window) { rollCh = 0.0 },
	glfw.KeyW:            func(_ *glfw.Window) { pitchCh = 0.0 },
	glfw.KeyS:            func(_ *glfw.Window) { pitchCh = 0.0 },
	glfw.KeyX:            func(_ *glfw.Window) { yawCh = 0.0 },
	glfw.KeyZ:            func(_ *glfw.Window) { yawCh = 0.0 },
}

func onKey(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, mods glfw.ModifierKey) {
	switch act {
	case glfw.Press:
		actFunc, ok := onKeyPressMap[key]
		if ok {
			actFunc(w)
		}
	case glfw.Release:
		actFunc, ok := onKeyReleaseMap[key]
		if ok {
			actFunc(w)
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

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	rkt.InitTextureUnit(0)
	rkt.InitTextureUnit(1)
	rkt.InitTextureUnit(2)
	rkt.InitTextureUnit(3)

	rkt.LoadPkg("res/base.zip")

	sky := rkt.NewSky()

	hud := rkt.NewHud()
	hud.SetViewport(w, h)

	rkt.InitCamera()
	camera = rkt.NewCamera(1.0, 8000.0, 10.0)
	camera.SetViewport(w, h)
	camera.CaptureMouse(window)

	rkt.ActiveLightEnv.AmbColor = rkt.Vec3{X: 0.1, Y: 0.1, Z: 0.15}
	rkt.ActiveLightEnv.DirLights[0] = rkt.DirLight{
		Dir:   rkt.Vec3{X: 0.7, Y: 0.3, Z: 0.5},
		Color: rkt.Vec3{X: 0.9, Y: 0.9, Z: 1.0}}
	rkt.ActiveLightEnv.DirLights[1] = rkt.DirLight{
		Dir:   rkt.Vec3{X: 0.0, Y: 0.0, Z: -1.0},
		Color: rkt.Vec3{X: 0.0, Y: 0.1, Z: 0.1}}

	radius = 10.0

	airplane = rkt.NewL39()
	airplane.Pos.Z = 100.0
	airplane.Vel.X = 100.0

	patch00 := rkt.NewPatch("base/patch/00")
	patch00.Scale = 1600.0

	t := time.Time{}.Add(time.Hour * 14)

	rkt.InitDraw()
	rkt.SetLineColor(1.0, 0.0, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera.Radius = radius
		camera.Target = &airplane.RigidBody

		dt := time.Millisecond * 25
		dtf := float32(dt.Seconds())

		radius *= 1.0 + radiusCh*dtf
		if !pause {
			pitch += pitchCh * dtf
			pitch = rkt.Clampf(pitch, -1.0, +1.0)
			yaw += yawCh * dtf
			yaw = rkt.Clampf(yaw, -1.0, +1.0)
			throttle += throttleCh * dtf
			throttle = rkt.Clampf(throttle, 0.0, +1.0)
			if rollCh != 0.0 {
				roll += rollCh * dtf
				roll = rkt.Clampf(roll, -1.0, +1.0)
			} else {
				if roll < -0.5 {
					roll += dtf
				}
				if roll > +0.5 {
					roll -= dtf
				}
			}
			airplane.Wings[2].SetDeflection(-roll)
			airplane.Wings[3].SetDeflection(+roll)
			airplane.Wings[4].SetDeflection(+pitch)
			airplane.Wings[5].SetDeflection(+pitch)
			airplane.Wings[6].SetDeflection(+yaw)
			airplane.Engine.Throttle = throttle

			airplane.Update(dtf)
			log.Printf("v = %3.0f km/h    h = %5.0f m\n", airplane.Vel.Len()*3.6, airplane.Pos.Z)

			t = t.Add(dt * 60)
			rkt.ActiveLightEnv.SetFromExtern(t, rkt.Vec3{})
		}

		x, y := window.GetCursorPos()
		mousePos := rkt.Vec2{X: float32(x), Y: float32(y)}
		camera.SetProjection()
		camera.Update(mousePos)
		rkt.ActivePV = &camera.PVMatrix

		gl.DepthFunc(gl.LEQUAL)
		sky.Draw()
		gl.DepthFunc(gl.LESS)

		patch00.Draw()
		airplane.Draw()

		rkt.ActivePV = &hud.PVMatrix
		hud.Draw(airplane.Rot, rkt.Vec4{X: pitch, Y: roll, Z: yaw, W: throttle})

		time.Sleep(dt)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
