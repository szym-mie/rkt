package rkt

import (
	"math"
)

func buildRingVertices(ringCount int, radius, z float32) []Vec3 {
	ringVts := make([]Vec3, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		a := d * 2 * math.Pi
		x := radius * float32(math.Cos(a))
		y := radius * float32(math.Sin(a))
		ringVts[i] = Vec3{x, y, z}
	}

	return ringVts
}

func buildRingEndTexCoords(ringCount int, pageX uint8) []Vec2 {
	ringTCs := make([]Vec2, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		a := d * 2 * math.Pi
		x := float32(math.Cos(a))*0.25 + 0.25
		y := float32(math.Sin(a))*0.25 + 0.25
		ringTCs[i] = Vec2{x + float32(pageX)*0.5, y}
	}

	return ringTCs
}

func buildRingSideTexCoords(ringCount int, pageY uint8) []Vec2 {
	sideTCs := make([]Vec2, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		x := float32(math.Abs(d*2.0 - 1.0))
		y := float32(pageY) * 0.5
		sideTCs[i] = Vec2{x, y}
	}

	return sideTCs
}

func buildTextVertices(len int) []Vec3 {
	count := len * 6
	textVts := make([]Vec3, count)
	for i := range len {
		j := i * 6
		d := float32(i) - float32(len-1)*0.5
		textVts[j+0] = Vec3{-0.5 + d, -0.5, 0.0}
		textVts[j+1] = Vec3{+0.5 + d, -0.5, 0.0}
		textVts[j+2] = Vec3{+0.5 + d, +0.5, 0.0}
		textVts[j+3] = Vec3{+0.5 + d, +0.5, 0.0}
		textVts[j+4] = Vec3{-0.5 + d, +0.5, 0.0}
		textVts[j+5] = Vec3{-0.5 + d, -0.5, 0.0}
	}

	return textVts
}

func buildQuadVec2(lower Vec2, upper Vec2) []Vec2 {
	pts := make([]Vec2, 6)
	uu := upper
	ll := lower
	ul := Vec2{uu.X, ll.Y}
	lu := Vec2{ll.X, uu.Y}
	pts[0] = ll
	pts[1] = ul
	pts[2] = uu
	pts[3] = uu
	pts[4] = lu
	pts[5] = ll
	return pts
}

func buildQuadVec3(lower Vec3, upper Vec3) []Vec3 {
	pts := make([]Vec3, 6)
	uu := upper
	ll := lower
	ul := Vec3{uu.X, ll.Y, uu.Z}
	lu := Vec3{ll.X, uu.Y, uu.Z}
	pts[0] = ll
	pts[1] = ul
	pts[2] = uu
	pts[3] = uu
	pts[4] = lu
	pts[5] = ll
	return pts
}
