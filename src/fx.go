package rkt

import "log"

type PlumeDef struct {
	ShaderName  string  `json:"shader"`
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
	PtCnt       uint32  `json:"pt_cnt"`
	EmitRate    uint32  `json:"emit_rate"`
	InitVel     float32 `json:"init_vel"`
	DragCoeff   float32 `json:"drag"`
}

var plumeBufferAttrs = []BufferAttr{
	{BufferAttrPos, "a_Pos", 3},
}

func (d *PlumeDef) create() *Plume {
	p := new(Plume)
	// TODO: revamp fx with particles
	p.shader = getShader(d.ShaderName)
	p.buffer = NewBuffer(p.shader, plumeBufferAttrs)
	p.offset = d.Offset
	p.local = Matrix4{}
	p.local.SetPos(p.offset)
	p.ptPos = make([]Vec3, d.PtCnt)
	p.ptVel = make([]Vec3, d.PtCnt)
	p.ptCnt = d.PtCnt
	p.emitRate = d.EmitRate
	p.emitIndex = 0
	return p
}

type Plume struct {
	shader    Shader
	buffer    *Buffer
	offset    Vec3
	local     Matrix4
	ptPos     []Vec3
	ptVel     []Vec3
	ptCnt     uint32
	emitRate  uint32
	emitIndex uint32
}

func (p *Plume) draw(model Matrix4) {
	p.shader.active()
	uPMatrix := p.shader.getUniform("u_PMatrix")
	uVMatrix := p.shader.getUniform("u_VMatrix")
	ActivePV.ProjMatrix.uniform(uPMatrix)
	ActivePV.ViewMatrix.uniform(uVMatrix)
	p.buffer.bind()
	p.buffer.drawPts()
}
func (p *Plume) update(pos Vec3, dt float32) {
	// TODO: add sparks
	log.Printf("%v\n", pos)
	p.ptPos[p.emitIndex] = pos
	p.emitIndex++
	p.emitIndex %= p.ptCnt
	p.buffer.arrayVec3(p.ptPos)
}
