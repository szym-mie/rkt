package rkt

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type Quat struct {
	x, y, z, w float32
}

func (q Quat) apply() {
	angle := float32(2.0 * math.Acos(float64(q.w)))
	comp := float32(math.Sqrt(float64(1.0 - q.w*q.w)))
	if comp > 0.001 {
		x := q.x / comp
		y := q.y / comp
		z := q.z / comp

		gl.Rotatef(angle, x, y, z)
	}
}
