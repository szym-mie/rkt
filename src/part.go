package rkt

import (
	"log"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
)

type Part interface {
	draw(offset *Vec3)
	update(v *Vehicle, n *PartNode, dt float32)
	getMass() float32
	getAttachPts() (Vec3, Vec3)
	GetName() string
	Activate()
}

func NewPart(name string) Part {
	partDef, ok := partDefMap[name]
	if !ok {
		log.Fatalf("new_part: no partdef %s", name)
	}

	return partDef.create()
}

type PartBase struct {
	Def              *PartDef
	Geom             []Geom1
	IsActive, IsDead bool
}

func (p *PartBase) getAttachPts() (Vec3, Vec3) {
	return p.Def.Attach.Upper, p.Def.Attach.Lower
}
func (p *PartBase) GetName() string {
	return p.Def.Name
}
func (p *PartBase) drawModel(offset *Vec3) {
	offset.apply()
	for _, g := range p.Geom {
		g.draw()
	}
}
func (p *PartBase) drawAttachPts() {
	drawAttachPt(0.4, &p.Def.Attach.Lower)
	drawAttachPt(0.4, &p.Def.Attach.Upper)
}
func (p *PartBase) Activate() {
	p.IsActive = true
}

type PartCtrl struct {
	PartBase
}

func (p *PartCtrl) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	gl.PopMatrix()
}
func (p *PartCtrl) update(v *Vehicle, n *PartNode, dt float32) {
	// TODO: do nothing?
}
func (p *PartCtrl) getMass() float32 {
	return p.Def.Mass
}

type PartEngine struct {
	PartBase
	Plume    Plume
	FuelFlow float32
	FuelMass float32
}

func (p *PartEngine) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	if p.IsActive && p.FuelMass > 0.0 {
		p.Plume.draw()
	}

	gl.PopMatrix()
}
func (p *PartEngine) update(v *Vehicle, n *PartNode, dt float32) {
	e := p.Def.Engine
	if p.IsActive {
		p.FuelFlow = e.FuelDef.Flow
	} else if e.CanShutdown {
		p.FuelFlow = 0.0
	}

	if p.FuelMass > 0.0 {
		force := e.FuelDef.Impulse * p.FuelFlow * dt
		forceVec := Vec3{0.0, 0.0, force}
		log.Printf("%v", v.Rot.rotate(forceVec))
		v.Vel = v.Vel.add(v.Rot.rotate(forceVec).scale(1 / v.Mass))
		fuelCons := p.FuelFlow * dt
		p.FuelMass -= min(fuelCons, p.FuelMass)
		p.Plume.update(dt)
	} else {
		// engine shutdown
		p.FuelFlow = 0.0
		p.FuelMass = 0.0
	}
}
func (p *PartEngine) getMass() float32 {
	return p.Def.Mass + p.FuelMass
}

type PartDecoup struct {
	PartBase
	IsUsed bool
}

func (p *PartDecoup) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	gl.PopMatrix()
}
func (p *PartDecoup) update(v *Vehicle, n *PartNode, dt float32) {
	if p.IsActive && !p.IsUsed {
		un := n.Upper
		n.Upper = nil
		un.Lower = nil
		// TODO: create new vehicle, apply force
		p.IsUsed = true
	}
}
func (p *PartDecoup) getMass() float32 {
	return p.Def.Mass
}

type PartWing struct {
	PartBase
	stallRatio float32
}

func (p *PartWing) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	gl.PopMatrix()
}
func (p *PartWing) update(v *Vehicle, n *PartNode, dt time.Duration) {
	// TODO: no aerodynamics yet
}
func (p *PartWing) getMass() float32 {
	return p.Def.Mass
}
