package rkt

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type CameraMode uint

const (
	CameraOrbitMode CameraMode = iota + 1
	CameraFlyByMode
	CameraFixedMode
)

type Camera struct {
	PVMatrix
	Target       *RigidBody
	FocusPos     Vec3
	PivotPos     Vec3
	lastMousePos Vec2
	mode         CameraMode
	mouseSpeed   float32
	depthNear    float32
	depthFar     float32
	width        uint16
	height       uint16
	Radius       float32
	pitch        float32
	yaw          float32
}

var axisSwapMatrix Mat4

func InitCamera() {
	axisSwapMatrix.SetIdentity()
	axisSwapMatrix.RotZ(math.Pi * 0.5)
	axisSwapMatrix.RotY(math.Pi * 0.5)
}

func NewCamera(depthNear, depthFar float32, mouseSpeed float32) *Camera {
	c := new(Camera)
	c.mode = CameraOrbitMode
	c.depthNear = depthNear
	c.depthFar = depthFar
	c.mouseSpeed = mouseSpeed
	c.Radius = 10.0
	return c
}

func (c *Camera) ToOrbit() {
	c.mode = CameraOrbitMode
}
func (c *Camera) ToFlyBy() {
	c.mode = CameraFlyByMode
	if c.Target != nil {
		c.PivotPos.From(c.Target.Pos)
		c.PivotPos.AddSelf(c.Target.Vel.MulSca(2))
	}
}
func (c *Camera) ToFixed() {
	c.mode = CameraFixedMode
}
func (c *Camera) SetViewport(width, height uint16) {
	c.width = width
	c.height = height
	c.SetProjection()
}
func (c *Camera) SetProjection() {
	aspect := float32(c.width) / float32(c.height)
	c.ProjMatrix.Frustum(aspect, c.depthNear, c.depthFar)
}
func (c *Camera) CaptureMouse(window *glfw.Window) {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
func (c *Camera) UpdateViewOrbit() {
	c.ViewMatrix.From(&axisSwapMatrix)
	c.ViewMatrix.AddPosSelf(Vec3{0, 0, -c.Radius})
	c.ViewMatrix.RotY(-c.pitch)
	c.ViewMatrix.RotZ(c.yaw)
	focusMatrix := NewMat4Pos(c.FocusPos.Invert())
	c.ViewMatrix.MulSelf(focusMatrix)
}
func (c *Camera) UpdateViewFlyBy() {
	c.ViewMatrix.From(&axisSwapMatrix)
	lookDir := c.FocusPos.Sub(c.PivotPos).Norm()
	pitch := Asinf(lookDir.Z)
	yaw := Atan2f(lookDir.Y, lookDir.X)
	c.ViewMatrix.RotY(pitch)
	c.ViewMatrix.RotZ(-yaw)
	pivotMatrix := NewMat4Pos(c.PivotPos.Invert())
	c.ViewMatrix.MulSelf(pivotMatrix)
}
func (c *Camera) UpdateFixedMode() {
	c.ViewMatrix.From(&axisSwapMatrix)
	c.ViewMatrix.AddPosSelf(Vec3{c.yaw, -c.pitch, -c.Radius})
	if c.Target != nil {
		c.ViewMatrix.MulSelf(c.Target.Rot.Conj().ToMat4())
	}
	focusMatrix := NewMat4Pos(c.FocusPos.Invert())
	c.ViewMatrix.MulSelf(focusMatrix)
}
func (c *Camera) Update(mousePos Vec2) {
	if c.Target != nil {
		c.FocusPos.From(c.Target.Pos)
	}

	diffPos := mousePos.Sub(c.lastMousePos)
	c.lastMousePos = mousePos
	c.yaw += diffPos.X / float32(c.width) * c.mouseSpeed
	c.pitch += diffPos.Y / float32(c.height) * c.mouseSpeed
	c.pitch = min(max(c.pitch, -math.Pi*0.5), math.Pi*0.5)
	c.Radius = Maxf(c.Radius, 0.5)
	switch c.mode {
	case CameraOrbitMode:
		c.UpdateViewOrbit()
	case CameraFlyByMode:
		c.UpdateViewFlyBy()
	case CameraFixedMode:
		c.UpdateFixedMode()
	}
}
