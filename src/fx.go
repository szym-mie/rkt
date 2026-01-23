package rkt

type Plume struct {
	geom      *Geom1
	offset    Vec3
	size      float32
	texOffset float32
}

func (p *Plume) draw() {
	p.offset.apply()
	p.geom.drawTexOffset(Vec2{p.texOffset, 0.0})
}
func (p *Plume) update(dt float32) {
	// TODO: add sparks
	p.texOffset += 0.831
	if p.texOffset > 9.0 {
		p.texOffset -= 8.0
	}
}
