package rkt

import "github.com/go-gl/gl/v2.1/gl"

type Vehicle struct {
	Name     string
	Parts    *PartNode
	Mass     float32
	Pos, Vel Vec3
	Rot, Ang Quat
	Stages   *StageNode
}

func NewVehicle(name string, root Part) *Vehicle {
	v := new(Vehicle)
	v.Name = name
	v.Parts = NewPartNode(root)
	v.Stages = &StageNode{nil, nil}
	return v
}

func (v *Vehicle) Draw() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	v.Pos.apply()
	v.Rot.apply()
	node := v.Parts
	for node != nil {
		node.Part.draw(&node.Offset)
		node = node.Lower
	}
	gl.PopMatrix()
}
func (v *Vehicle) Update(dt float32) {
	mass := float32(0.0)
	for node := v.Parts; node != nil; node = node.Lower {
		mass += node.Part.getMass()
	}

	v.Mass = mass
	for node := v.Parts; node != nil; node = node.Lower {
		node.Part.update(v, node, dt)
	}

	v.Vel.Z -= 9.0 * dt
	v.Pos = v.Pos.add(v.Vel.scale(dt))
	if v.Pos.Z < 0.0 {
		v.Pos.Z = 0.0
		if v.Vel.Z < 0.0 {
			v.Vel.X *= 0.8
			v.Vel.Y *= 0.8
			v.Vel.Z = 0.0
		}
	}
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
	p.Offset.Z = n.Offset.Z + nAttachPt.Z - pAttachPt.Z
	return p
}

type StageNode struct {
	Part Part
	Next *StageNode
}
