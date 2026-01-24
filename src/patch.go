package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type PatchObjDef struct {
	ObjName string `json:"obj"`
	Pos     Vec3   `json:"pos"`
}

func (d *PatchObjDef) create() *PatchObj {
	p := new(PatchObj)
	geom, ok := geom1DefMap[d.ObjName]
	if !ok {
		log.Fatalf("create: no such geom1def %s", d.ObjName)
	}

	// TODO: use cloning in parts instead of create, here just reference the
	// same geometry object
	p.geom = geom.create()
	p.pos = d.Pos
	return p
}

type PatchObj struct {
	geom *Geom1
	pos  Vec3
}

type PatchDef struct {
	Lat      float32       `json:"lat"`
	Lon      float32       `json:"lon"`
	Scale    float32       `json:"scale"`
	GeomName string        `json:"geom"`
	Objs     []PatchObjDef `json:"objs"`
}

func (d *PatchDef) create() *Patch {
	p := new(Patch)
	geom, ok := geom2DefMap[d.GeomName]
	if !ok {
		log.Fatalf("create: no such geom2def %s", d.GeomName)
	}

	p.Scale = d.Scale
	p.geom = geom.create()
	p.objs = make([]PatchObj, len(d.Objs))
	for i, def := range d.Objs {
		p.objs[i] = *def.create()
	}

	return p
}

type Patch struct {
	Pos   Vec3
	Scale float32
	geom  *Geom2
	objs  []PatchObj
}

func NewPatch(patchName string) *Patch {
	patchDef, ok := patchDefMap[patchName]
	if !ok {
		log.Fatalf("new_patch: no such patchdef %s", patchName)
	}

	return patchDef.create()
}

func (p *PatchObj) draw() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.pos.apply()
	p.geom.draw()
	gl.PopMatrix()
}

func (p *Patch) Draw() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.Pos.apply()
	for _, obj := range p.objs {
		obj.draw()
	}

	gl.Scalef(p.Scale, p.Scale, p.Scale)
	p.geom.draw()
	gl.PopMatrix()
}
