package rkt

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Mat3 [9]float32
type Mat4 [16]float32

const invEps = 0.000001

func NewMat3() *Mat3 {
	m := new(Mat3)
	m[0] = 1.0
	m[4] = 1.0
	m[8] = 1.0
	return m
}

func (m *Mat3) SetZero() {
	m[0] = 0.0
	m[1] = 0.0
	m[2] = 0.0
	m[3] = 0.0
	m[4] = 0.0
	m[5] = 0.0
	m[6] = 0.0
	m[7] = 0.0
	m[8] = 0.0
}
func (m *Mat3) SetIdentity() {
	m.SetZero()
	m[0] = 1.0
	m[4] = 1.0
	m[8] = 1.0
}
func (m *Mat3) SetScale1(k float32) {
	m[0] = k
	m[4] = k
	m[8] = k
}
func (m *Mat3) SetScale3(p Vec3) {
	m[0] = p.X
	m[4] = p.Y
	m[8] = p.Z
}
func (m *Mat3) Scale1(k float32) {
	n := new(Mat3)
	n.SetIdentity()
	n.SetScale1(k)
	m.MulSelf(n)
}
func (m *Mat3) Scale3(p Vec3) {
	n := new(Mat3)
	n.SetIdentity()
	n.SetScale3(p)
	m.MulSelf(n)
}
func (m *Mat3) Invert() *Mat3 {
	c0 := m[4]*m[8] - m[5]*m[7]
	c1 := m[5]*m[6] - m[3]*m[8]
	c2 := m[3]*m[7] - m[4]*m[6]
	det := m[0]*c0 + m[1]*c1 + m[2]*c2
	if det > -invEps && det < invEps {
		return m
	}
	det = 1.0 / det

	n := new(Mat3)
	n[0] = c0 * det
	n[1] = (m[2]*m[7] - m[1]*m[8]) * det
	n[2] = (m[1]*m[5] - m[2]*m[4]) * det
	n[3] = c1 * det
	n[4] = (m[0]*m[8] - m[2]*m[6]) * det
	n[5] = (m[2]*m[3] - m[0]*m[5]) * det
	n[6] = c2 * det
	n[7] = (m[3]*m[7] - m[4]*m[6]) * det
	n[8] = (m[1]*m[6] - m[0]*m[7]) * det
	return n
}
func (m *Mat3) Add(n *Mat3) *Mat3 {
	return &Mat3{
		m[0] + n[0], m[1] + n[1], m[2] + n[2],
		m[3] + n[3], m[4] + n[4], m[5] + n[5],
		m[6] + n[6], m[7] + n[7], m[8] + n[8],
	}
}
func (m *Mat3) Mul(n *Mat3) *Mat3 {
	r0 := Vec3{n[0], n[1], n[2]}
	r1 := Vec3{n[3], n[4], n[5]}
	r2 := Vec3{n[6], n[7], n[8]}
	c0 := Vec3{m[0], m[3], m[6]}
	c1 := Vec3{m[1], m[4], m[7]}
	c2 := Vec3{m[2], m[5], m[8]}
	return &Mat3{
		r0.Dot(c0), r0.Dot(c1), r0.Dot(c2),
		r1.Dot(c0), r1.Dot(c1), r1.Dot(c2),
		r2.Dot(c0), r2.Dot(c1), r2.Dot(c2),
	}
}
func (m *Mat3) MulSelf(n *Mat3) {
	r0 := Vec3{n[0], n[1], n[2]}
	r1 := Vec3{n[3], n[4], n[5]}
	r2 := Vec3{n[6], n[7], n[8]}
	c0 := Vec3{m[0], m[3], m[6]}
	c1 := Vec3{m[1], m[4], m[7]}
	c2 := Vec3{m[2], m[5], m[8]}
	m[0] = r0.Dot(c0)
	m[1] = r0.Dot(c1)
	m[2] = r0.Dot(c2)
	m[3] = r1.Dot(c0)
	m[4] = r1.Dot(c1)
	m[5] = r1.Dot(c2)
	m[6] = r2.Dot(c0)
	m[7] = r2.Dot(c1)
	m[8] = r2.Dot(c2)
}
func (m *Mat3) Apply(v Vec3) Vec3 {
	return Vec3{
		v.X*m[0] + v.Y*m[1] + v.Z*m[2],
		v.X*m[3] + v.Y*m[4] + v.Z*m[5],
		v.X*m[6] + v.Y*m[7] + v.Z*m[8],
	}
}
func (m *Mat3) SetRotX(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[4] = c
	m[5] = s
	m[7] = -s
	m[8] = c
}
func (m *Mat3) SetRotY(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[0] = c
	m[2] = -s
	m[6] = s
	m[8] = c
}
func (m *Mat3) SetRotZ(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[0] = c
	m[1] = s
	m[4] = -s
	m[5] = c
}
func (m *Mat3) RotX(theta float32) {
	n := new(Mat3)
	n.SetIdentity()
	n.SetRotX(theta)
	m.MulSelf(n)
}
func (m *Mat3) RotY(theta float32) {
	n := new(Mat3)
	n.SetIdentity()
	n.SetRotY(theta)
	m.MulSelf(n)
}
func (m *Mat3) RotZ(theta float32) {
	n := new(Mat3)
	n.SetIdentity()
	n.SetRotZ(theta)
	m.MulSelf(n)
}

