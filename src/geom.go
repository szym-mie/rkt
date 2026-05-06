package rkt

type Geom1Def struct {
	ShaderName  string
	TextureName string
	BufferAttrs []BufferAttr
	Array       []float32
}

func (d *Geom1Def) create() *Geom1 {
	g := new(Geom1)
	g.Shader = getShader(d.ShaderName)
	g.Texture = getTexture(d.TextureName)
	g.Buffer = NewBuffer(g.Shader, d.BufferAttrs)
	g.Buffer.arrayFloat(d.Array)
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
	g.Texture.bind2D()
	uTexture := g.Shader.getUniform("u_Texture")
	uAmbLightColor := g.Shader.getUniform("u_AmbLightColor")
	uDirLightDir := g.Shader.getUniform("u_DirLightDir")
	uDirLightColor := g.Shader.getUniform("u_DirLightColor")
	uPMatrix := g.Shader.getUniform("u_PMatrix")
	uVMatrix := g.Shader.getUniform("u_VMatrix")
	uMMatrix := g.Shader.getUniform("u_MMatrix")

	ActiveLightEnv.uniform(uAmbLightColor, uDirLightDir, uDirLightColor)
	g.Texture.uniform(uTexture, 0)
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
	Texture3Name string
	BufferAttrs  []BufferAttr
	Array        []float32
}

func (d *Geom2Def) create() *Geom2 {
	g := new(Geom2)
	g.Shader = getShader(d.ShaderName)
	g.Texture0 = getTexture(d.Texture0Name)
	g.Texture1 = getTexture(d.Texture1Name)
	g.Texture2 = getTexture(d.Texture2Name)
	g.Texture3 = getTexture(d.Texture3Name)
	g.Buffer = NewBuffer(g.Shader, d.BufferAttrs)
	g.Buffer.arrayFloat(d.Array)
	return g
}

// Well it should be Geom3 but I don't feel like changing every name now
type Geom2 struct {
	Shader   Shader
	Texture0 Texture
	Texture1 Texture
	Texture2 Texture
	Texture3 Texture
	Buffer   *Buffer
}

func (g *Geom2) clone() *Geom2 {
	n := new(Geom2)
	n.Shader = g.Shader
	n.Texture0 = g.Texture0
	n.Texture1 = g.Texture1
	n.Texture2 = g.Texture2
	n.Texture3 = g.Texture3
	n.Buffer = g.Buffer
	return n
}
func (g *Geom2) draw(m *Matrix4) {
	g.Shader.active()
	g.Texture0.bindTo2D(0)
	g.Texture1.bindTo2D(1)
	g.Texture1.bindTo2D(2)
	uTexture0 := g.Shader.getUniform("u_Texture0")
	uTexture1 := g.Shader.getUniform("u_Texture1")
	uTexture2 := g.Shader.getUniform("u_Texture2")
	uTexture3 := g.Shader.getUniform("u_Texture3")
	uAmbLightColor := g.Shader.getUniform("u_AmbLightColor")
	uDirLightDir := g.Shader.getUniform("u_DirLightDir")
	uDirLightColor := g.Shader.getUniform("u_DirLightColor")
	uPMatrix := g.Shader.getUniform("u_PMatrix")
	uVMatrix := g.Shader.getUniform("u_VMatrix")
	uMMatrix := g.Shader.getUniform("u_MMatrix")

	ActiveLightEnv.uniform(uAmbLightColor, uDirLightDir, uDirLightColor)
	g.Texture0.uniform(uTexture0, 0)
	g.Texture1.uniform(uTexture1, 1)
	g.Texture2.uniform(uTexture2, 2)
	g.Texture3.uniform(uTexture3, 3)
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
	m.uniform(uMMatrix)
	g.Buffer.bind()
	g.Buffer.draw()
}
