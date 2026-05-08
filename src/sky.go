package rkt

type Sky struct {
	Shader  Shader
	Texture Texture
	Buffer  *Buffer
}

func NewSky() *Sky {
	s := new(Sky)
	s.Shader = getShader("base/glsl/sky1")
	s.Buffer = NewBuffer(s.Shader, []BufferAttr{{BufferAttrPos, "a_Pos", 2}})
	s.Buffer.arrayFloat([]float32{
		-1, -1, +1, -1, -1, +1,
		-1, +1, +1, -1, +1, +1,
	})
	return s
}

func (s *Sky) Draw() {
	s.Shader.active()
	// s.Texture.bind()
	uPMatrix := s.Shader.getUniform("u_PMatrix")
	uVMatrix := s.Shader.getUniform("u_VMatrix")
	uSunDir := s.Shader.getUniform("u_SunDir")
	uSunColor := s.Shader.getUniform("u_SunColor")
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
	sun := ActiveLightEnv.DirLights[0]
	sun.Dir.uniform(uSunDir)
	sun.Color.uniform(uSunColor)
	s.Buffer.bind()
	s.Buffer.draw()
}
