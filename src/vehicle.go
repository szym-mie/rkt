package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type Vehicle struct {
	Name          string
	Parts         *PartNode
	Mass          float32
	Com           Vec3
	Inertia       Vec3
	Height        float32
	Pos, Vel, Ang Vec3
	Rot           Quat
	Stages        *StageNode
}

func NewVehicle(name string, root Part) *Vehicle {
	v := new(Vehicle)
	v.Name = name
	v.Parts = NewPartNode(root)
	v.Stages = &StageNode{nil, nil}
	v.Rot = ZeroQuat()
	return v
}

func (v *Vehicle) Fork(nodes *PartNode) *Vehicle {
	w := new(Vehicle)
	w.Name = "_debris"
	w.Parts = nodes
	w.Stages = &StageNode{}
	offset := nodes.Offset
	for node := nodes; node != nil; node = node.Lower {
		node.Offset = node.Offset.Sub(offset)
	}

	w.Pos = v.Pos.Add(v.Rot.Rotate(offset))
	w.Vel = v.Vel
	w.Rot = v.Rot
	w.Ang = v.Ang
	return w
}
func (v *Vehicle) Draw() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	v.Pos.Apply()
	v.Rot.Apply()
	for node := v.Parts; node != nil; node = node.Lower {
		node.Part.draw(&node.Offset)
	}
	for node := v.Parts.Upper; node != nil; node = node.Upper {
		node.Part.draw(&node.Offset)
	}
	DrawDiamond(v.Com)
	gl.PopMatrix()
}
func (v *Vehicle) Update(dt float32) {
	v.UpdateHeight()
	v.UpdateCoM()
	v.UpdateInertia()

	for node := v.Parts; node != nil; node = node.Lower {
		node.Part.update(v, node, dt)
	}
	for node := v.Parts.Upper; node != nil; node = node.Upper {
		node.Part.update(v, node, dt)
	}

	v.Vel.Z -= 9.0 * dt
	v.Pos = v.Pos.Add(v.Vel.MulSca(dt))
	if v.Pos.Z < v.Height {
		v.Pos.Z = v.Height
		v.Rot = ZeroQuat()
		v.Ang.X = 0.0
		v.Ang.Y = 0.0
		v.Ang.Z = 0.0
		if v.Vel.Z < 0.0 {
			v.Vel.X *= 0.8
			v.Vel.Y *= 0.8
			v.Vel.Z = 0.0
		}
	}
	w := Quat{0.0, v.Ang.X, v.Ang.Y, v.Ang.Z}.Scale(dt * 0.5)
	w.a += 1.0
	v.Rot = w.Product(v.Rot).Norm()
	v.Ang = v.Ang.MulSca(0.99)
}
func (v *Vehicle) AddToStage(part Part) {
	s := new(StageNode)
	s.Part = part
	s.Next = v.Stages
	v.Stages = s
}
func (v *Vehicle) AddStage() {
	v.AddToStage(nil)
}
func (v *Vehicle) AttachAbove(node *PartNode, part Part) *PartNode {
	v.AddToStage(part)
	return node.AttachAbove(part)
}
func (v *Vehicle) AttachBelow(node *PartNode, part Part) *PartNode {
	v.AddToStage(part)
	return node.AttachBelow(part)
}
func (v *Vehicle) ApplyStage() {
	for s := v.Stages; s != nil; s = s.Next {
		if s.Part != nil {
			s.Part.Activate()
		} else {
			v.Stages = s.Next
			break
		}
	}
}
func (v *Vehicle) UpdateHeight() {
	last := v.Parts
	for node := v.Parts; node != nil; node = node.Lower {
		last = node
	}
	_, lAttachPt := last.Part.getAttachPts()
	v.Height = -last.Offset.Z - lAttachPt.Z
}
func (v *Vehicle) UpdateCoM() {
	com := Vec3{}
	totalMass := v.Parts.Part.getMass()
	for node := v.Parts.Lower; node != nil; node = node.Lower {
		nodeMass := node.Part.getMass()
		totalMass += nodeMass
		com = com.Lerp(node.Offset, nodeMass/totalMass)
	}
	for node := v.Parts.Upper; node != nil; node = node.Upper {
		nodeMass := node.Part.getMass()
		totalMass += nodeMass
		com = com.Lerp(node.Offset, nodeMass/totalMass)
	}

	v.Mass = totalMass
	v.Com = com
	log.Println(v.Name, com)
}
func (v *Vehicle) UpdateInertia() {
	totalInertia := Vec3{}
	for node := v.Parts.Lower; node != nil; node = node.Lower {
		radiusSq := node.Offset.Sub(v.Com).AxisLenSq()
		mass := node.Part.getMass()
		offset := radiusSq.MulSca(mass)
		inertia := node.Part.getInertiaCoeff().MulSca(mass).Add(offset)
		totalInertia = totalInertia.Add(inertia)
	}
	for node := v.Parts; node != nil; node = node.Upper {
		radiusSq := node.Offset.Sub(v.Com).AxisLenSq()
		mass := node.Part.getMass()
		offset := radiusSq.MulSca(mass)
		inertia := node.Part.getInertiaCoeff().MulSca(mass).Add(offset)
		totalInertia = totalInertia.Add(inertia)
	}

	v.Inertia = totalInertia
	log.Println(v.Name, totalInertia)
}
func (v *Vehicle) ApplyImpulse(imp, pt Vec3) {
	// angular impulse
	torque := pt.Sub(v.Com).Cross(v.Rot.Conj().Rotate(imp))
	// TODO: only remaining energy should go into linear motion not full
	v.Ang = v.Ang.Add(v.Rot.Rotate(torque).Div(v.Inertia))
	v.Vel = v.Vel.Add(imp.MulSca(1.0 / v.Mass))
	// TODO: draw force vector because why not
	DrawVector(imp.MulSca(0.01), v.Pos.Add(v.Rot.Rotate(pt)))
}

// BEGIN STUPID
var Vehicles []*Vehicle = make([]*Vehicle, 128)
var VehiclesIndex uint = 0

func (v *Vehicle) Link() {
	if VehiclesIndex < 128 {
		Vehicles[VehiclesIndex] = v
		VehiclesIndex++
	}
}

// END STUPID

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

func (n *PartNode) AttachAbove(part Part) *PartNode {
	p := NewPartNode(part)
	// link both parts
	n.Upper = p
	p.Lower = n
	// calculate offset based on attachment points (height only for now)
	nAttachPt, _ := n.Part.getAttachPts()
	_, pAttachPt := part.getAttachPts()
	p.Offset.Z = n.Offset.Z + nAttachPt.Z - pAttachPt.Z
	return p
}
func (n *PartNode) AttachBelow(part Part) *PartNode {
	p := NewPartNode(part)
	// link both parts
	n.Lower = p
	p.Upper = n
	// calculate offset based on attachment points (height only for now)
	_, nAttachPt := n.Part.getAttachPts()
	pAttachPt, _ := part.getAttachPts()
	p.Offset.Z = n.Offset.Z + nAttachPt.Z - pAttachPt.Z
	return p
}

type StageNode struct {
	Part Part
	Next *StageNode
}
