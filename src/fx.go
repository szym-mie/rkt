package rkt

type PlumeDef struct {
	ShaderName  string  `json:"shader"`
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
	PtCnt       uint32  `json:"pt_cnt"`
	EmitRate    uint32  `json:"emit_rate"`
	InitVel     Vec3    `json:"init_vel"`
	SideVel     Vec3    `json:"side_vel"`
	DragCoeff   float32 `json:"drag"`
}

var plumeBufferAttrs = []BufferAttr{
	{BufferAttrPos, "a_Pos", 4},
}

func (d *PlumeDef) create() *Plume {
	p := new(Plume)
	// TODO: revamp fx with particles
	p.shader = getShader(d.ShaderName)
	p.buffer = NewBuffer(p.shader, plumeBufferAttrs)
	p.ptPos = make([]Vec4, d.PtCnt)
	p.ptVel = make([]Vec3, d.PtCnt)
	p.ptDecay = float32(d.EmitRate) / float32(d.PtCnt)
	p.def = d
	return p
}

type Plume struct {
	shader    Shader
	buffer    *Buffer
	ptPos     []Vec4
	ptVel     []Vec3
	ptDecay   float32
	emitIndex uint32
	def       *PlumeDef
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
func (p *Plume) emit(pos, vel, offset Vec3, rot Quat) {
	for range p.def.EmitRate {
		totalOffset := offset.Add(p.def.Offset)
		worldPos := pos.Add(rot.Rotate(totalOffset))
		sideVel := p.def.SideVel.RandomSphere()
		localVel := p.def.InitVel.Add(sideVel)
		worldVel := vel.Add(rot.Rotate(localVel))
		p.ptPos[p.emitIndex] = worldPos.To4(0.0)
		p.ptVel[p.emitIndex] = worldVel
		p.emitIndex++
		p.emitIndex %= p.def.PtCnt
	}
}
func (p *Plume) update(dt float32) {
	// TODO: add sparks
	drag := 1.0 - p.def.DragCoeff*dt
	for i := range p.ptPos {
		vel := p.ptVel[i].MulSca(dt)
		p.ptPos[i] = p.ptPos[i].Add(vel.To4(p.ptDecay))
		p.ptVel[i] = p.ptVel[i].MulSca(drag)
	}
	p.buffer.arrayVec4(p.ptPos)
}
