package rkt

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type Bitmap image.RGBA

type TextureFilter int

const (
	TextureFilterNearest = iota
	TextureFilterLinear
)

type Texture uint32

func (t Texture) setFilter(filterType TextureFilter) {
	t.bind()
	switch filterType {
	case TextureFilterNearest:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	case TextureFilterLinear:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
}
func (t Texture) setRepeat(repeatEnable bool) {
	t.bind()
	var param int32
	if repeatEnable {
		param = gl.REPEAT
	} else {
		param = gl.CLAMP
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, param)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, param)
}
func (t Texture) bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}
func (t Texture) bindTo(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}

func InitTextureUnit(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.Enable(gl.TEXTURE_2D)
}

func (b *Bitmap) createTexture() Texture {
	var handle uint32
	size := b.Rect.Size()
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &handle)
	texture := Texture(handle)
	texture.bind()
	texture.setFilter(TextureFilterNearest)
	texture.setRepeat(true)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(b.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}