func NewMat4() *Mat4 {
	m := new(Mat4)
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[15] = 1.0
	return m
}

func NewMat4Pos(v Vec3) *Mat4 {
	m := new(Mat4)
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[12] = v.X
	m[13] = v.Y
	m[14] = v.Z
	m[15] = 1.0
	return m
}

func (m *Mat4) SetZero() {
	m[0] = 0.0
	m[1] = 0.0
	m[2] = 0.0
	m[3] = 0.0
	m[4] = 0.0
	m[5] = 0.0
	m[6] = 0.0
	m[7] = 0.0
	m[8] = 0.0
	m[9] = 0.0
	m[10] = 0.0
	m[11] = 0.0
	m[12] = 0.0
	m[13] = 0.0
	m[14] = 0.0
	m[15] = 0.0
}
func (m *Mat4) SetIdentity() {
	m.SetZero()
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[15] = 1.0
}
func (m *Mat4) SetPos(p Vec3) {
	m[12] = p.X
	m[13] = p.Y
	m[14] = p.Z
}
func (m *Mat4) AddPosSelf(p Vec3) {
	m[12] += p.X
	m[13] += p.Y
	m[14] += p.Z
}
func (m *Mat4) From(n *Mat4) {
	m[0] = n[0]
	m[1] = n[1]
	m[2] = n[2]
	m[3] = n[3]
	m[4] = n[4]
	m[5] = n[5]
	m[6] = n[6]
	m[7] = n[7]
	m[8] = n[8]
	m[9] = n[9]
	m[10] = n[10]
	m[11] = n[11]
	m[12] = n[12]
	m[13] = n[13]
	m[14] = n[14]
	m[15] = n[15]
}
func (m *Mat4) SetScale1(k float32) {
	m[0] = k
	m[5] = k
	m[10] = k
}
func (m *Mat4) SetScale3(p Vec3) {
	m[0] = p.X
	m[5] = p.Y
	m[10] = p.Z
}
func (m *Mat4) Scale1(k float32) {
	n := new(Mat4)
	n.SetIdentity()
	n.SetScale1(k)
	m.MulSelf(n)
}
func (m *Mat4) Scale3(p Vec3) {
	n := new(Mat4)
	n.SetIdentity()
	n.SetScale3(p)
	m.MulSelf(n)
}

