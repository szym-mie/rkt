package rkt

var drawVectorGeom *Geom1
var drawDiamondGeom *Geom1
var lineColor Vec3

func InitDraw() {
	drawVectorGeom = geom1DefMap["base/geom/drawvector"].create()
	drawDiamondGeom = geom1DefMap["base/geom/drawdiamond"].create()
}

func SetLineColor(r, g, b float32) {
	lineColor = Vec3{r, g, b}
}

func DrawVector(dir, pt Vec3) {
	// gl.Disable(gl.DEPTH_TEST)
	// gl.PushMatrix()
	// gl.Color4f(lineColor.X, lineColor.Y, lineColor.Z, 1.0)
	// q := NewVecDiffQuat(dir.Norm(), Vec3{0, 0, 1}).Norm()
	// pt.Apply()
	// q.Conj().Apply()
	// gl.Scalef(1.0, 1.0, dir.Len())
	// drawVectorGeom.draw()
	// gl.Color4f(1.0, 1.0, 1.0, 1.0)
	// gl.PopMatrix()
	// gl.Enable(gl.DEPTH_TEST)
}

func DrawDiamond(pt Vec3) {
	// gl.Disable(gl.DEPTH_TEST)
	// gl.PushMatrix()
	// gl.Color4f(lineColor.X, lineColor.Y, lineColor.Z, 1.0)
	// pt.Apply()
	// drawDiamondGeom.draw()
	// gl.Color4f(1.0, 1.0, 1.0, 1.0)
	// gl.PopMatrix()
	// gl.Enable(gl.DEPTH_TEST)
}
