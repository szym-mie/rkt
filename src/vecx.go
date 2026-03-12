package rkt

import (
	"encoding/json"
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type VecAxis uint

const (
	XAxis VecAxis = iota + 1
	YAxis
	ZAxis
	WAxis
)

type Vec2 struct {
	X, Y float32
}

func (v *Vec2) UnmarshalJSON(data []byte) error {
	val := new([2]float32)
	if err := json.Unmarshal(data, val); err != nil {
		return err
	}

	v.FromArray(*val)
	return nil
}
func (v *Vec2) FromArray(val [2]float32) {
	v.X = val[0]
	v.Y = val[1]
}
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{v.X + u.X, v.Y + u.Y}
}
func (v Vec2) Sub(u Vec2) Vec2 {
	return Vec2{v.X - u.X, v.Y - u.Y}
}

type Vec3 struct {
	X, Y, Z float32
}

func (v *Vec3) UnmarshalJSON(data []byte) error {
	val := new([3]float32)
	if err := json.Unmarshal(data, val); err != nil {
		return err
	}

	v.fromArray(*val)
	return nil
}
func (v *Vec3) fromArray(val [3]float32) {
	v.X = val[0]
	v.Y = val[1]
	v.Z = val[2]
}
func (v Vec3) LenSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
func (v Vec3) Len() float32 {
	lenSq := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return float32(math.Sqrt(float64(lenSq)))
}
func (v Vec3) Apply() {
	gl.Translatef(v.X, v.Y, v.Z)
}
func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z}
}
func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z}
}
func (v Vec3) Scale(k float32) Vec3 {
	return Vec3{v.X * k, v.Y * k, v.Z * k}
}
func (v Vec3) Norm() Vec3 {
	lenSqr := float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if lenSqr < 0.00001 {
		return Vec3{}
	}

	return v.Scale(1 / float32(math.Sqrt(lenSqr)))
}
func (v Vec3) Product(u Vec3) Vec3 {
	return Vec3{v.X * u.X, v.Y * u.Y, v.Z * u.Z}
}
func (v Vec3) Dot(u Vec3) float32 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}
func (v Vec3) Cross(u Vec3) Vec3 {
	x := v.Y*u.Z - v.Z*u.Y
	y := v.Z*u.X - v.X*u.Z
	z := v.X*u.Y - v.Y*u.X
	return Vec3{x, y, z}
}
func (v Vec3) Ortho() Vec3 {
	other := Vec3{}
	x := math.Abs(float64(v.X))
	y := math.Abs(float64(v.Y))
	z := math.Abs(float64(v.Z))
	if x < y {
		if x < z {
			other.X = 1.0
		} else {
			other.Z = 1.0
		}
	} else {
		if y < z {
			other.Y = 1.0
		} else {
			other.Z = 1.0
		}
	}

	return v.Cross(other)
}

func Max(x float32, y float32) float32 {
	if x < y {
		return y
	}

	return x
}

func Clamp(x float32, min, max float32) float32 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
