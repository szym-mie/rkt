package rkt

type PlumeDef struct {
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
}

const fxShaderName = "base/glsl/nolight1"

func (d *PlumeDef) create() *Plume {
	p := new(Plume)

	// TODO: revamp fx with particles
	// p.geom = g
	p.offset = d.Offset
	p.size = d.Size
	p.local = Matrix4{}
	p.local.SetPos(p.offset)
	p.local.SetScale1(p.size)
	return p
}

type Plume struct {
	geom      *Geom1
	offset    Vec3
	local     Matrix4
	size      float32
	texOffset float32
}

func (p *Plume) draw(model Matrix4) {
	// p.geom.drawTexOffset(Vec2{p.texOffset, 0.0})
}
func (p *Plume) update(dt float32) {
	// TODO: add sparks
	p.texOffset += 0.831
	if p.texOffset > 9.0 {
		p.texOffset -= 8.0
	}
}