func (m *Mat4) Invert() *Mat4 {
	b00 := m[0]*m[5] - m[1]*m[4]
	b01 := m[0]*m[6] - m[2]*m[4]
	b02 := m[0]*m[7] - m[3]*m[4]
	b03 := m[1]*m[6] - m[2]*m[5]
	b04 := m[1]*m[7] - m[3]*m[5]
	b05 := m[2]*m[7] - m[3]*m[6]
	b06 := m[8]*m[13] - m[9]*m[12]
	b07 := m[8]*m[14] - m[10]*m[12]
	b08 := m[8]*m[15] - m[11]*m[12]
	b09 := m[9]*m[14] - m[10]*m[13]
	b10 := m[9]*m[15] - m[11]*m[13]
	b11 := m[10]*m[15] - m[11]*m[14]

	det := b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
	if det > -invEps && det < invEps {
		return m
	}
	det = 1.0 / det

	n := new(Mat4)
	n[0] = (m[5]*b11 - m[6]*b10 + m[7]*b09) * det
	n[1] = (m[2]*b10 - m[1]*b11 - m[3]*b09) * det
	n[2] = (m[13]*b05 - m[14]*b04 + m[15]*b03) * det
	n[3] = (m[10]*b04 - m[9]*b05 - m[11]*b03) * det
	n[4] = (m[6]*b08 - m[4]*b11 - m[7]*b07) * det
	n[5] = (m[0]*b11 - m[2]*b08 + m[3]*b07) * det
	n[6] = (m[14]*b02 - m[12]*b05 - m[15]*b01) * det
	n[7] = (m[8]*b05 - m[10]*b02 + m[11]*b01) * det
	n[8] = (m[4]*b10 - m[5]*b08 + m[7]*b06) * det
	n[9] = (m[1]*b08 - m[0]*b10 - m[3]*b06) * det
	n[10] = (m[12]*b04 - m[13]*b02 + m[15]*b00) * det
	n[11] = (m[9]*b02 - m[8]*b04 - m[11]*b00) * det
	n[12] = (m[5]*b07 - m[4]*b09 - m[6]*b06) * det
	n[13] = (m[0]*b09 - m[1]*b07 + m[2]*b06) * det
	n[14] = (m[13]*b01 - m[12]*b03 - m[14]*b00) * det
	n[15] = (m[8]*b03 - m[9]*b01 + m[10]*b00) * det
	return n
}
func (m *Mat4) Add(n *Mat4) *Mat4 {
	return &Mat4{
		m[0] + n[0], m[1] + n[1], m[2] + n[2], m[3] + n[3],
		m[4] + n[4], m[5] + n[5], m[6] + n[6], m[7] + n[7],
		m[8] + n[8], m[9] + n[9], m[10] + n[10], m[11] + n[11],
		m[12] + n[12], m[13] + n[13], m[14] + n[14], m[15] + n[15],
	}
}
func (m *Mat4) Mul(n *Mat4) *Mat4 {
	r0 := Vec4{n[0], n[1], n[2], n[3]}
	r1 := Vec4{n[4], n[5], n[6], n[7]}
	r2 := Vec4{n[8], n[9], n[10], n[11]}
	r3 := Vec4{n[12], n[13], n[14], n[15]}
	c0 := Vec4{m[0], m[4], m[8], m[12]}
	c1 := Vec4{m[1], m[5], m[9], m[13]}
	c2 := Vec4{m[2], m[6], m[10], m[14]}
	c3 := Vec4{m[3], m[7], m[11], m[15]}
	return &Mat4{
		r0.Dot(c0), r0.Dot(c1), r0.Dot(c2), r0.Dot(c3),
		r1.Dot(c0), r1.Dot(c1), r1.Dot(c2), r1.Dot(c3),
		r2.Dot(c0), r2.Dot(c1), r2.Dot(c2), r2.Dot(c3),
		r3.Dot(c0), r3.Dot(c1), r3.Dot(c2), r3.Dot(c3),
	}
}
func (m *Mat4) MulSelf(n *Mat4) {
	r0 := Vec4{n[0], n[1], n[2], n[3]}
	r1 := Vec4{n[4], n[5], n[6], n[7]}
	r2 := Vec4{n[8], n[9], n[10], n[11]}
	r3 := Vec4{n[12], n[13], n[14], n[15]}
	c0 := Vec4{m[0], m[4], m[8], m[12]}
	c1 := Vec4{m[1], m[5], m[9], m[13]}
	c2 := Vec4{m[2], m[6], m[10], m[14]}
	c3 := Vec4{m[3], m[7], m[11], m[15]}
	m[0] = r0.Dot(c0)
	m[1] = r0.Dot(c1)
	m[2] = r0.Dot(c2)
	m[3] = r0.Dot(c3)
	m[4] = r1.Dot(c0)
	m[5] = r1.Dot(c1)
	m[6] = r1.Dot(c2)
	m[7] = r1.Dot(c3)
	m[8] = r2.Dot(c0)
	m[9] = r2.Dot(c1)
	m[10] = r2.Dot(c2)
	m[11] = r2.Dot(c3)
	m[12] = r3.Dot(c0)
	m[13] = r3.Dot(c1)
	m[14] = r3.Dot(c2)
	m[15] = r3.Dot(c3)
}
func (m *Mat4) SetRotX(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[5] = c
	m[6] = s
	m[9] = -s
	m[10] = c
}
func (m *Mat4) SetRotY(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[0] = c
	m[2] = -s
	m[8] = s
	m[10] = c
}
func (m *Mat4) SetRotZ(theta float32) {
	s := Sinf(theta)
	c := Cosf(theta)
	m[0] = c
	m[1] = s
	m[4] = -s
	m[5] = c
}
func (m *Mat4) RotX(theta float32) {
	n := new(Mat4)
	n.SetIdentity()
	n.SetRotX(theta)
	m.MulSelf(n)
}
func (m *Mat4) RotY(theta float32) {
	n := new(Mat4)
	n.SetIdentity()
	n.SetRotY(theta)
	m.MulSelf(n)
}
func (m *Mat4) RotZ(theta float32) {
	n := new(Mat4)
	n.SetIdentity()
	n.SetRotZ(theta)
	m.MulSelf(n)
}
func (m *Mat4) Frustum(aspect, depthNear, depthFar float32) {
	f := aspect - 1.0
	depthDiff := depthNear - depthFar
	m.SetZero()
	m[0] = depthNear / (1.0 + f)
	m[5] = depthNear
	m[10] = (depthNear + depthFar) / depthDiff
	m[11] = -1.0
	m[14] = 2.0 * depthNear * depthFar / depthDiff
}
func (m *Mat4) Ortho(aspect, depthNear, depthFar float32) {
	f := aspect - 1.0
	depthDiff := depthNear - depthFar
	m.SetZero()
	m[0] = 1.0 / (1.0 + f)
	m[5] = 1.0
	m[10] = -2.0 / depthDiff
	m[14] = (depthNear + depthFar) / depthDiff
	m[15] = 1.0
}
func (m *Mat4) uniform(location int32) {
	gl.UniformMatrix4fv(location, 1, false, &m[0])
}
