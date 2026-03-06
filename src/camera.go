package rkt

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Camera struct {
	Target       *Vehicle
	FocusPos     Vec3
	lastMousePos Vec2
	mouseSpeed   float32
	depthNear    float64
	depthFar     float64
	width        uint16
	height       uint16
	Radius       float32
	pitch        float32
	yaw          float32
}

func NewCamera(depthNear, depthFar float64, mouseSpeed float32) *Camera {
	c := new(Camera)
	c.depthNear = depthNear
	c.depthFar = depthFar
	c.mouseSpeed = mouseSpeed
	c.Radius = 10.0

	return c
}

func (c *Camera) SetViewport(width, height uint16) {
	c.width = width
	c.height = height
	c.SetProjection()
}
func (c *Camera) SetProjection() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.LoadIdentity()
	f := float64(c.width)/float64(c.height) - 1.0
	gl.Frustum(-1.0-f, 1.0+f, -1.0, 1.0, c.depthNear, c.depthFar)
	gl.Rotatef(90.0, 0.0, 0.0, 1.0) // y <-> z
	gl.Rotatef(90.0, 0.0, 1.0, 0.0) // y <-> z
	gl.PushMatrix()
}
func (c *Camera) CaptureMouse(window *glfw.Window) {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
func (c *Camera) Apply() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PushMatrix()
	gl.Translatef(c.Radius, 0.0, 0.0)
	gl.Rotatef(-c.pitch, 0.0, 1.0, 0.0)
	gl.Rotatef(c.yaw, 0.0, 0.0, 1.0)
	c.FocusPos.Apply()
}
func (c *Camera) Update(mousePos Vec2) {
	if c.Target != nil {
		c.FocusPos.X = -c.Target.Pos.X
		c.FocusPos.Y = -c.Target.Pos.Y
		c.FocusPos.Z = -c.Target.Pos.Z
	}

	diffPos := mousePos.Sub(c.lastMousePos)
	c.lastMousePos = mousePos
	c.yaw += diffPos.X / float32(c.width) * c.mouseSpeed
	c.pitch += diffPos.Y / float32(c.height) * c.mouseSpeed
	c.pitch = min(max(c.pitch, -90.0), 90.0)
}
