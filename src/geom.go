package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type GeomDef struct {
	TextureName string `json:"texture"`
	Vertices    []Vec3 `json:"vertex"`
	TexCoords   []Vec2 `json:"texcoord"`
}

func (d *GeomDef) create() *Geom {
	g := new(Geom)
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s", d.TextureName)
	}

	count := len(d.Vertices)
	if count != len(d.TexCoords) {
		log.Fatalf("build_geom: vertices and texcoords lens mismatch")
	}

	g.Texture = texture
	g.Vertices = make([]Vec3, count)
	g.TexCoords = make([]Vec2, count)
	copy(g.Vertices, d.Vertices)
	copy(g.TexCoords, d.TexCoords)
	g.Count = count
	return g
}

type Geom struct {
	Texture   Texture
	Vertices  []Vec3
	TexCoords []Vec2
	Count     int
}

func NewGeom(texture Texture, triCount int) *Geom {
	g := new(Geom)
	g.Texture = texture
	g.Count = triCount * 3
	g.Vertices = make([]Vec3, g.Count)
	g.TexCoords = make([]Vec2, g.Count)
	return g
}

func (g *Geom) clone() *Geom {
	n := new(Geom)
	n.Texture = g.Texture
	n.Count = g.Count
	n.Vertices = make([]Vec3, g.Count)
	n.TexCoords = make([]Vec2, g.Count)
	copy(g.Vertices, n.Vertices)
	copy(g.TexCoords, n.TexCoords)
	return n
}

func (g *Geom) draw() {
	g.Texture.bind()
	gl.Begin(gl.TRIANGLES)
	for i := range g.Count {
		v := g.Vertices[i]
		t := g.TexCoords[i]
		gl.TexCoord2f(t.X, t.Y)
		gl.Vertex3f(v.X, v.Y, v.Z)
	}

	gl.End()
}
func (g *Geom) drawTexOffset(texCoordOffset Vec2) {
	g.Texture.bind()
	gl.Begin(gl.TRIANGLES)
	for i := range g.Count {
		v := g.Vertices[i]
		t := g.TexCoords[i]
		gl.TexCoord2f(t.X+texCoordOffset.X, t.Y+texCoordOffset.Y)
		gl.Vertex3f(v.X, v.Y, v.Z)
	}

	gl.End()
}
