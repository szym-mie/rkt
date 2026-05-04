package rkt

import "log"

type Sky struct {
	Shader  Shader
	Texture Texture
	Buffer  *Buffer
}

func NewSky() *Sky {
	s := new(Sky)
	shader, ok := shaderMap["base/glsl/sky1"]
	if !ok {
		log.Fatalf("create: no such shader %s", "TODO")
	}

	texture, ok := textureMap["TODO"]
	if !ok {
		log.Fatalf("create: no such texture %s", "TODO")
	}

	s.Shader = shader
	s.Texture = texture
	s.Buffer = NewBuffer(shader, []BufferAttr{{BufferAttrPos, "a_Pos", 2}})
	bufferData := s.Buffer.newDataArray(6)
	bufferData[0] = -1
	bufferData[1] = -1
	bufferData[2] = -1
	bufferData[3] = +1
	bufferData[4] = +1
	bufferData[5] = -1
	bufferData[6] = +1
	bufferData[7] = +1
	bufferData[8] = +1
	bufferData[9] = -1
	bufferData[10] = -1
	bufferData[11] = +1
	s.Buffer.data(bufferData)
	return s
}

func (s *Sky) Draw() {
	s.Shader.active()
	s.Texture.bind()

	s.Buffer.draw()
}
