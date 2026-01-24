package rkt

import (
	"log"
)

type PlumeDef struct {
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
}

func (d *PlumeDef) create() *Plume {
	p := new(Plume)
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s", d.TextureName)
	}

	g := NewGeom1(texture, 12)
	ringVts := buildRingVertices(6, 1.4, -2.0)
	for i := range 6 {
		j := (i + 1) % 6
		ri := i * 6
		g.Vertices[ri+0] = Vec3{0.0, 0.0, 0.0}
		g.Vertices[ri+1] = ringVts[i]
		g.Vertices[ri+2] = ringVts[j]
		g.Vertices[ri+3] = Vec3{0.0, 0.0, -6.0}
		g.Vertices[ri+4] = ringVts[j]
		g.Vertices[ri+5] = ringVts[i]

		g.TexCoords[ri+0] = Vec2{0.5, 0.0}
		g.TexCoords[ri+1] = Vec2{0.0, 0.5}
		g.TexCoords[ri+2] = Vec2{1.0, 0.5}
		g.TexCoords[ri+3] = Vec2{0.5, 1.0}
		g.TexCoords[ri+4] = Vec2{1.0, 0.5}
		g.TexCoords[ri+5] = Vec2{0.0, 0.5}
	}

	p.geom = g
	p.offset = d.Offset
	p.size = d.Size
	return p
}

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
