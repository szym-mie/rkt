package rkt

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
)

type Hud struct {
	width  uint16
	height uint16
	vdi    *Vdi
}

type Vdi struct {
	pos       Vec2
	scale     float32
	tapeGeom  *Geom1
	datumGeom *Geom1
}

const vdiTapeName = "base/hud/vditape"
const vdiDatumName = "base/hud/vdidatum"

func (c *Hud) SetViewport(width, height uint16) {
	c.width = width
	c.height = height
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.LoadIdentity()
	f := float64(width)/float64(height) - 1.0
	// same as default projection, but with aspect correction
	gl.Ortho(-1.0-f, 1.0+f, -1.0, 1.0, -1.0, 1.0)
	gl.PushMatrix()
}

func NewHud() *Hud {
	h := new(Hud)
	h.vdi = buildVdi()
	return h
}

func buildVdi() *Vdi {
	v := new(Vdi)
	v.pos = Vec2{-0.9, -0.9}
	v.scale = 0.2

	tapeTexture, ok := textureMap[vdiTapeName]
	if !ok {
		log.Fatalf("build_vdi: no such texture %s\n", vdiTapeName)
	}
	tapeTexture.setRepeat(false)
	v.tapeGeom = NewGeom1(tapeTexture, 2)
	v.tapeGeom.Vertices = buildQuadVec3(
		Vec3{-4.0, -4.0, 0.0}, Vec3{4.0, 4.0, 0.0})
	v.tapeGeom.TexCoords = buildQuadVec2(
		Vec2{-1.5, -1.5}, Vec2{2.5, 2.5})

	datumTexture, ok := textureMap[vdiDatumName]
	if !ok {
		log.Fatalf("build_vdi: no such texture %s\n", vdiDatumName)
	}
	datumTexture.setRepeat(false)
	v.datumGeom = NewGeom1(datumTexture, 2)
	v.datumGeom.Vertices = buildQuadVec3(
		Vec3{-1.0, -1.0, 0.0}, Vec3{1.0, 1.0, 0.0})
	v.datumGeom.TexCoords = buildQuadVec2(
		Vec2{0.0, 0.0}, Vec2{1.0, 1.0})
	return v
}

func (v *Vdi) draw(input Quat) {
	// roll := input.Roll()
	// pitch := input.Pitch()
}
