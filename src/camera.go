package rkt

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var ActiveCamera *Camera

type Camera struct {
	Target       *Vehicle
	FocusPos     Vec3
	lastMousePos Vec2
	mouseSpeed   float32
	depthNear    float32
	depthFar     float32
	width        uint16
	height       uint16
	projMatrix   Matrix4
	viewMatrix   Matrix4
	pvMatrix     Matrix4
	Radius       float32
	pitch        float32
	yaw          float32
}

func NewCamera(depthNear, depthFar float32, mouseSpeed float32) *Camera {
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
	aspect := float32(c.width) / float32(c.height)
	c.projMatrix.Frustum(aspect, c.depthNear, c.depthFar)
	c.updatePv()
}
func (c *Camera) CaptureMouse(window *glfw.Window) {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
func (c *Camera) UpdateView() {
	c.viewMatrix.SetIdentity()
	c.viewMatrix.SetPos(Vec3{0.0, 0.0, -c.Radius})
	c.viewMatrix.RotZ(math.Pi * 0.5)
	c.viewMatrix.RotY(math.Pi * 0.5)
	c.viewMatrix.RotY(-c.pitch)
	c.viewMatrix.RotZ(c.yaw)
	trans := NewMatrix4Pos(c.FocusPos)
	c.viewMatrix.MulSelf(trans)
	c.updatePv()
}
func (c *Camera) updatePv() {
	c.pvMatrix = *c.projMatrix.Mul(&c.viewMatrix)
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
	c.pitch = min(max(c.pitch, -math.Pi*0.5), math.Pi*0.5)
	c.UpdateView()
}
