package rkt

import (
	"log"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Hud struct {
	PVMatrixPair
	width  uint16
	height uint16
	adi    *Adi
}

type Adi struct {
	pos       Vec3
	scale     float32
	tapeGeom  *Geom1
	frameGeom *Geom1
}

const adiTapeGeomName = "base/geom/aditape"
const adiFrameGeomName = "base/geom/adiframe"

func NewHud() *Hud {
	h := new(Hud)
	h.adi = buildAdi()
	return h
}

func (h *Hud) SetViewport(width, height uint16) {
	h.width = width
	h.height = height
	h.SetProjection()
}
func (h *Hud) SetProjection() {
	aspect := float32(h.width) / float32(h.height)
	h.ProjMatrix.Ortho(aspect, -1.0, 1.0)
	h.ViewMatrix.SetIdentity()
}
func (h *Hud) Draw(input Quat) {
	gl.Disable(gl.DEPTH_TEST)
	h.SetProjection()
	h.adi.draw(input)
	gl.Enable(gl.DEPTH_TEST)
}

func buildAdi() *Adi {
	v := new(Adi)
	v.pos = Vec3{-0.8, -0.8, 0.0}
	v.scale = 0.2

	tapeGeomDef, ok := geom1DefMap[adiTapeGeomName]
	if !ok {
		log.Fatalf("build_adi: no such geom1def %s\n", adiTapeGeomName)
	}

	frameGeomDef, ok := geom1DefMap[adiFrameGeomName]
	if !ok {
		log.Fatalf("build_adi: no such geom1def %s\n", adiFrameGeomName)
	}

	v.tapeGeom = tapeGeomDef.create()
	v.frameGeom = frameGeomDef.create()
	return v
}

func (v *Adi) draw(input Quat) {
	model := Matrix4{}
	model.SetIdentity()
	model.SetPos(v.pos)
	pr := input.Rotate(Vec3{0, 0, 1})
	up := input.Rotate(Vec3{1, 0, 0})
	lf := input.Rotate(Vec3{0, 1, 0})
	pitch := math.Asin(float64(Clamp(pr.Z, -1.0, 1.0)))
	roll := -math.Atan2(float64(lf.Z), float64(up.Z))
	yaw := -math.Atan2(float64(pr.X), float64(pr.Y))
	model.RotZ(float32(roll))
	model.RotX(float32(pitch - math.Pi*0.5))
	model.RotZ(float32(yaw))
	model.Scale1(v.scale)
	v.tapeGeom.draw(&model)

	model.SetIdentity()
	model.SetPos(v.pos.Add(Vec3{Z: 1.0}))
	model.Scale1(v.scale)
	v.frameGeom.draw(&model)
}
