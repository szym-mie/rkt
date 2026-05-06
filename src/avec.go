package rkt

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v3.3-core/gl"
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
func (v Vec3) AxisLenSq() Vec3 {
	x, y, z := v.X*v.X, v.Y*v.Y, v.Z*v.Z
	return Vec3{y + z, x + z, x + y}
}
func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z}
}
func (v Vec3) AddSca(k float32) Vec3 {
	return Vec3{v.X + k, v.Y + k, v.Z + k}
}
func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z}
}
func (v Vec3) Norm() Vec3 {
	lenSqr := float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if lenSqr < 0.00001 {
		return Vec3{}
	}

	return v.MulSca(1 / float32(math.Sqrt(lenSqr)))
}
func (v Vec3) Mul(u Vec3) Vec3 {
	return Vec3{v.X * u.X, v.Y * u.Y, v.Z * u.Z}
}
func (v Vec3) MulSca(k float32) Vec3 {
	return Vec3{v.X * k, v.Y * k, v.Z * k}
}
func (v Vec3) Div(u Vec3) Vec3 {
	x, y, z := float32(0.0), float32(0.0), float32(0.0)
	if u.X != 0.0 {
		x = v.X / u.X
	}
	if u.Y != 0.0 {
		y = v.Y / u.Y
	}
	if u.Z != 0.0 {
		z = v.Z / u.Z
	}
	return Vec3{x, y, z}
}
func (v Vec3) Lerp(u Vec3, w float32) Vec3 {
	return v.MulSca(1.0 - w).Add(u.MulSca(w))
}
func (v Vec3) RandomSphere() Vec3 {
	x := rand.Float32()*2.0*v.X - v.X
	y := rand.Float32()*2.0*v.Y - v.Y
	z := rand.Float32()*2.0*v.Z - v.Z
	return Vec3{x, y, z}
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
	x := Absf(v.X)
	y := Absf(v.Y)
	z := Absf(v.Z)
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
func (v Vec3) To4(w float32) Vec4 {
	return Vec4{v.X, v.Y, v.Z, w}
}
func (v Vec3) uniform(location int32) {
	gl.Uniform3f(location, v.X, v.Y, v.Z)
}

type Vec4 struct {
	X, Y, Z, W float32
}

func (v Vec4) LenSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
}
func (v Vec4) Len() float32 {
	lenSq := v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
	return float32(math.Sqrt(float64(lenSq)))
}
func (v Vec4) AxisLenSq() Vec4 {
	x, y, z, w := v.X*v.X, v.Y*v.Y, v.Z*v.Z, v.W*v.W
	xw, yz := x+w, y+z
	return Vec4{yz + w, xw + z, xw + y, yz + x}
}
func (v Vec4) Add(u Vec4) Vec4 {
	return Vec4{v.X + u.X, v.Y + u.Y, v.Z + u.Z, v.W + u.W}
}
func (v Vec4) AddSca(k float32) Vec4 {
	return Vec4{v.X + k, v.Y + k, v.Z + k, v.W + k}
}
func (v Vec4) Sub(u Vec4) Vec4 {
	return Vec4{v.X - u.X, v.Y - u.Y, v.Z - u.Z, v.W - u.W}
}
func (v Vec4) Norm() Vec4 {
	lenSqr := float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if lenSqr < 0.00001 {
		return Vec4{}
	}

	return v.MulSca(1 / float32(math.Sqrt(lenSqr)))
}
func (v Vec4) Mul(u Vec4) Vec4 {
	return Vec4{v.X * u.X, v.Y * u.Y, v.Z * u.Z, v.W * u.W}
}
func (v Vec4) MulSca(k float32) Vec4 {
	return Vec4{v.X * k, v.Y * k, v.Z * k, v.W * k}
}
func (v Vec4) Div(u Vec4) Vec4 {
	x, y, z, w := float32(0.0), float32(0.0), float32(0.0), float32(0.0)
	if u.X != 0.0 {
		x = v.X / u.X
	}
	if u.Y != 0.0 {
		y = v.Y / u.Y
	}
	if u.Z != 0.0 {
		z = v.Z / u.Z
	}
	if u.W != 0.0 {
		w = v.W / u.W
	}
	return Vec4{x, y, z, w}
}
func (v Vec4) Lerp(u Vec4, w float32) Vec4 {
	return v.MulSca(1.0 - w).Add(u.MulSca(w))
}
func (v Vec4) Dot(u Vec4) float32 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z + v.W*u.W
}
func (v Vec4) uniform(location int32) {
	gl.Uniform4f(location, v.X, v.Y, v.Z, v.W)
}

func Absf(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Minf(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
}

func Maxf(x, y float32) float32 {
	if x < y {
		return y
	}
	return x
}

func Clampf(x float32, min, max float32) float32 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
