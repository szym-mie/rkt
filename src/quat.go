package rkt

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type Quat struct {
	a, b, c, d float32
}

func NewAxisAngleQuat(theta float32, axis Vec3) Quat {
	rad := float64(theta) / 90 * math.Pi
	a := float32(math.Cos(rad))
	s := float32(math.Sin(rad))
	return Quat{a, axis.X * s, axis.Y * s, axis.Z * s}
}

func ZeroQuat() Quat {
	return Quat{1.0, 0.0, 0.0, 0.0}
}

func (q Quat) conj() Quat {
	return Quat{q.a, -q.b, -q.c, -q.d}
}

func (q Quat) norm() Quat {
	sum := float64(q.a*q.a + q.b*q.b + q.c*q.c + q.d*q.d)
	m := float32(math.Sqrt(sum))
	if m < 0.0001 {
		return Quat{1.0, 0.0, 0.0, 0.0}
	}

	return Quat{q.a / m, q.b / m, q.c / m, q.d / m}
}

func (q Quat) Product(p Quat) Quat {
	a := q.a*p.a - q.b*p.b - q.c*p.c - q.d*p.d
	b := q.a*p.b + q.b*p.a + q.c*p.d - q.d*p.c
	c := q.a*p.c - q.b*p.d + q.c*p.a + q.d*p.b
	d := q.a*p.d + q.b*p.c - q.c*p.b + q.d*p.a
	return Quat{a, b, c, d}
}

func (q Quat) rotate(v Vec3) Vec3 {
	p := Quat{0.0, v.X, v.Y, v.Z}
	o := q.Product(p).Product(q.conj())
	return Vec3{o.b, o.c, o.d}
}

func (q Quat) apply() {
	xw := 2 * q.b * q.a
	xx := 2 * q.b * q.b
	xy := 2 * q.b * q.c
	xz := 2 * q.b * q.d
	yw := 2 * q.c * q.a
	yy := 2 * q.c * q.c
	yz := 2 * q.c * q.d
	zw := 2 * q.d * q.a
	zz := 2 * q.d * q.d

	v := [16]float32{
		1 - yy - zz, xy + zw, xz - yw, 0.0,
		xy - zw, 1 - xx - zz, yz + xw, 0.0,
		xz + yw, yz - xw, 1 - xx - yy, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
	gl.MultMatrixf(&v[0])
}
