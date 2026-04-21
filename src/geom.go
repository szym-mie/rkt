package rkt

import (
	"log"
)

type Geom1Def struct {
	ShaderName  string `json:"shader"`
	TextureName string `json:"texture"`
	Vertices    []Vec3 `json:"vertex"`
	TexCoords   []Vec2 `json:"texcoord"`
}

var geom1BufferAttrs = []BufferAttr{
	{BufferAttrPos, "a_Pos", 3},
	{BufferAttrTexCoord0, "a_TexCoord0", 2},
}

func (d *Geom1Def) create() *Geom1 {
	g := new(Geom1)
	shader, ok := shaderMap[d.ShaderName]
	if !ok {
		log.Fatalf("create: no such shader %s", d.ShaderName)
	}

	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("create: no such texture %s", d.TextureName)
	}

	count := len(d.Vertices)
	if count != len(d.TexCoords) {
		log.Fatalf("create: vertices and texcoords lens mismatch")
	}

	g.Shader = shader
	g.Texture = texture
	g.Buffer = NewBuffer(shader, geom1BufferAttrs)
	bufferData := g.Buffer.newDataArray(count)
	i := int32(0)
	for j := range count {
		bufferData[i+0] = d.Vertices[j].X
		bufferData[i+1] = d.Vertices[j].Y
		bufferData[i+2] = d.Vertices[j].Z
		i += g.Buffer.Attrs[0].Size
		bufferData[i+0] = d.TexCoords[j].X
		bufferData[i+1] = d.TexCoords[j].Y
		i += g.Buffer.Attrs[1].Size
	}
	g.Buffer.data(bufferData)

	return g
}

type Geom1 struct {
	Shader  Shader
	Texture Texture
	Buffer  *Buffer
	Count   int // unused
}

func NewGeom1(shader Shader, texture Texture) *Geom1 {
	g := new(Geom1)
	g.Shader = shader
	g.Texture = texture
	g.Buffer = NewBuffer(shader, geom1BufferAttrs)
	return g
}

func (g *Geom1) clone() *Geom1 {
	n := new(Geom1)
	n.Shader = g.Shader
	n.Texture = g.Texture
	n.Buffer = g.Buffer
	n.Count = g.Count
	return n
}
func (g *Geom1) draw(m *Matrix4) {
	g.Shader.active()
	g.Texture.bind()
	uDiffTexture := g.Shader.getUniform("u_DiffTexture")
	uVPMatrix := g.Shader.getUniform("u_VPMatrix")
	uMMatrix := g.Shader.getUniform("u_MMatrix")
	g.Texture.uniform(uDiffTexture, 0)
	ActiveCamera.pvMatrix.uniform(uVPMatrix)
	m.uniform(uMMatrix)
	g.Buffer.bind()
	g.Buffer.draw()
}
func (g *Geom1) drawTexOffset(texCoordOffset Vec2) {
	// g.Texture.bind()
	// for i := range g.Count {
	// 	v := g.Vertices[i]
	// 	t := g.TexCoords[i]
	// 	gl.TexCoord2f(t.X+texCoordOffset.X, t.Y+texCoordOffset.Y)
	// 	gl.Vertex3f(v.X, v.Y, v.Z)
	// }
	// TODO: bind and draw
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
	texture0, ok := textureMap[d.Texture0Name]
	if !ok {
		log.Fatalf("create: no such texture %s", d.Texture0Name)
	}

	texture1, ok := textureMap[d.Texture1Name]
	if !ok {
		log.Fatalf("create: no such texture %s", d.Texture1Name)
	}

	count := len(d.Vertices)
	if count != len(d.TexCoords0) || count != len(d.TexCoords1) {
		log.Fatalf("create: vertices and texcoords lens mismatch")
	}

	g.Texture0 = texture0
	g.Texture1 = texture1
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
	// g.Texture0.bindTo(0)
	// g.Texture1.bindTo(1)
	// g.Texture0.setFilter(TextureFilterLinear)
	// gl.Begin(gl.TRIANGLES)
	// for i := range g.Count {
	// 	v := g.Vertices[i]
	// 	t0 := g.TexCoords0[i]
	// 	t1 := g.TexCoords1[i]
	// 	gl.MultiTexCoord2f(gl.TEXTURE0, t0.X, t0.Y)
	// 	gl.MultiTexCoord2f(gl.TEXTURE1, t1.X, t1.Y)
	// 	gl.Vertex3f(v.X, v.Y, v.Z)
	// }

	// gl.End()
}
