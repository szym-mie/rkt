package rkt

import (
	"encoding/json"

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
	lenSqr := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	if lenSqr < 0.00001 {
		return Vec3{}
	}

	return v.Scale(1 / lenSqr)
}
