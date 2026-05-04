package rkt

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Matrix4 [16]float32

func NewMatrix4() *Matrix4 {
	m := &Matrix4{}
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[15] = 1.0
	return m
}

func NewMatrix4Pos(v Vec3) *Matrix4 {
	m := &Matrix4{}
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[12] = v.X
	m[13] = v.Y
	m[14] = v.Z
	m[15] = 1.0
	return m
}

func (m *Matrix4) SetZero() {
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
func (m *Matrix4) SetIdentity() {
	m.SetZero()
	m[0] = 1.0
	m[5] = 1.0
	m[10] = 1.0
	m[15] = 1.0
}
func (m *Matrix4) SetPos(p Vec3) {
	m[12] = p.X
	m[13] = p.Y
	m[14] = p.Z
}
func (m *Matrix4) AddPosSelf(p Vec3) {
	m[12] += p.X
	m[13] += p.Y
	m[14] += p.Z
}
func (m *Matrix4) AddPos(p Vec3) *Matrix4 {
	n := new(Matrix4)
	n[0] = m[0]
	n[1] = m[1]
	n[2] = m[2]
	n[3] = m[3]
	n[4] = m[4]
	n[5] = m[5]
	n[6] = m[6]
	n[7] = m[7]
	n[8] = m[8]
	n[9] = m[9]
	n[10] = m[10]
	n[11] = m[11]
	n[12] = m[12] + p.X
	n[13] = m[13] + p.Y
	n[14] = m[14] + p.Z
	n[15] = m[15]
	return n
}
func (m *Matrix4) SetScale1(k float32) {
	m[0] = k
	m[5] = k
	m[10] = k
}
func (m *Matrix4) SetScale3(p Vec3) {
	m[0] = p.X
	m[5] = p.Y
	m[10] = p.Z
}
func (m *Matrix4) Scale1(k float32) {
	n := new(Matrix4)
	n.SetIdentity()
	n.SetScale1(k)
	m.MulSelf(n)
}
func (m *Matrix4) Scale3(p Vec3) {
	n := new(Matrix4)
	n.SetIdentity()
	n.SetScale3(p)
	m.MulSelf(n)
}
func (m *Matrix4) Add(n *Matrix4) *Matrix4 {
	return &Matrix4{
		m[0] + n[0], m[1] + n[1], m[2] + n[2], m[3] + n[3],
		m[4] + n[4], m[5] + n[5], m[6] + n[6], m[7] + n[7],
		m[8] + n[8], m[9] + n[9], m[10] + n[10], m[11] + n[11],
		m[12] + n[12], m[13] + n[13], m[14] + n[14], m[15] + n[15],
	}
}
func (m *Matrix4) Mul(n *Matrix4) *Matrix4 {
	r0 := Vec4{n[0], n[1], n[2], n[3]}
	r1 := Vec4{n[4], n[5], n[6], n[7]}
	r2 := Vec4{n[8], n[9], n[10], n[11]}
	r3 := Vec4{n[12], n[13], n[14], n[15]}
	c0 := Vec4{m[0], m[4], m[8], m[12]}
	c1 := Vec4{m[1], m[5], m[9], m[13]}
	c2 := Vec4{m[2], m[6], m[10], m[14]}
	c3 := Vec4{m[3], m[7], m[11], m[15]}

	return &Matrix4{
		r0.Dot(c0), r0.Dot(c1), r0.Dot(c2), r0.Dot(c3),
		r1.Dot(c0), r1.Dot(c1), r1.Dot(c2), r1.Dot(c3),
		r2.Dot(c0), r2.Dot(c1), r2.Dot(c2), r2.Dot(c3),
		r3.Dot(c0), r3.Dot(c1), r3.Dot(c2), r3.Dot(c3),
	}
}
func (m *Matrix4) MulSelf(n *Matrix4) {
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
func (m *Matrix4) SetRotX(theta float32) {
	s := float32(math.Sin(float64(theta)))
	c := float32(math.Cos(float64(theta)))

	m[5] = c
	m[6] = s
	m[9] = -s
	m[10] = c
}
func (m *Matrix4) SetRotY(theta float32) {
	s := float32(math.Sin(float64(theta)))
	c := float32(math.Cos(float64(theta)))

	m[0] = c
	m[2] = -s
	m[8] = s
	m[10] = c
}
func (m *Matrix4) SetRotZ(theta float32) {
	s := float32(math.Sin(float64(theta)))
	c := float32(math.Cos(float64(theta)))

	m[0] = c
	m[1] = s
	m[4] = -s
	m[5] = c
}
func (m *Matrix4) RotX(theta float32) {
	n := new(Matrix4)
	n.SetIdentity()
	n.SetRotX(theta)
	m.MulSelf(n)
}
func (m *Matrix4) RotY(theta float32) {
	n := new(Matrix4)
	n.SetIdentity()
	n.SetRotY(theta)
	m.MulSelf(n)
}
func (m *Matrix4) RotZ(theta float32) {
	n := new(Matrix4)
	n.SetIdentity()
	n.SetRotZ(theta)
	m.MulSelf(n)
}
func (m *Matrix4) Frustum(aspect, depthNear, depthFar float32) {
	f := aspect - 1.0
	depthDiff := depthNear - depthFar
	m.SetZero()
	m[0] = depthNear / (1.0 + f)
	m[5] = depthNear
	m[10] = (depthNear + depthFar) / depthDiff
	m[11] = -1.0
	m[14] = 2.0 * depthNear * depthFar / depthDiff
}
func (m *Matrix4) Ortho(aspect, depthNear, depthFar float32) {
	f := aspect - 1.0
	depthDiff := depthNear - depthFar
	m.SetZero()
	m[0] = 1.0 / (1.0 + f)
	m[5] = 1.0
	m[10] = -2.0 / depthDiff
	m[14] = (depthNear + depthFar) / depthDiff
	m[15] = 1.0
}
func (m *Matrix4) uniform(location int32) {
	gl.UniformMatrix4fv(location, 1, false, &m[0])
}
