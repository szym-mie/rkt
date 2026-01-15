package rkt

import (
	"encoding/json"

	"github.com/go-gl/gl/v2.1/gl"
)

type Vec2 struct {
	x, y float32
}

type Vec3 struct {
	x, y, z float32
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
	v.x = val[0]
	v.y = val[1]
	v.z = val[2]
}

func (v Vec3) apply() {
	gl.Translatef(v.x, v.y, v.z)
}

func (v Vec3) add(u Vec3) Vec3 {
	return Vec3{v.x + u.x, v.y + u.y, v.z + u.z}
}
