package rkt

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type Quat struct {
	a, b, c, d float32
}

func NewAxisAngleQuat(theta float32, axis Vec3) Quat {
	rad := float64(theta) / 180 * math.Pi
	a := float32(math.Cos(rad * 0.5))
	s := float32(math.Sin(rad * 0.5))
	return Quat{a, axis.X * s, axis.Y * s, axis.Z * s}
}

func ZeroQuat() Quat {
	return Quat{1.0, 0.0, 0.0, 0.0}
}
func (q Quat) conj() Quat {
	return Quat{q.a, -q.b, -q.c, -q.d}
}
func (q Quat) Norm() Quat {
	sum := float64(q.a*q.a + q.b*q.b + q.c*q.c + q.d*q.d)
	m := float32(math.Sqrt(sum))
	if m < 0.0001 {
		return Quat{1.0, 0.0, 0.0, 0.0}
	}

	return Quat{q.a / m, q.b / m, q.c / m, q.d / m}
}
func (q Quat) Add(p Quat) Quat {
	return Quat{q.a + p.a, q.b + p.b, q.c + p.c, q.d + p.d}
}
func (q Quat) Scale(k float32) Quat {
	return Quat{q.a * k, q.b * k, q.c * k, q.d * k}
}
func (q Quat) Product(p Quat) Quat {
	a := q.a*p.a - q.b*p.b - q.c*p.c - q.d*p.d
	b := q.a*p.b + q.b*p.a + q.c*p.d - q.d*p.c
	c := q.a*p.c - q.b*p.d + q.c*p.a + q.d*p.b
	d := q.a*p.d + q.b*p.c - q.c*p.b + q.d*p.a
	return Quat{a, b, c, d}
}
func (q Quat) Slerp(p Quat, w float32) Quat {
	cht := float64(q.a*p.a + q.b*p.b + q.c*p.c + q.d*p.d)
	if cht >= 1.0 || cht <= -1.0 {
		return q
	}

	ht := math.Acos(cht)
	sht := math.Sqrt(1.0 - cht*cht)
	if sht < 0.0001 && sht > -0.0001 {
		a := q.a*0.5 + p.a*0.5
		b := q.b*0.5 + p.b*0.5
		c := q.c*0.5 + p.c*0.5
		d := q.d*0.5 + p.d*0.5
		return Quat{a, b, c, d}
	}

	r := float64(w)
	rq := float32(math.Sin((1-r)*ht) / sht)
	rp := float32(math.Sin(r*ht) / sht)
	a := q.a*rq + p.a*rp
	b := q.b*rq + p.b*rp
	c := q.c*rq + p.c*rp
	d := q.d*rq + p.d*rp
	return Quat{a, b, c, d}
}
func (q Quat) ZeroSlerp(w float32) Quat {
	return q.Slerp(ZeroQuat(), w)
}
func (q Quat) Rotate(v Vec3) Vec3 {
	p := Quat{0.0, v.X, v.Y, v.Z}
	o := q.Product(p).Product(q.conj())
	return Vec3{o.b, o.c, o.d}
}
func (q Quat) Apply() {
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

// order or Euler axis vectors: Z (roll), Y (pitch), X (heading)
/*
x - heading
y - pitch
z - roll

y = asin(te8)
if abs(te8) < 0.999999:
x = atan2(-te9, te10)
z = atan2(-te4, te0)
else:
x = atan2(te6, te5)
z = 0

te0 = 1.0 - 2*q.c*q.c - 2*q.d*q.d
te1 = 2*q.b*q.c + 2*q.a*q.d
te2 = 2*q.b*q.d - 2*q.a*q.c
te4 = 2*q.b*q.c - 2*q.a*q.d
te5 = 1.0 - 2*q.b*q.b - 2*q.d*q.d
te6 = 2*q.c*q.d + 2*q.a*q.b
te8 = 2*q.b*q.d + 2*q.a*q.c
te9 = 2*q.c*q.d - 2*q.a*q.b
te10 = 1.0 - 2*q.b*q.b - 2*q.c*q.c
*/
func (q Quat) Heading() float32 {
	x := float64(2*q.c*q.d - 2*q.a*q.b)
	rad := math.Asin(-x)
	return float32(rad * 180 / math.Pi)
}
func (q Quat) Pitch() float32 {
	y := float64(2*q.b*q.d + 2*q.a*q.c)
	x := float64(1.0 - 2*q.b*q.b - 2*q.c*q.c)
	rad := math.Atan2(y, x)
	return float32(rad * 180 / math.Pi)
}
func (q Quat) Roll() float32 {
	y := float64(2*q.b*q.c + 2*q.a*q.d)
	x := float64(1.0 - 2*q.b*q.b - 2*q.d*q.d)
	rad := math.Atan2(y, x)
	return float32(rad * 180 / math.Pi)
}
