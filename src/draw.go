package rkt

import "github.com/go-gl/gl/v2.1/gl"

func drawAttachPt(size float32, pos *Vec3) {
	w, b := float32(1.0), float32(0.0)

	gl.Begin(gl.TRIANGLES)

	gl.Color4f(w, w, w, 0.0)
	gl.Vertex3f(pos.x+size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x, pos.y+size, pos.z)
	gl.Vertex3f(pos.x+size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y+size, pos.z)

	gl.Color4f(b, b, b, 0.0)
	gl.Vertex3f(pos.x+size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x+size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y+size, pos.z)
	gl.Vertex3f(pos.x-size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y+size, pos.z)

	gl.Color4f(w, w, w, 0.0)
	gl.Vertex3f(pos.x+size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x, pos.y-size, pos.z)
	gl.Vertex3f(pos.x+size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y-size, pos.z)

	gl.Color4f(b, b, b, 0.0)
	gl.Vertex3f(pos.x+size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x+size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y-size, pos.z)
	gl.Vertex3f(pos.x-size, pos.y, pos.z+size)
	gl.Vertex3f(pos.x-size, pos.y, pos.z-size)
	gl.Vertex3f(pos.x, pos.y-size, pos.z)

	gl.End()
}

type Geom struct {
	texture   Texture
	vertices  []Vec3
	texCoords []Vec2
	count     uint32
}

func (g *Geom) draw() {
	g.texture.bind()
	gl.Begin(gl.TRIANGLES)
	for i := range g.count {
		v := g.vertices[i]
		t := g.texCoords[i]
		gl.TexCoord2f(t.x, t.y)
		gl.Vertex3f(v.x, v.y, v.z)
	}

	gl.End()
}
