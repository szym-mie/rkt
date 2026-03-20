package rkt

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func drawAttachPt(size float32, pos *Vec3) {
	w, b := float32(1.0), float32(0.0)

	gl.Begin(gl.TRIANGLES)

	gl.Color4f(w, w, w, 0.0)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X, pos.Y+size, pos.Z)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y+size, pos.Z)

	gl.Color4f(b, b, b, 0.0)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y+size, pos.Z)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y+size, pos.Z)

	gl.Color4f(w, w, w, 0.0)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X, pos.Y-size, pos.Z)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y-size, pos.Z)

	gl.Color4f(b, b, b, 0.0)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X+size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y-size, pos.Z)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z+size)
	gl.Vertex3f(pos.X-size, pos.Y, pos.Z-size)
	gl.Vertex3f(pos.X, pos.Y-size, pos.Z)

	gl.End()
}

var lineForceGeom *Geom1
var lineColor Vec3

func InitDraw() {
	lineForceGeom = geom1DefMap["base/geom/forcevec"].create()
}

func SetLineColor(r, g, b float32) {
	lineColor = Vec3{r, g, b}
}

func DrawForce(dir, pt Vec3) {
	gl.Disable(gl.DEPTH_TEST)
	gl.PushMatrix()
	gl.Color4f(lineColor.X, lineColor.Y, lineColor.Z, 1.0)
	q := NewVecDiffQuat(dir.Norm(), Vec3{0, 0, 1}).Norm()
	pt.Apply()
	q.Conj().Apply()
	gl.Scalef(1.0, 1.0, dir.Len())
	lineForceGeom.draw()
	gl.Color4f(1.0, 1.0, 1.0, 1.0)
	gl.PopMatrix()
	gl.Enable(gl.DEPTH_TEST)
}
