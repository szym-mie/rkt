package rkt

import (
	"log"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const quadGeomName = "base/geom/quad"

type Hud struct {
	PVMatrix
	width   uint16
	height  uint16
	quadDef *Geom1Def
	adi     *Adi
	chelp   *CtrlsHelper
}

type Adi struct {
	pos       Vec3
	scale     float32
	ballGeom  *Geom1
	frameGeom *Geom1
}

type CtrlsHelper struct {
	pos          Vec3
	scale        float32
	gridGeom     *Geom1
	stickGeom    *Geom1
	rudderGeom   *Geom1
	throttleGeom *Geom1
}

func NewHud() *Hud {
	h := new(Hud)
	quadGeomDef, ok := geom1DefMap[quadGeomName]
	if !ok {
		log.Fatalf("new_hud: no such geom1def %s\n", quadGeomName)
	}

	h.quadDef = quadGeomDef
	h.adi = buildAdi(h)
	h.chelp = buildCtrlsHelper(h)
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
func (h *Hud) Draw(orient Quat, input Vec4) {
	gl.Disable(gl.DEPTH_TEST)
	h.SetProjection()
	h.adi.draw(orient)
	h.chelp.draw(input)
	gl.Enable(gl.DEPTH_TEST)
}

const adiBallGeomName = "base/geom/adiball"
const adiFrameTextureName = "base/hud/adiframe"

func buildAdi(hud *Hud) *Adi {
	a := new(Adi)
	a.pos = Vec3{-0.8, -0.8, 0.0}
	a.scale = 0.2

	ballGeomDef, ok := geom1DefMap[adiBallGeomName]
	if !ok {
		log.Fatalf("build_adi: no such geom1def %s\n", adiBallGeomName)
	}

	a.ballGeom = ballGeomDef.create()
	a.frameGeom = hud.quadDef.toQuad(getTexture(adiFrameTextureName))
	return a
}

func (a *Adi) draw(orient Quat) {
	model := Mat4{}
	model.SetIdentity()
	model.SetPos(a.pos)
	pr := orient.Rotate(Vec3{1, 0, 0})
	up := orient.Rotate(Vec3{0, 0, 1})
	lf := orient.Rotate(Vec3{0, 1, 0})
	pitch := math.Asin(float64(Clampf(pr.Z, -1.0, 1.0)))
	roll := math.Atan2(float64(lf.Z), float64(up.Z))
	yaw := -math.Atan2(float64(pr.X), float64(pr.Y))
	model.RotZ(float32(roll))
	model.RotX(float32(pitch - math.Pi*0.5))
	model.RotZ(float32(yaw))
	model.Scale1(a.scale)
	a.ballGeom.draw(&model)

	model.SetIdentity()
	model.SetPos(a.pos.Add(Vec3{Z: 1.0}))
	model.Scale1(a.scale)
	a.frameGeom.draw(&model)
}

const ctrlsHelperGridTextureName = "base/hud/stick_grid"
const ctrlsHelperStickTextureName = "base/hud/stick_pitchroll"
const ctrlsHelperRudderTextureName = "base/hud/stick_yaw"
const ctrlsHelperThrottleTextureName = "base/hud/stick_throttle"

func buildCtrlsHelper(hud *Hud) *CtrlsHelper {
	c := new(CtrlsHelper)
	c.pos = Vec3{0.9, -0.9, 0.0}
	c.scale = 0.1

	c.gridGeom = hud.quadDef.toQuad(getTexture(ctrlsHelperGridTextureName))
	c.stickGeom = hud.quadDef.toQuad(getTexture(ctrlsHelperStickTextureName))
	c.rudderGeom = hud.quadDef.toQuad(getTexture(ctrlsHelperRudderTextureName))
	c.throttleGeom = hud.quadDef.toQuad(getTexture(ctrlsHelperThrottleTextureName))
	return c
}

func (c *CtrlsHelper) draw(input Vec4) {
	model := Mat4{}
	model.SetIdentity()
	model.SetPos(c.pos)
	model.Scale1(c.scale)
	c.gridGeom.draw(&model)

	deflScale := c.scale * 0.8
	model.SetIdentity()
	model.SetPos(c.pos.Add(Vec3{-input.Y * deflScale, input.X * deflScale, 1.0}))
	model.Scale1(c.scale)
	c.stickGeom.draw(&model)

	model.SetIdentity()
	model.SetPos(c.pos.Add(Vec3{input.Z * deflScale, 0.0, 1.0}))
	model.Scale1(c.scale)
	c.rudderGeom.draw(&model)

	model.SetIdentity()
	throttlePos := 2.0*input.W - 1.0
	model.SetPos(c.pos.Add(Vec3{0.0, throttlePos * deflScale, 1.0}))
	model.Scale1(c.scale)
	c.throttleGeom.draw(&model)
}
