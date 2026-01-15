package rkt

import "github.com/go-gl/gl/v2.1/gl"

type Plume struct {
	texture Texture
	offset  Vec3
	size    float32
}

func (p *Plume) draw() {
	p.texture.bind()
	p.offset.apply()

	gl.Begin(gl.TRIANGLES)

	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, p.size)
	gl.TexCoord2f(0.5, 0.0)
	gl.Vertex3f(0.0, 0.0, 0.0)
	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, -p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 0.0)
	gl.Vertex3f(0.0, 0.0, 0.0)

	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 0.0)
	gl.Vertex3f(-p.size, -1.0, p.size)
	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(0.0, 0.0, 0.0)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 0.0)
	gl.Vertex3f(0.0, 0.0, 0.0)

	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, p.size)
	gl.TexCoord2f(0.5, 1.0)
	gl.Vertex3f(0.0, -4.0, 0.0)
	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, -p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 1.0)
	gl.Vertex3f(0.0, -4.0, 0.0)

	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(p.size, -1.0, p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 1.0)
	gl.Vertex3f(0.0, -4.0, 0.0)
	gl.TexCoord2f(0.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, p.size)
	gl.TexCoord2f(1.0, 0.5)
	gl.Vertex3f(-p.size, -1.0, -p.size)
	gl.TexCoord2f(0.5, 1.0)
	gl.Vertex3f(0.0, -4.0, 0.0)

	gl.End()
}
func (p *Plume) update() {

}
