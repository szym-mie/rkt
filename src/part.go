package rkt

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type PartBase struct {
	Def              *PartDef
	Geom             []Geom
	IsActive, IsDead bool
}

func (p *PartBase) getAttachPts() (Vec3, Vec3) {
	return p.Def.Attach.Upper, p.Def.Attach.Lower
}

type Vehicle struct {
	Name     string
	PartTree *PartNode
	Mass     float32
	Pos, Vel Vec3
	Rot, Ang Quat
}

func NewVehicle(name string, partTree *PartNode) *Vehicle {
	v := new(Vehicle)
	v.Name = name
	v.PartTree = partTree
	return v
}

func (v *Vehicle) Draw() {
	v.Pos.apply()
	v.Rot.apply()
	node := v.PartTree
	for node != nil {
		node.Part.draw(&node.Offset)
		node = node.Lower
	}
}
func (v *Vehicle) Update() {
	mass := float32(0.0)
	node := v.PartTree
	for node != nil {
		mass += node.Part.getMass()
		node = node.Lower
	}

	v.Mass = mass
	node = v.PartTree
	for node != nil {
		node.Part.update(v, node)
		node = node.Lower
	}

	v.Pos = v.Pos.add(v.Vel)
}

type Part interface {
	draw(offset *Vec3)
	update(v *Vehicle, n *PartNode)
	getMass() float32
	getAttachPts() (Vec3, Vec3)
}

type PartNode struct {
	Part         Part
	Lower, Upper *PartNode
	Offset       Vec3
}

func NewPartNode(part Part) *PartNode {
	n := new(PartNode)
	n.Part = part
	return n
}

func (n *PartNode) AttachBelow(part Part) *PartNode {
	p := NewPartNode(part)
	// link both parts
	n.Lower = p
	p.Upper = n
	// calculate offset based on attachment points (height only for now)
	_, nAttachPt := n.Part.getAttachPts()
	pAttachPt, _ := part.getAttachPts()
	p.Offset.z = n.Offset.z + nAttachPt.z - pAttachPt.z
	return p
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
func (p *PartCtrl) update(v *Vehicle, n *PartNode) {
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
	if p.IsActive {
		p.Plume.draw()
	}

	gl.PopMatrix()
}
func (p *PartEngine) update(v *Vehicle, n *PartNode) {
	e := p.Def.Engine
	if p.IsActive {
		p.FuelFlow = e.FuelDef.Flow
	} else if e.CanShutdown {
		p.FuelFlow = 0.0
	}

	if p.FuelMass > 0.0 {
		force := e.FuelDef.Impulse * p.FuelFlow
		v.Vel.y += force / v.Mass
		p.FuelMass -= min(p.FuelFlow, p.FuelMass)
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
func (p *PartDecoup) update(v *Vehicle, n *PartNode) {
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
func (p *PartWing) update(v *Vehicle, n *PartNode) {
	// TODO: no aerodynamics yet
}
func (p *PartWing) getMass() float32 {
	return p.Def.Mass
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
