package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type Geom1Def struct {
	TextureName string `json:"texture"`
	Vertices    []Vec3 `json:"vertex"`
	TexCoords   []Vec2 `json:"texcoord"`
}

func (d *Geom1Def) create() *Geom1 {
	g := new(Geom1)
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("create: no such texture %s", d.TextureName)
	}

	count := len(d.Vertices)
	if count != len(d.TexCoords) {
		log.Fatalf("create: vertices and texcoords lens mismatch")
	}

	g.Texture = texture
	g.Vertices = make([]Vec3, count)
	g.TexCoords = make([]Vec2, count)
	copy(g.Vertices, d.Vertices)
	copy(g.TexCoords, d.TexCoords)
	g.Count = count
	return g
}

type Geom1 struct {
	Texture   Texture
	Vertices  []Vec3
	TexCoords []Vec2
	Count     int
}

func NewGeom1(texture Texture, triCount int) *Geom1 {
	g := new(Geom1)
	g.Texture = texture
	g.Count = triCount * 3
	g.Vertices = make([]Vec3, g.Count)
	g.TexCoords = make([]Vec2, g.Count)
	return g
}

func (g *Geom1) clone() *Geom1 {
	n := new(Geom1)
	n.Texture = g.Texture
	n.Count = g.Count
	n.Vertices = make([]Vec3, g.Count)
	n.TexCoords = make([]Vec2, g.Count)
	copy(g.Vertices, n.Vertices)
	copy(g.TexCoords, n.TexCoords)
	return n
}
func (g *Geom1) draw() {
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
func (g *Geom1) drawTexOffset(texCoordOffset Vec2) {
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

type Geom2Def struct {
	Texture0Name string `json:"texture0"`
	Texture1Name string `json:"texture1"`
	Vertices     []Vec3 `json:"vertex"`
	TexCoords0   []Vec2 `json:"texcoord0"`
	TexCoords1   []Vec2 `json:"texcoord1"`
}

func NewGeom2(texture0, texture1 Texture, triCount int) *Geom2 {
	g := new(Geom2)
	g.Texture0 = texture0
	g.Texture1 = texture1
	g.Count = triCount * 3
	g.Vertices = make([]Vec3, g.Count)
	g.TexCoords0 = make([]Vec2, g.Count)
	g.TexCoords1 = make([]Vec2, g.Count)
	return g
}

func (d *Geom2Def) create() *Geom2 {
	g := new(Geom2)
	texture1, ok := textureMap[d.Texture0Name]
	if !ok {
		log.Fatalf("create: no such texture %s", d.Texture0Name)
	}

	texture2, ok := textureMap[d.Texture1Name]
	if !ok {
		log.Fatalf("create: no such texture %s", d.Texture1Name)
	}

	count := len(d.Vertices)
	if count != len(d.TexCoords0) || count != len(d.TexCoords1) {
		log.Fatalf("create: vertices and texcoords lens mismatch")
	}

	g.Texture0 = texture1
	g.Texture1 = texture2
	g.Vertices = make([]Vec3, count)
	g.TexCoords0 = make([]Vec2, count)
	g.TexCoords1 = make([]Vec2, count)
	copy(g.Vertices, d.Vertices)
	copy(g.TexCoords0, d.TexCoords0)
	copy(g.TexCoords1, d.TexCoords1)
	g.Count = count
	return g
}

type Geom2 struct {
	Texture0   Texture
	Texture1   Texture
	Vertices   []Vec3
	TexCoords0 []Vec2
	TexCoords1 []Vec2
	Count      int
}

func (g *Geom2) clone() *Geom2 {
	n := new(Geom2)
	n.Texture0 = g.Texture0
	n.Texture1 = g.Texture1
	n.Count = g.Count
	n.Vertices = make([]Vec3, g.Count)
	n.TexCoords0 = make([]Vec2, g.Count)
	n.TexCoords1 = make([]Vec2, g.Count)
	copy(g.Vertices, n.Vertices)
	copy(g.TexCoords0, n.TexCoords0)
	copy(g.TexCoords1, n.TexCoords1)
	return n
}
func (g *Geom2) draw() {
	g.Texture0.bindTo(0)
	g.Texture1.bindTo(1)
	gl.Begin(gl.TRIANGLES)
	for i := range g.Count {
		v := g.Vertices[i]
		t0 := g.TexCoords0[i]
		t1 := g.TexCoords1[i]
		gl.MultiTexCoord2f(gl.TEXTURE0, t0.X, t0.Y)
		gl.MultiTexCoord2f(gl.TEXTURE1, t1.X, t1.Y)
		gl.Vertex3f(v.X, v.Y, v.Z)
	}

	gl.End()
}
