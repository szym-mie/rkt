package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type Patch struct {
	Pos   Vec3
	Scale float32
	geom  *Geom
}

func NewPatch(geomName string) *Patch {
	p := new(Patch)
	geomDef, ok := geomDefMap[geomName]
	if !ok {
		log.Fatalf("new_patch: no such geomdef %s", geomName)
	}

	p.geom = geomDef.create()
	p.Scale = 1.0
	return p
}

func (p *Patch) Draw() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()

	p.Pos.apply()
	gl.Scalef(p.Scale, p.Scale, p.Scale)
	p.geom.draw()
	gl.PopMatrix()
}
