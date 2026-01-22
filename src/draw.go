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
