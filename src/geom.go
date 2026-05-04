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

type Geom2Def struct {
	ShaderName   string
	Texture0Name string
	Texture1Name string
	Texture2Name string
	BufferAttrs  []BufferAttr
	RawArray     []float32
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

	texture2, ok := textureMap[d.Texture2Name]
	if !ok {
		log.Fatalf("create: no such texture %s", d.Texture2Name)
	}

	g.Shader = shader
	g.Texture0 = texture0
	g.Texture1 = texture1
	g.Texture2 = texture2
	g.Buffer = NewBuffer(shader, d.BufferAttrs)
	g.Buffer.data(d.RawArray)
	return g
}

// Well it should be Geom3 but I don't feel like changing every name now
type Geom2 struct {
	Shader   Shader
	Texture0 Texture
	Texture1 Texture
	Texture2 Texture
	Buffer   *Buffer
}

func (g *Geom2) clone() *Geom2 {
	n := new(Geom2)
	n.Shader = g.Shader
	n.Texture0 = g.Texture0
	n.Texture1 = g.Texture1
	n.Texture2 = g.Texture2
	n.Buffer = g.Buffer
	return n
}
func (g *Geom2) draw(m *Matrix4) {
	g.Shader.active()
	g.Texture0.bindTo(0)
	g.Texture1.bindTo(1)
	g.Texture1.bindTo(2)
	// TODO: use less descriptive texture names in shaders
	uDiffTexture0 := g.Shader.getUniform("u_DiffTexture0")
	uDiffTexture1 := g.Shader.getUniform("u_DiffTexture1")
	uNormTexture := g.Shader.getUniform("u_NormTexture")
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
	g.Texture0.uniform(uDiffTexture0, 0)
	g.Texture1.uniform(uDiffTexture1, 1)
	g.Texture2.uniform(uNormTexture, 2)
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
	m.uniform(uMMatrix)
	g.Buffer.bind()
	g.Buffer.draw()
}
