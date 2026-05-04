package rkt

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Geom1Def struct {
	ShaderName  string
	TextureName string
	BufferAttrs []BufferAttr
	RawArray    []float32
}

var geom2BufferAttrs = []BufferAttr{
	{BufferAttrPos, "a_Pos", 3},
	{BufferAttrUV0, "a_UV0", 2},
	{BufferAttrUV1, "a_UV1", 2},
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

	g.Shader = shader
	g.Texture = texture
	g.Buffer = NewBuffer(shader, d.BufferAttrs)
	g.Buffer.data(d.RawArray)

	return g
}

type Geom1 struct {
	Shader  Shader
	Texture Texture
	Buffer  *Buffer
}

func (g *Geom1) clone() *Geom1 {
	n := new(Geom1)
	n.Shader = g.Shader
	n.Texture = g.Texture
	n.Buffer = g.Buffer
	return n
}
func (g *Geom1) draw(m *Matrix4) {
	g.Shader.active()
	g.Texture.bind()
	uDiffTexture := g.Shader.getUniform("u_DiffTexture")
	uAmbLightColor := g.Shader.getUniform("u_AmbLightColor")
	uDirLightDir := g.Shader.getUniform("u_DirLightDir")
	uDirLightColor := g.Shader.getUniform("u_DirLightColor")
	uPMatrix := g.Shader.getUniform("u_PMatrix")
	uVMatrix := g.Shader.getUniform("u_VMatrix")
	uMMatrix := g.Shader.getUniform("u_MMatrix")
	// TODO: source light from somewhere else
	dirLightDir := []float32{0.7, 0.3, 0.5, 0.0, 0.0, -1.0}
	dirLightColor := []float32{0.9, 0.9, 1.0, 0.0, 0.2, 0.3}
	gl.Uniform3f(uAmbLightColor, 0.4, 0.4, 0.5)
	gl.Uniform3fv(uDirLightDir, 2, &dirLightDir[0])
	gl.Uniform3fv(uDirLightColor, 2, &dirLightColor[0])
	g.Texture.uniform(uDiffTexture, 0)
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
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
	ShaderName   string `json:"shader"`
	Texture0Name string `json:"texture0"`
	Texture1Name string `json:"texture1"`
	Vertices     []Vec3 `json:"vertex"`
	TexCoords0   []Vec2 `json:"texcoord0"`
	TexCoords1   []Vec2 `json:"texcoord1"`
}

func (d *Geom2Def) create() *Geom2 {
	g := new(Geom2)
	shader, ok := shaderMap[d.ShaderName]
	if !ok {
		log.Fatalf("create: no such shader %s", d.ShaderName)
	}

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

	g.Shader = shader
	g.Texture0 = texture0
	g.Texture1 = texture1
	g.Buffer = NewBuffer(shader, geom2BufferAttrs)
	bufferData := g.Buffer.newDataArray(count)
	i := int32(0)
	for j := range count {
		bufferData[i+0] = d.Vertices[j].X
		bufferData[i+1] = d.Vertices[j].Y
		bufferData[i+2] = d.Vertices[j].Z
		i += g.Buffer.Attrs[0].Cocnt
		bufferData[i+0] = d.TexCoords0[j].X
		bufferData[i+1] = d.TexCoords0[j].Y
		i += g.Buffer.Attrs[1].Cocnt
		bufferData[i+0] = d.TexCoords1[j].X
		bufferData[i+1] = d.TexCoords1[j].Y
		i += g.Buffer.Attrs[2].Cocnt
	}
	g.Buffer.data(bufferData)

	return g
}

type Geom2 struct {
	Shader   Shader
	Texture0 Texture
	Texture1 Texture
	Buffer   *Buffer
	Count    int // unused
}

func (g *Geom2) clone() *Geom2 {
	n := new(Geom2)
	n.Shader = g.Shader
	n.Texture0 = g.Texture0
	n.Texture1 = g.Texture1
	n.Buffer = g.Buffer
	n.Count = g.Count
	return n
}
func (g *Geom2) draw(m *Matrix4) {
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

	g.Shader.active()
	g.Texture0.bindTo(0)
	g.Texture1.bindTo(1)
	uDiffTexture0 := g.Shader.getUniform("u_DiffTexture0")
	uDiffTexture1 := g.Shader.getUniform("u_DiffTexture1")
	uPMatrix := g.Shader.getUniform("u_PMatrix")
	uVMatrix := g.Shader.getUniform("u_VMatrix")
	uMMatrix := g.Shader.getUniform("u_MMatrix")
	g.Texture0.uniform(uDiffTexture0, 0)
	g.Texture1.uniform(uDiffTexture1, 1)
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
	m.uniform(uMMatrix)
	g.Buffer.bind()
	g.Buffer.draw()

	// gl.End()
}
