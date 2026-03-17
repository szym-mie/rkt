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
	getInertiaCoeff() Vec3
	getDrag(aoa float32) float32
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
	offset.Apply()
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
func (p *PartBase) getMass() float32 {
	return p.Def.Mass
}
func (p *PartBase) getInertiaCoeff() Vec3 {
	return p.Def.Body.InertiaCoeff
}
func (p *PartBase) getDrag(aoa float32) float32 {
	return 1.0 * 0.5
}

type PartHull struct {
	PartBase
}

func (p *PartHull) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	gl.PopMatrix()
}
func (p *PartHull) update(v *Vehicle, n *PartNode, dt float32) {
	// do nothing
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
		nv := v.Fork(n)
		fv := nv.Rot.Rotate(Vec3{0, 0, -10})
		nv.Vel = nv.Vel.Add(fv)
		nv.Link()
		p.IsUsed = true
	}
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
		v.Vel = v.Vel.Add(v.Rot.Rotate(forceVec).MulSca(1 / v.Mass))
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

type PartChute struct {
	PartBase
	IsUsed       bool
	IsCut        bool
	ExtraDrag    float32
	ChuteGeom    *Geom1
	ChuteAntiRot Quat
	ChuteRot     Quat
	DeployTime   float32
	Height       float32
	Radius       float32
}

func (p *PartChute) draw(offset *Vec3) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	p.drawModel(offset)
	r := p.Radius
	h := p.Height
	p.ChuteAntiRot.Apply()
	p.ChuteRot.Conj().Apply()
	gl.Scalef(r, r, h)
	if p.IsActive && !p.IsCut {
		p.ChuteGeom.draw()
	}

	gl.PopMatrix()
}
func (p *PartChute) update(v *Vehicle, n *PartNode, dt float32) {
	c := p.Def.Chute
	if p.IsActive && !p.IsCut {
		p.ChuteAntiRot = v.Rot.Conj()
		p.ChuteRot = NewVecDiffQuat(v.Vel, Vec3{0, 0, -1}).Norm()
		if p.IsUsed {
			p.DeployTime += dt
		} else {
			p.DeployTime = 0
			p.IsUsed = true
		}
	}

	if p.IsUsed && v.Vel.LenSq() < 0.5 {
		p.IsCut = true
		p.DeployTime = 0.0
	}

	deployFrac := p.DeployTime / c.DeployTime
	if deployFrac > 1.0 {
		deployFrac = 1.0
	}
	deployFracSq := deployFrac * deployFrac

	p.Height = deployFrac * c.Height
	p.Radius = deployFracSq * c.Height
	dragMag := (0.5 + deployFracSq*c.Area*c.Drag) * v.Vel.LenSq() * 0.5
	// linear impulse
	dragForce := v.Vel.Norm().MulSca(-dragMag * dt)

	// EXPERIMENTAL: impose rotation (simpified cylinder)
	// radius from centre of mass (only works if part above root)
	radius := Vec3{0, 0, 1}
	// angular impulse
	dragTorque := radius.Cross(v.Rot.Conj().Rotate(dragForce))
	if v.Vel.Len() > 1.0 {
		// log.Println(v.Vel.Len())
		log.Println( /*v.Ang,*/ v.Rot.Conj().Rotate(dragForce), v.Inertia) //, dragTorque, v.Inertia)
		v.Ang = v.Ang.Add(v.Rot.Rotate(dragTorque).Div(v.Inertia))
		v.Vel = v.Vel.Add(dragForce.MulSca(1 / v.Mass))
	}
}
func (p *PartChute) getMass() float32 {
	return p.Def.Mass + p.Def.Chute.Mass
}
func (p *PartChute) getDrag(aoa float32) float32 {
	drag := p.PartBase.getDrag(aoa)
	if p.IsActive {
		drag += p.ExtraDrag
	}

	return drag
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
