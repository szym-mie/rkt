package rkt

import (
	"unicode/utf8"
)

const fontDx = 1.0 / 16.0
const fontDy = 1.0 / 8.0

type Label struct {
	Pos   Vec3
	Rot   Quat
	Scale Vec2
	len   int
	geom  *Geom1
}

const fontShaderName = "base/glsl/phong1"

func NewLabel(fontName string, len int) *Label {
	t := new(Label)
	// TODO: gl3
	// shader, ok := shaderMap[fontShaderName]
	// if !ok {
	// 	log.Fatalf("build_geom: no such shader %s\n", fontShaderName)
	// }

	// texture, ok := textureMap[fontName]
	// if !ok {
	// 	log.Fatalf("new_label: no such texture %s", fontName)
	// }

	// t.len = len
	// t.geom = NewGeom1(shader, texture, len*2)
	// copy(t.geom.Vertices, buildTextVertices(len))
	t.Rot = ZeroQuat()
	t.Scale = Vec2{0.5, 1.0}
	return t
}

func NewLabelFor(fontName string, msg string) *Label {
	len := utf8.RuneCountInString(msg)
	l := NewLabel(fontName, len)
	l.Write(msg)
	return l
}

func (l *Label) Write(msg string) {
	i := 0
	for _, c := range msg[:l.len] {
		if c > '\x7f' {
			c = '\x7f'
		}

		// TODO: gl3
		// xl := float32(c%16) * fontDx
		// yl := float32(c/16) * fontDy
		// j := i * 6
		// xu := xl + fontDx
		// yu := yl + fontDy
		// l.geom.TexCoords[j+0] = Vec2{xl, yu}
		// l.geom.TexCoords[j+1] = Vec2{xu, yu}
		// l.geom.TexCoords[j+2] = Vec2{xu, yl}
		// l.geom.TexCoords[j+3] = Vec2{xu, yl}
		// l.geom.TexCoords[j+4] = Vec2{xl, yl}
		// l.geom.TexCoords[j+5] = Vec2{xl, yu}
		i++
	}
}
func (l *Label) Draw() {
	// gl.MatrixMode(gl.MODELVIEW)
	// gl.PushMatrix()
	// l.Rot.Apply()
	// l.Pos.Apply()
	// gl.Scalef(l.Scale.X, l.Scale.Y, 0.0)
	// l.geom.draw()
	// gl.PopMatrix()
}
